package internal

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var Log *Logger

// The supervisor should only be created once for each program invocation.
// This struct manages the state of all worker threads, passes work to them,
// and maintains the state of the program. All work is started via Start() and
// can be stopped gracefully with Terminate()
type Supervisor struct {
	// Pipeline steps to call for each request. Each worker will use the same
	// pipeline and creates a new goroutine for each completed request. The result
	// of each pipeline func is passed to the next and is finally saved at the end
	Pipeline []PipelineFunc

	// The LinkJudge function is used to determine whether a link passed to the discovered
	// channel should be added to the url queue or not.
	ShouldAddLink LinkJudge

	client *http.Client // http.Client to use when making requests

	// Urls are placed in a queue and then passed to workers one at a time.
	// This queue may be added to over the course of the program if links should
	// be followed. A mutex is necessary to ensure that no race conditions occur when
	// passing urls to workers
	urlQueue []string
	qMutex   *sync.Mutex

	sent     chan CompletedRequest // Channel used by workers to notify supervisor of completed requests
	sentUrls map[string]int        // map of completed requests and their status codes

	// Any links discovered when parsing html can be passed to the supervisor to be added to the queue.
	// All pipeline funcs have access to this channel to allow for custom rules when following links.
	discovered DiscoveredChan

	errors chan error // any non-fatal errors which occur in worker threads are passed to the supervisor

	// References to each worker thread and its state. The worker state is used to
	// determine whether the program should exit based on what all threads are doing.
	// If every thread has completed its work and there are no more urls to process,
	// then the done channel is triggered and the Start() function returns after terminating the threads.
	// Threads will gracefully exit and wait for all pipelines to complete before stopping
	workers       []*worker        // list of worker threads
	currentStates []int            // The current state of each worker
	workerStates  chan workerState // Channel for workers to notify supervisor of their state

	config *Config   // configuration object used
	done   chan bool // Used to block the main thread until all work is done
}

func NewSupervisor(config *Config) *Supervisor {
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	Log = &Logger{os.Stdout, os.Stderr, 0}

	return &Supervisor{
		ShouldAddLink: func(f DiscoveredLink) bool { return true }, // By default, all links will be added
		client:        client,
		urlQueue:      []string{},
		qMutex:        &sync.Mutex{},
		sent:          make(chan CompletedRequest),
		sentUrls:      make(map[string]int),
		discovered:    make(DiscoveredChan),
		errors:        make(chan error),
		workers:       []*worker{},
		currentStates: make([]int, config.Workers),
		workerStates:  make(chan workerState),
		config:        config,
		done:          make(chan bool, 1),
	}
}

// Generate all seed urls and start all worker threads
func (s *Supervisor) Start() error {
	Log.Info("Starting Anubis...")
	err := s.buildSeed()
	if err != nil {
		return err
	}

	var stats *os.File = nil

	if s.config.Stats != "" {
		stats, err = os.OpenFile(s.config.Stats, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer stats.Close()
		if err != nil {
			return err
		}
	}

	// Write all errors to stderr
	go func() {
		err, ok := <-s.errors

		for ok {
			Log.LogError(err)
			err, ok = <-s.errors
		}
	}()

	// Save all responses into the global url map to prevent duplicates
	go func() {
		res, ok := <-s.sent
		for ok {
			Log.Info(res.U + " " + strconv.Itoa(res.S))

			if stats != nil {
				// What else would we want to include here?
				_, err = stats.WriteString(res.U + " " + strconv.Itoa(res.S) + "\n")
				if err != nil {
					s.errors <- err
				}
			}

			s.sentUrls[res.U] = res.S
			res, ok = <-s.sent
		}
	}()

	// Add links from discovered chan to URL queue
	go func() {
		f, ok := <-s.discovered
		for ok {
			Log.Info("Found link " + f.Url.String())

			if s.ShouldAddLink(f) {
				s.addLink(f.Url.String())
			}

			f, ok = <-s.discovered
		}
	}()

	go s.monitorWorkers()
	s.startWorkers()

	// Wait until we're done to return.
	// We're done when all workers are done processing and the url queue is empty
	<-s.done
	Log.Info("No more urls to process. Stopping Anubis...")

	return nil
}

// Gracefully stop all workers. This causes each worker to exit its main loop, and once each worker's pipelines
// have completed the done channel will be triggered
func (s *Supervisor) Terminate() {
	// We terminate workers by closing their queue
	for _, worker := range s.workers {
		Log.Info("Stopping worker " + strconv.Itoa(worker.id))
		close(worker.queue)
	}
	// Empty url queue
	s.urlQueue = []string{}
}

// Build seed urls based on configuration. Checks sitemap first if configured to add urls.
// By default, the base url will be part of the seed urls
func (s *Supervisor) buildSeed() error {
	var seed []string

	// If there are no seed urls, we will try sitemap.xml.
	// The program should exit with an error if we can't get this file
	if s.config.Sitemap != "" {
		resp, err := SendRequest(s.client, s.config.Sitemap, s.config)
		if err != nil {
			return err
		}

		if resp.StatusCode < 300 {
			seed, err = ParseSiteMap(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}
		}
	}

	if len(s.config.Seed) != 0 {
		seed = append(seed, s.config.Seed...)
	}

	if s.config.Url != "" {
		seed = append(seed, s.config.Url)
	}

	s.urlQueue = seed

	if len(seed) > 0 {
		Log.Info("Found " + strconv.Itoa(len(seed)) + " seed urls")
	} else {
		return errors.New("No seed urls found. Exiting...")
	}
	return nil
}

// Start each worker and monitor their progress
func (s *Supervisor) startWorkers() {
	for i := 0; i < s.config.Workers; i += 1 {
		w := &worker{
			id:         i,
			queue:      make(chan string, 1),
			errors:     s.errors,
			state:      s.workerStates,
			sent:       s.sent,
			discovered: s.discovered,
			pipeline:   s.Pipeline,
			client:     s.client,
			config:     s.config,
		}
		s.workers = append(s.workers, w)
		Log.Info("Starting worker " + strconv.Itoa(i))
		go w.Start()
	}
}

func (s *Supervisor) addLink(u string) {
	if _, ok := s.sentUrls[u]; !ok {
		// Make sure link isn't already in queue
		// TODO maybe use something other than a slice here
		for _, qUrl := range s.urlQueue {
			if qUrl == u {
				Log.Info("Link " + u + " already in queue, ignoring")
				return
			}
		}

		s.qMutex.Lock()
		s.urlQueue = append(s.urlQueue, u)
		s.qMutex.Unlock()

		// Check for idle workers to send new url to.
		// If none are idle, then the url will sit in the queue
		for _, worker := range s.workers {
			if s.currentStates[worker.id] == WorkerStateInactive || s.currentStates[worker.id] == WorkerStateWaiting {
				s.sendWork(worker)
				break
			}
		}
	} else {
		Log.Info("Already processed link " + u + ", ignoring")
	}
}

// Monitor the state of the worker. Determine whether to stop program or send worker a new url
func (s *Supervisor) monitorWorkers() {
	state, ok := <-s.workerStates

	for ok {
		switch state.state {
		case WorkerStateInactive:
			s.checkProgramState(state.id)
			Log.Info("STATUS Worker " + strconv.Itoa(state.id) + ": inactive")
			fallthrough
		case WorkerStateWaiting:
			s.sendWork(s.workers[state.id])
			// Whenever the worker receives the url, it will notify the supervisor of its new state
		case WorkerStateStopping:
			Log.Info("STATUS Worker " + strconv.Itoa(state.id) + ": stopping")
			s.currentStates[state.id] = WorkerStateStopping
		case WorkerStateFinished:
			Log.Info("STATUS Worker " + strconv.Itoa(state.id) + ": finished")
			s.currentStates[state.id] = WorkerStateFinished
		default:
			Log.Info("STATUS Worker " + strconv.Itoa(state.id) + ": running")
			s.currentStates[state.id] = WorkerStateRunning
		}
		state, ok = <-s.workerStates
	}
}

// Check the state of the program to make sure that there is still more work to be done
// This is accomplished by checking if all workers are inactive or finished.
// Then if the url queue is still empty, we should exit the program
func (s *Supervisor) checkProgramState(id int) {
	s.currentStates[id] = WorkerStateInactive
	for _, state := range s.currentStates {
		if state != WorkerStateInactive && state != WorkerStateFinished {
			return
		}
	}

	if len(s.urlQueue) == 0 {
		s.done <- true
	}
}

// Thread-safe array shift
func (s *Supervisor) shiftQueue() (string, bool) {
	u := ""
	ok := false

	s.qMutex.Lock()
	if len(s.urlQueue) > 0 {
		u, ok = s.urlQueue[0], true
		s.urlQueue = s.urlQueue[1:]
	}
	s.qMutex.Unlock()
	return u, ok
}

// Send work to the specified worker. Immediately set the worker's state to receiving to prevent
// race condition in program's state
func (s *Supervisor) sendWork(w *worker) {
	u, ok := s.shiftQueue()
	if ok {
		Log.Info("Sending " + u + " to worker " + strconv.Itoa(w.id))
		s.currentStates[w.id] = WorkerStateReceiving
		Log.Info("STATUS Worker " + strconv.Itoa(w.id) + ": receiving")
		s.sentUrls[u] = 0 // Placeholder to prevent this link from being added again
		w.queue <- u
	}
}
