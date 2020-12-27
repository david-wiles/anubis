package internal

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var UrlRegex = regexp.MustCompile("https?://[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b[-a-zA-Z0-9()@:%_+.~#?&/=]*")

func ParseUrl(u string) (*url.URL, error) {
	if UrlRegex.MatchString(u) {
		return url.Parse(u)
	}
	return nil, errors.New("Could not parse " + u + " as URL")
}

func ParseSiteMap(r io.Reader) ([]string, error) {
	urls := []string{}
	sitemap := struct {
		xml.Name `xml:"urlset"`
		Urls     []struct {
			Loc      string `xml:"loc"`
			Priority string `xml:"priority"`
			Lastmod  string `xml:"lastmod"`
		} `xml:"url"`
	}{}

	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&sitemap)
	if err != nil {
		return nil, err
	}

	for _, u := range sitemap.Urls {
		urls = append(urls, u.Loc)
	}

	return urls, nil
}

func SendRequest(client *http.Client, url string, config *Config) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//config.Auth.AddAuth(req)
	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
