package anubis

import (
	"errors"
	"log"
	"net/url"
	"path"
	"regexp"
)

var (
	URLRE    = regexp.MustCompile("https?://[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b[-a-zA-Z0-9()@:%_+.~#?&/=]*")
	LinkRE   = regexp.MustCompile("<(link|LINK)(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+href=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*/>")
	ScriptRE = regexp.MustCompile("<(script|SCIRPT)(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+src=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*(/>|>)")
	ImageRE  = regexp.MustCompile("<img(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+src=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*/?>")
)

func getFullURL(parent string, link string) (string, error) {
	if !URLRE.MatchString(parent) {
		return "", errors.New("Could not parse " + parent + " as URL")
	}

	u, err := url.Parse(parent)
	if err != nil {
		return "", err
	}

	if !URLRE.MatchString(link) {
		link = u.Scheme + "://" + path.Join(u.Hostname(), link)
	}

	if !URLRE.MatchString(link) {
		return "", errors.New("Could not parse " + link + " as URL")
	}

	return link, nil
}

func GetLinkURLs(parent string, html string) []string {
	urls := []string{}

	matches := LinkRE.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			if parsedUrl, err := getFullURL(parent, match[4]); err == nil {
				urls = append(urls, parsedUrl)
			} else {
				log.Println(err)
			}
		}
	}

	return urls
}

func GetScriptURLs(parent string, html string) []string {
	urls := []string{}

	matches := ScriptRE.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			if parsedUrl, err := getFullURL(parent, match[4]); err == nil {
				urls = append(urls, parsedUrl)
			} else {
				log.Println(err)
			}
		}
	}

	return urls
}

func GetImageURLs(parent string, html string) []string {
	urls := []string{}

	matches := ImageRE.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			if parsedUrl, err := getFullURL(parent, match[3]); err == nil {
				urls = append(urls, parsedUrl)
			} else {
				log.Println(err)
			}
		}
	}

	return urls
}
