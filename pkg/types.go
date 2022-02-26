package anubis

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type WebDriver interface {
	DoRequest(*http.Request) (*http.Response, error)
}

type DefaultWebDriver struct {
	client http.Client
}

func (driver DefaultWebDriver) DoRequest(req *http.Request) (*http.Response, error) {
	return driver.client.Do(req)
}

type ResponseHandler interface {
	Handle(*http.Request, *http.Response) error
}

// DefaultResponseHandler is assigned to Anubis if none is assigned otherwise.
// This handler will parse an HTML page and grab all script, stylesheet, and image URLs
// and add those to the queue for the anubis instance.
//
// This handler requires a reference to the Anubis instance to add it to the queue, and
// a similarly functioning crawler based on Anubis would also need to function in the same way.
//
// The DefaultResponseHandler also determines when the instance is finished based on default
// behavior. If the currently-handled response is from a start URL, we will store all the found
// links on that page to check after each successive step. Once all links are retrieved, the
// cancel function will be called and the program will exit.
type DefaultResponseHandler struct {
	Anubis *Anubis // The owning Anubis instance

	// NeededLinks stores each link which should be processed before the program can exit.
	// Once a link is processed, it is removed from the map. Finally, when the map is empty
	// we can call the cancel function
	NeededLinks map[string]bool
}

func (handler DefaultResponseHandler) Handle(req *http.Request, resp *http.Response) error {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check whether this is an HTML response. If it is, then we should initialize the conditions
	// to stop the anubis instance once all files have been downloaded
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		parentURL := req.URL.String()
		bodyString := string(body)

		urls := []string{}

		// Get all links
		urls = append(urls, GetLinkURLs(parentURL, bodyString)...)
		urls = append(urls, GetScriptURLs(parentURL, bodyString)...)
		urls = append(urls, GetImageURLs(parentURL, bodyString)...)

		for _, u := range urls {
			if handler.Anubis.AddURL(u) {
				handler.NeededLinks[u] = true
			}
		}
	}

	delete(handler.NeededLinks, req.URL.String())

	// Signal to all workers to exit if all work is finished
	if len(handler.NeededLinks) == 0 {
		handler.Anubis.Cancel()
	}

	filename := req.URL.Path

	// If this is not a valid url, attempt to write to file system in a directory with the host name
	if len(filename) == 0 || filename[len(filename)-1] == '/' {
		filename = "index.html"
	}

	p := path.Join(handler.Anubis.Output, req.URL.Hostname(), filename)

	if err := os.MkdirAll(path.Dir(p), 0774); err != nil {
		return err
	}

	if err := os.WriteFile(p, body, 0644); err != nil {
		return err
	}

	return nil
}

type RequestProcessor interface {
	Process(string, map[string]string, WebDriver, ResponseHandler) error
}

type DefaultRequestProcessor struct{}

func (*DefaultRequestProcessor) Process(url string, headers map[string]string, webdriver WebDriver, handler ResponseHandler) error {
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

type DuplicateFilter interface {
	TestURL(string) bool
}

// DefaultDuplicateFilter uses a *sync.Map to store URLs in memory
type DefaultDuplicateFilter struct {
	store *sync.Map
}

// TestURL uses the LoadOrStore function in sync.Map to simultaneously determine whether the key exists and set its
// value. If the key was already in the map, we return true, otherwise we'll return false
func (filter *DefaultDuplicateFilter) TestURL(u string) bool {
	_, exists := filter.store.LoadOrStore(u, true)
	return exists
}
