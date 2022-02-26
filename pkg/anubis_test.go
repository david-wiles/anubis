package anubis

import (
	"net/http"
	"reflect"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

type NopWebDriver struct{}

func (NopWebDriver) DoRequest(*http.Request) (*http.Response, error) {
	return nil, nil
}

type NopResponseHandler struct{}

func (NopResponseHandler) Handle(*http.Request, *http.Response) error {
	return nil
}

func NewTestAnubis() *Anubis {
	a := NewAnubis()
	a.Driver = NopWebDriver{}
	a.Handler = NopResponseHandler{}
	a.processor = &StringProcessor{mu: &sync.Mutex{}}
	return a
}

func TestAnubis_AddURL(t *testing.T) {
	tests := []struct {
		name string
		urls []string
		want []bool
	}{
		{
			name: "Returns true when URL was added",
			urls: []string{"test"},
			want: []bool{true},
		},
		{
			name: "Returns true for each URL added",
			urls: []string{"a1", "b2", "c3"},
			want: []bool{true, true, true},
		},
		{
			name: "Returns false when a duplicate URL is added",
			urls: []string{"a1", "b2", "a1"},
			want: []bool{true, true, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnubis()

			got := []bool{}

			for _, u := range tt.urls {
				got = append(got, a.AddURL(u))
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

type StringProcessor struct {
	results []string
	mu      *sync.Mutex
}

func (processor *StringProcessor) Process(url string, h map[string]string, d WebDriver, r ResponseHandler) error {
	processor.mu.Lock()
	processor.results = append(processor.results, url)
	processor.mu.Unlock()
	return nil
}

func findInSlice(sl []string, target string) bool {
	for i, end := 0, len(sl); i < end; i++ {
		if sl[i] == target {
			return true
		}
	}
	return false
}

func TestAnubis_Start(t *testing.T) {
	t.Run("Cancel will panic if start hasn't been called", func(t *testing.T) {
		a := NewTestAnubis()

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Anubis didn't panic.")
			}
		}()

		a.Cancel()
	})

	t.Run("Cancel func will stop execution", func(t *testing.T) {
		a := NewTestAnubis()

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Anubis panicked: %v", err)
			}
		}()

		a.Start()

		time.AfterFunc(1*time.Second, func() {
			a.Cancel()
		})

		a.Wait()
	})

	t.Run("Each URL pushed to anubis before start will be processed", func(t *testing.T) {
		a := NewTestAnubis()

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Anubis panicked: %v", err)
			}
		}()

		time.AfterFunc(2*time.Second, func() {
			a.Cancel()
		})

		a.AddURL("a")
		a.AddURL("b")
		a.AddURL("c")
		a.AddURL("d")

		a.Start()

		a.Wait()

		for _, u := range []string{"a", "b", "c", "d"} {
			if !findInSlice(a.processor.(*StringProcessor).results, u) {
				t.Errorf("URL not processed: %v", u)
			}
		}
	})

	t.Run("URLs added after start should be processed", func(t *testing.T) {
		a := NewTestAnubis()

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Anubis panicked: %v\n", err)
				debug.PrintStack()
			}
		}()

		a.Start()

		time.AfterFunc(1*time.Second, func() {
			a.AddURL("a")
			a.AddURL("b")
			a.AddURL("c")
			a.AddURL("d")
			a.Cancel()
		})

		a.Wait()

		for _, u := range []string{"a", "b", "c", "d"} {
			if !findInSlice(a.processor.(*StringProcessor).results, u) {
				t.Errorf("URL not processed: %v", u)
			}
		}

		for _, u := range []string{"e", "f"} {
			if findInSlice(a.processor.(*StringProcessor).results, u) {
				t.Errorf("URL processed: %v", u)
			}
		}
	})

	t.Run("URLs added after cancel should not be processed", func(t *testing.T) {
		a := NewTestAnubis()

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Anubis panicked: %v", err)
			}
		}()

		a.Start()
		a.AddURL("a")
		a.AddURL("b")
		a.AddURL("c")
		a.AddURL("d")

		time.AfterFunc(1*time.Second, func() {
			a.Cancel()
			a.AddURL("e")
			a.AddURL("f")
		})

		a.Wait()

		for _, u := range []string{"a", "b", "c", "d"} {
			if !findInSlice(a.processor.(*StringProcessor).results, u) {
				t.Errorf("URL not processed: %v", u)
			}
		}

		for _, u := range []string{"e", "f"} {
			if findInSlice(a.processor.(*StringProcessor).results, u) {
				t.Errorf("URL processed: %v", u)
			}
		}
	})
}

func TestAnubis_worker(t *testing.T) {
	t.Run("Worker should process URLs from queue until it is closed", func(t *testing.T) {
		a := NewTestAnubis()

		queue := make(chan string)
		processor := StringProcessor{mu: &sync.Mutex{}}

		a.wg.Add(1)
		go a.worker(&processor, queue)

		queue <- "a"
		queue <- "b"
		queue <- "c"
		close(queue)

		for _, u := range []string{"a", "b", "c"} {
			if !findInSlice(processor.results, u) {
				t.Errorf("URL not processed: %v", u)
			}
		}
	})
}
