package pkg

import (
	"anubis/internal"
	"regexp"
)

var (
	anchorRe = regexp.MustCompile("<[aA](\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+href=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*>")
	linkRe   = regexp.MustCompile("<(link|LINK)(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+href=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*/>")
	scriptRe = regexp.MustCompile("<(script|SCIRPT)(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+src=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*(/>|>)")
	imageRe  = regexp.MustCompile("<img(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s+src=\"([-a-zA-Z0-9()@:%_+.~#?&/=]+)\"(\\s+[a-zA-Z0-9]+(=\"[-a-zA-Z0-9 .]+\")?)*?\\s*/?>")
)

// Finds all links and resources on the page and adds them to the queue
func FollowLinks(bytes []byte, curr string, found chan internal.FoundUrl) ([]byte, error) {
	html := string(bytes)

	matches := anchorRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			found <- internal.FoundUrl{curr, match[3]}
		}
	}

	matches = linkRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			found <- internal.FoundUrl{curr, match[4]}
		}
	}

	matches = scriptRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 8 {
			found <- internal.FoundUrl{curr, match[4]}
		}
	}

	matches = imageRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			found <- internal.FoundUrl{curr, match[3]}
		}
	}

	return bytes, nil
}

// Finds all links and resources on the page but only follows links if they are to the same host
func FollowLocalLinks(bytes []byte, curr string, found chan internal.FoundUrl) ([]byte, error) {
	html := string(bytes)
	currUrl, err := internal.ParseUrl(curr)
	if err != nil {
		return bytes, err
	}

	matches := anchorRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			m := match[3]
			parsedMatch, err := internal.ParseUrl(m)
			if err != nil {
				// We couldn't parse the url, so it is a relative link
				found <- internal.FoundUrl{curr, m}
			} else if parsedMatch.Hostname() == currUrl.Hostname() {
				found <- internal.FoundUrl{curr, m}
			}
		}
	}

	matches = linkRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			found <- internal.FoundUrl{curr, match[4]}
		}
	}

	matches = scriptRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 8 {
			found <- internal.FoundUrl{curr, match[4]}
		}
	}

	matches = imageRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			found <- internal.FoundUrl{curr, match[3]}
		}
	}

	return bytes, nil
}

// Adds stylesheet, javascript, and image resources to the queue
func GetResources(bytes []byte, curr string, found chan internal.FoundUrl) ([]byte, error) {
	html := string(bytes)

	matches := linkRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			found <- internal.FoundUrl{curr, match[4]}
		}
	}

	matches = scriptRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 8 {
			found <- internal.FoundUrl{curr, match[4]}
		}
	}

	matches = imageRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			found <- internal.FoundUrl{curr, match[3]}
		}
	}

	return bytes, nil
}
