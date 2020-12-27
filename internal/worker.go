package internal

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type worker struct {
	id       int
	queue    chan string
	errors   chan error
	state    chan WorkerState
	sent     chan CompletedRequest
	found    chan FoundUrl
	pipeline []PipelineFunc
	config   *Config
}

type PipelineFunc func([]byte, string, chan FoundUrl) ([]byte, error)

type WorkerState int

var (
	WorkerStateInactive  WorkerState = 0 // No work for this worker to do
	WorkerStateWaiting   WorkerState = 1 // Pipeline running but waiting for another url to process
	WorkerStateRunning   WorkerState = 2 // Currently processing a url
	WorkerStateStopping  WorkerState = 3 // Received stop signal, waiting for pipelines to finish
	WorkerStateFinished  WorkerState = 4 // Finished all pending work
	WorkerStateReceiving WorkerState = 5 // Receiving work from supervisor
)

type CompletedRequest struct {
	U string
	S int
}

func (w *worker) Start() {
	notifier := make(chan bool)
	stopping := false
	isWorking := false
	pipelines := 0

	// Monitor the state of all pipelines to prevent exiting before work is done
	go func() {
		done, ok := <-notifier
		for ok {
			if done {
				pipelines -= 1
				if pipelines == 0 {
					if !isWorking {
						w.state <- WorkerStateInactive
					}
					if stopping {
						w.state <- WorkerStateFinished
						break
					}
				}
			}
			done, ok = <-notifier
		}
	}()

	w.state <- WorkerStateInactive
	u, ok := <-w.queue

	// Main loop for worker threads. Send a request and pass the response to the pipeline for processing
	for ok {
		isWorking = true
		w.state <- WorkerStateRunning

		if resp, err := SendRequest(u, w.config); err == nil {
			w.sent <- CompletedRequest{u, resp.StatusCode}
			pipelines += 1
			go w.runPipeline(resp, u, notifier)
		} else {
			w.errors <- err
		}

		time.Sleep(w.config.Delay)

		if pipelines > 0 {
			w.state <- WorkerStateWaiting
		} else {
			w.state <- WorkerStateInactive
		}

		isWorking = false
		u, ok = <-w.queue
	}

	// Set the state to 'stopping' so that the notifier knows what state to send to supervisor
	// when all pipelines complete
	w.state <- WorkerStateStopping
	stopping = true

	// If there are no pipelines open at this point, we can go ahead and stop the worker
	if pipelines == 0 {
		close(notifier)
		w.state <- WorkerStateFinished
	}
}

func (w *worker) runPipeline(r *http.Response, u string, notifier chan bool) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.errors <- err
		notifier <- true
		return
	}
	_ = r.Body.Close()

	for _, step := range w.pipeline {
		b, err = step(b, u, w.found)
		if err != nil {
			w.errors <- err
		}
	}

	parsed, err := ParseUrl(u)
	if err != nil {
		w.errors <- err
		return
	}

	// Write file
	err = w.writeBytes(b, parsed)
	if err != nil {
		w.errors <- err
	}

	notifier <- true
}

// Writes bytes to specified output, whether that is the file system or a remote url
func (w *worker) writeBytes(b []byte, uri *url.URL) error {
	if strings.HasPrefix(w.config.Output, "http") {
		req, err := http.NewRequest("POST", w.config.Output+uri.Path, bytes.NewReader(b))
		if err != nil {
			return err
		}

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
	} else {
		filename := uri.Path
		// If this is not a valid url, attempt to write to file system in a directory with the host name
		if len(uri.Path) == 0 || uri.Path[len(uri.Path)-1] == '/' {
			filename = "index.html"
		}
		p := path.Join(w.config.Output, uri.Hostname(), filename)

		// Create directories
		_, err := os.Stat(path.Dir(p))
		if os.IsNotExist(err) {
			err := os.MkdirAll(path.Dir(p), 0755)
			if err != nil {
				return err
			}
		}

		f, err := os.Create(p)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}

		_ = f.Close()
	}
	return nil
}
