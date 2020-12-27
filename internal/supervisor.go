package internal

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
)

// The supervisor should only be created once for each program invocation.
// This struct manages the state of all worker threads, passes work to them,
// and maintains the state of the program. All work is started via Start() and
// can be stopped gracefully with Terminate()
type Supervisor struct {
	// Pipeline steps to call for each request. Each worker will use the same
	// pipeline and creates a new goroutine for each completed request. The result
	// of each pipeline func is passed to the next and is finally saved at the end
	Pipeline []PipelineFunc

	client *http.Client // http.Client to use when making requests

	// Urls are placed in a queue and then passed to workers one at a time.
	// This queue may be added to over the course of the program if links should
	// be followed. A mutex is necessary to ensure that no race conditions occur when
	// passing urls to workers
	urlQueue []string
	qMutex   *sync.Mutex

	sent     chan CompletedRequest // Channel used by workers to notify supervisor of completed requests
	sentUrls map[string]int        // map of completed requests and their status codes
	// Any links found when parsing html can be passed to the supervisor to be added to the queue.
	// All pipeline funcs have access to this channel to allow for custom rules when following links
	found chan FoundUrl

	errors chan error // any non-fatal errors which occur in worker threads are passed to the supervisor

	// References to each worker thread and its state. The worker state is used to
	// determine whether the program should exit based on what all threads are doing.
	// If every thread has completed its work and there are no more urls to process,
	// then the done channel is triggered and the Start() function returns after terminating the threads.
	// Threads will gracefully exit and wait for all pipelines to complete before stopping
	workers     []*worker     // list of worker threads
	workerState []WorkerState // The current state of each worker

	config *Config   // configuration object used
	done   chan bool // Used to block the main thread until all work is done
	logger *Logger   // Logger instance, used for convenience when writing to stdout/stderr
}

type FoundUrl struct {
	Current string
	Url     string
}

func NewSupervisor(config *Config) *Supervisor {
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	return &Supervisor{
		client:      client,
		urlQueue:    []string{},
		qMutex:      &sync.Mutex{},
		sent:        make(chan CompletedRequest),
		sentUrls:    make(map[string]int),
		found:       make(chan FoundUrl),
		errors:      make(chan error),
		logger:      &Logger{os.Stdout, os.Stderr, 0},
		workers:     []*worker{},
		workerState: make([]WorkerState, config.Workers),
		config:      config,
		done:        make(chan bool, 1),
	}
}

// Generate all seed urls and start all worker threads
func (s *Supervisor) Start() error {
	s.logger.Info("Starting Anubis...")
	err := s.buildSeed()
	if err != nil {
		return err
	}

	go s.writeErrors() // writes all errors to stderr
	go s.manageUrls()  // Manages all sent urls and found urls to prevent duplicates and loops
	s.startWorkers()

	// Wait until we're done to return.
	// We're done when all workers are done processing and the url queue is empty
	<-s.done
	s.logger.Info("No more urls to process. Stopping Anubis...")
	// Call terminate to make sure all pipelines finish before program exits
	s.Terminate()

	return nil
}

// Gracefully stop all workers
func (s *Supervisor) Terminate() {
	// We terminate workers by closing their queue
	for _, worker := range s.workers {
		s.logger.Info("Stopping worker " + strconv.Itoa(worker.id))
		close(worker.queue)
		// Terminate only happens when all workers are inactive, so there are no pending pipelines to handle
	}
	close(s.sent)
	close(s.found)
	close(s.errors)
	close(s.done)
}

// Build seed urls based on configuration. Checks sitemap first if configured to add urls.
// By default, the base url will be part of the seed urls
func (s *Supervisor) buildSeed() error {
	var seed []string

	// If there are no seed urls, we will try sitemap.xml.
	// The program should exit with an error if we can't get this file
	if s.config.Sitemap != "" {
		resp, err := SendRequest(s.config.Sitemap, s.config)
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

	s.logger.Info("Found " + strconv.Itoa(len(seed)) + " seed urls")
	return nil
}

// Start each worker and monitor their progress
func (s *Supervisor) startWorkers() {
	for i := 0; i < s.config.Workers; i += 1 {
		w := &worker{
			id:       i,
			queue:    make(chan string, 1),
			errors:   s.errors,
			state:    make(chan WorkerState, 1),
			sent:     s.sent,
			found:    s.found,
			pipeline: s.Pipeline,
			config:   s.config,
		}
		s.workers = append(s.workers, w)
		s.logger.Info("Starting worker " + strconv.Itoa(i))
		go w.Start()
		go s.monitorWorker(w)
	}
}

func (s *Supervisor) writeErrors() {
	err, ok := <-s.errors

	for ok {
		s.logger.LogError(err)
		err, ok = <-s.errors
	}
}

func (s *Supervisor) manageUrls() {
	// Save all responses into the global url map to prevent duplicates
	go func() {
		res, ok := <-s.sent
		for ok {
			s.logger.Info(res.U + " " + strconv.Itoa(res.S))
			s.sentUrls[res.U] = res.S
			res, ok = <-s.sent
		}
	}()

	// Add links from found chan to URL queue
	go func() {
		f, ok := <-s.found
		for ok {
			s.logger.Info("Found link " + f.Url)
			link := f.Url

			if !UrlRegex.MatchString(link) {
				currentUrl, err := ParseUrl(f.Current)
				if err == nil {
					link = currentUrl.Scheme + "://" + path.Join(currentUrl.Hostname(), f.Url)
				} else {
					s.errors <- err
					continue
				}
			}

			s.addLink(link)

			f, ok = <-s.found
		}
	}()
}

func (s *Supervisor) addLink(u string) {
	if _, ok := s.sentUrls[u]; !ok {
		// Make sure link isn't already in queue
		// TODO maybe use something other than a slice here
		for _, qUrl := range s.urlQueue {
			if qUrl == u {
				s.logger.Info("Link " + u + " already in queue, ignoring")
				return
			}
		}

		s.qMutex.Lock()
		s.urlQueue = append(s.urlQueue, u)
		s.qMutex.Unlock()
	} else {
		s.logger.Info("Already processed link " + u + ", ignoring")
		return
	}

	// Check for idle workers to send new url to.
	// If none are idle, then the url will sit in the queue
	for _, worker := range s.workers {
		if s.workerState[worker.id] == WorkerStateInactive || s.workerState[worker.id] == WorkerStateWaiting {
			s.sendWork(worker)
			break
		}
	}
}

// Monitor the state of the worker. Determine whether to stop program or send worker a new url
func (s *Supervisor) monitorWorker(w *worker) {
	state, ok := <-w.state

	for ok {
		switch state {
		case WorkerStateInactive:
			s.checkProgramState(w.id)
			s.logger.Info("STATUS Worker " + strconv.Itoa(w.id) + ": inactive")
			fallthrough
		case WorkerStateWaiting:
			s.sendWork(w)
			// Whenever the worker receives the url, it will notify the supervisor of its new state
		case WorkerStateStopping:
			s.logger.Info("STATUS Worker " + strconv.Itoa(w.id) + ": stopping")
			s.workerState[w.id] = WorkerStateStopping
		case WorkerStateFinished:
			s.logger.Info("STATUS Worker " + strconv.Itoa(w.id) + ": finished")
			s.workerState[w.id] = WorkerStateFinished
		default:
			s.logger.Info("STATUS Worker " + strconv.Itoa(w.id) + ": running")
			s.workerState[w.id] = WorkerStateRunning
		}
		state, ok = <-w.state
	}
}

// Check the state of the program to make sure that there is still more work to be done
// This is accomplished by checking if all workers are inactive or finished.
// Then if the url queue is still empty, we should exit the program
func (s *Supervisor) checkProgramState(id int) {
	s.workerState[id] = WorkerStateInactive
	for _, state := range s.workerState {
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
		s.logger.Info("Sending " + u + " to worker " + strconv.Itoa(w.id))
		s.workerState[w.id] = WorkerStateReceiving
		s.logger.Info("STATUS Worker " + strconv.Itoa(w.id) + ": receiving")
		s.sentUrls[u] = 0 // Placeholder to prevent this link from being added again
		w.queue <- u
	}
}
