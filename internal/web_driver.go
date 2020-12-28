package internal

import "net/http"

type WebDriver interface {
	SendRequest(string) (*http.Response, error)
	Authorize()
}

type DefaultWebDriver struct {
	client    *http.Client
	userAgent string
}

func (wd *DefaultWebDriver) SendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", wd.userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (wd *DefaultWebDriver) Authorize() {
	// No auth for default driver
}
