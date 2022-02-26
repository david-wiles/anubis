package anubis

import (
	"context"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Anubis struct {
	// Output specifies the root output directory.
	//
	// The output file will match the URL, so if Output is './src' and the perserved URL is
	// 'www.test.com/about.html', then the final output will be ./src/www.test.com/about.html.
	// This allows for links in the HTML page to be changed to relative paths, so when the page
	// is visited locally, no network requests will be necessary.
	Output string

	Workers int               // Workers indicates how many worker goroutines to use
	Headers map[string]string // Headers specifies all headers used during each network request

	Driver  WebDriver       // Driver is a WebDriver instance which will dictate how the network requests are made
	Handler ResponseHandler // Handler controls how the responses are handled before copied to a file

	wg    *sync.WaitGroup // wg is used to ensure that all workers finish before the program exits
	queue chan string     // queue is used to pass URLs to worker goroutines

	Cancel func() // Cancel should be called when the program should finish work
}

// NewAnubis creates an Anubis instance from the given arguments. If not overwritten, all fields will use
// their default values
func NewAnubis(options ...Option) *Anubis {
	a := &Anubis{
		Output:  ".",
		Workers: 4,
		Headers: make(map[string]string),
		Driver:  DefaultWebDriver{client: *http.DefaultClient},
		Handler: nil,
		wg:      &sync.WaitGroup{},
		queue:   make(chan string, 256),
		Cancel:  nil,
	}

	a.Handler = DefaultResponseHandler{a, make(map[string]bool)}
	for _, opt := range options {
		opt.SetOpt(a)
	}
	return a
}

// Start will create the context for the instance and start all workers. Workers will immediately begin
// making network requests and processing responses
func (a *Anubis) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.Cancel = func() {
		close(a.queue)
		cancel()
	}

	for n := 0; n < a.Workers; n++ {
		a.wg.Add(1)
		go a.worker(ctx, a.queue)
	}

	a.wg.Wait()

	return nil
}

// AddURL will push a new url to the queue if it is not a duplicate. This function may block the caller
// until the queue's buffer is not full.
//
// The function will return true if the link was added to the queue, and false otherwise.
func (a *Anubis) AddURL(url string) bool {
	// TODO duplicate filter
	a.queue <- url
	return true
}

// Commit will use git to commit the files with the output directory specified by the start options.
// If Anubis is started as a crawler, then this would commit all files changed up to that point
func (a Anubis) Commit() error {
	// Add all changes
	cmd := exec.Command("git", "add", "-A")
	if err := cmd.Run(); err != nil {
		return err
	}

	commitBuilder := strings.Builder{}
	commitBuilder.WriteRune('"')
	commitBuilder.WriteString(time.Now().Format(time.RFC3339))
	commitBuilder.WriteRune('"')

	// Add additional information about the system

	// Commit changes
	cmd = exec.Command("git", "commit", "-m", commitBuilder.String())
	return cmd.Run()
}

func processURL(url string, headers map[string]string, webdriver WebDriver, handler ResponseHandler) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := webdriver.DoRequest(req)
	if err != nil {
		return err
	}

	if err := handler.Handle(req, resp); err != nil {
		log.Println(err)
	}

	return nil
}

// Each worker will read URLs from the channel until the context is cancelled or the queue is closed.
// This is started by calling anubis.Start(), so all start URLs should be added first
func (a Anubis) worker(ctx context.Context, queue chan string) {
	defer a.wg.Done()

	var (
		url string
		ok  bool = true
	)

	for ok {
		select {
		case <-ctx.Done():
			return
		default:
			url, ok = <-queue

			if err := processURL(url, a.Headers, a.Driver, a.Handler); err != nil {
				log.Println(err)
			}
		}
	}
}
