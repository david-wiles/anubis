package internal

import (
	"net/url"
	"path"
)

type LinkJudge func(DiscoveredLink) bool

type DiscoveredLink struct {
	Current *url.URL
	Url     *url.URL
}

type DiscoveredChan chan DiscoveredLink

func (f DiscoveredChan) SendLink(curr string, link string) error {
	cUrl, err := ParseUrl(curr)
	if err != nil {
		return err
	}

	if !UrlRegex.MatchString(link) {
		link = cUrl.Scheme + "://" + path.Join(cUrl.Hostname(), link)
	}

	lUrl, err := ParseUrl(link)
	if err != nil {
		return err
	}

	f <- DiscoveredLink{cUrl, lUrl}
	return nil
}
