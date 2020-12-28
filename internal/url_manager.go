package internal

import (
	"errors"
	"net/url"
	"sync"
)

// The UrlsManager should be used to track work to be done and completed work. The work to be queued should
// not be repeated if it has already been done or is already queued, and all responses from workers will
// notify the urls manager through record response to allow for the global state to update.
// One possible reason to overwrite the default UrlsManager would be to track the state with a database
// to share the queue and completed work between systems and persist the results more easily
type UrlsManager interface {
	// Indicates whether work is complete or we can expect more work to be sent.
	// This method is called whenever all workers are idle to determine
	// whether the program should exit
	IsComplete() bool

	// Close the manager to new links. Once closed, no new links will be added to
	// the queue but all url currently in the queue will be tested
	CloseQueue()

	// Empty the current queue. Useful when terminating the program gracefully
	EmptyQueue()

	// Removes a single url from the front of the queue. This must be thread-safe
	// since the queue is shared between all workers
	ShiftQueue() (string, bool)

	// Place a new url into the queue. This must be thread safe.
	QueueLinks(...*url.URL) []error

	// Saves the response from a completed request
	RecordResponse(CompletedRequest)
}

// Default implementation is an in-memory array for the queue and map to track history of sent requests
type DefaultUrlsManager struct {
	// Urls are placed in a queue and then passed to workers one at a time.
	// This queue may be added to over the course of the program if links should
	// be followed. A mutex is necessary to ensure that no race conditions occur when
	// passing urls to workers
	Queue     []string
	Completed map[string]int // map of completed requests and their status codes

	qmu *sync.Mutex

	isClosed bool // Indicates whether this queue is closed
}

func (u *DefaultUrlsManager) IsComplete() bool {
	return len(u.Queue) == 0
}

// Closes the current queue. This indicates that we will no longer accept new urls
func (u *DefaultUrlsManager) CloseQueue() {
	u.isClosed = true
}

// Empty the queue
func (u *DefaultUrlsManager) EmptyQueue() {
	u.Queue = []string{}
}

func (u *DefaultUrlsManager) ShiftQueue() (s string, ok bool) {
	u.qmu.Lock()
	if len(u.Queue) > 0 {
		s, ok = u.Queue[0], true
		u.Queue = u.Queue[1:]
	}
	u.qmu.Unlock()
	return s, ok
}

// Queue a link if it hasn't already been sent. QueueWork will be called once the link has
// successfully been added. This function will also place a placeholder response in the map
// to prevent it from being added again after it's removed from the queue
func (u *DefaultUrlsManager) QueueLinks(links ...*url.URL) []error {
	errs := []error{}
	if !u.isClosed {
		for _, link := range links {
			if _, ok := u.Completed[link.String()]; !ok {
				// Placeholder response
				u.Completed[link.String()] = 0
				u.qmu.Lock()
				u.Queue = append(u.Queue, link.String())
				u.qmu.Unlock()
			} else {
				errs = append(errs, errors.New("Already processed link "+link.String()+", ignoring"))
			}
		}
		return errs
	} else {
		return []error{errors.New("Queue is not accepting new urls")}
	}
}

func (u *DefaultUrlsManager) RecordResponse(c CompletedRequest) {
	u.Completed[c.Url] = c.StatusCode
}
