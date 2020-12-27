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

// Finds all links and resources on the page
func FollowLinks(bytes []byte, curr string, found internal.DiscoveredChan) ([]byte, error) {
	html := string(bytes)

	matches := anchorRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			if err := found.SendLink(curr, match[3]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	matches = linkRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			if err := found.SendLink(curr, match[4]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	matches = scriptRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 8 {
			if err := found.SendLink(curr, match[4]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	matches = imageRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			if err := found.SendLink(curr, match[3]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	return bytes, nil
}

// Adds stylesheet, javascript, and image resources to the queue
func GetResources(bytes []byte, curr string, found internal.DiscoveredChan) ([]byte, error) {
	html := string(bytes)

	matches := linkRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 7 {
			if err := found.SendLink(curr, match[4]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	matches = scriptRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 8 {
			if err := found.SendLink(curr, match[4]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	matches = imageRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 6 {
			if err := found.SendLink(curr, match[3]); err != nil {
				internal.Log.LogError(err)
			}
		}
	}

	return bytes, nil
}
