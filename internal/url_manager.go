package internal

import (
	"errors"
	"net/url"
	"sync"
)

type UrlsManager struct {
	// Urls are placed in a queue and then passed to workers one at a time.
	// This queue may be added to over the course of the program if links should
	// be followed. A mutex is necessary to ensure that no race conditions occur when
	// passing urls to workers
	Queue     []string
	Completed map[string]int // map of completed requests and their status codes

	qmu *sync.Mutex

	// Any links discovered when parsing html can be passed to the supervisor to be added to the queue.
	// All pipeline funcs have access to this channel to allow for custom rules when following links.
	discovered DiscoveredChan

	isClosed bool // Indicates whether this queue is closed

	errors chan error // Errors channel
}

func (u *UrlsManager) IsComplete() bool {
	// How to conditionally incorporate isClosed?
	return len(u.Queue) == 0 // && isClosed
}

// Closes the current queue. This indicates that we will no longer accept new urls
func (u *UrlsManager) Close() {
	u.isClosed = true
}

// Empty the queue
func (u *UrlsManager) Empty() {
	u.Queue = []string{}
}

func (u *UrlsManager) ShiftQueue() (s string, ok bool) {
	u.qmu.Lock()
	if len(u.Queue) > 0 {
		s, ok = u.Queue[0], true
		u.Queue = u.Queue[1:]
	}
	u.qmu.Unlock()
	return s, ok
}

// Queue a link if it hasn't already been sent. QueueWork will be called once the link has
// successfully been added
func (u *UrlsManager) QueueLink(link *url.URL, queueWork func()) error {
	if !u.isClosed {
		if _, ok := u.Completed[link.String()]; !ok {
			// Make sure link isn't already in queue
			// TODO maybe use something other than a slice here
			for _, qUrl := range u.Queue {
				if qUrl == link.String() {
					return errors.New("Link " + qUrl + " already in queue, ignoring")
				}
			}

		} else {
			return errors.New("Already processed link " + link.String() + ", ignoring")
		}

		u.qmu.Lock()
		u.Queue = append(u.Queue, link.String())
		u.qmu.Unlock()

		queueWork()
	}
	return nil
}

// Queue a list of strings without checking if they have been processed or are in the queue already
func (u *UrlsManager) QueueStrings(links ...string) {
	if !u.isClosed {
		u.qmu.Lock()
		u.Queue = append(u.Queue, links...)
		u.qmu.Unlock()
	}
}

func (u *UrlsManager) RecordResponse(c CompletedRequest) {
	u.Completed[c.Url] = c.StatusCode
}

func (u *UrlsManager) Monitor(shouldAdd LinkJudge, queueWork func()) {
	// Add links from discovered chan to URL queue
	go func() {
		f, ok := <-u.discovered
		for ok {
			Log.Info("Found link " + f.Url.String())

			if shouldAdd(f) {
				if err := u.QueueLink(f.Url, queueWork); err != nil {
					u.errors <- err
				}
			}

			f, ok = <-u.discovered
		}
	}()
}
