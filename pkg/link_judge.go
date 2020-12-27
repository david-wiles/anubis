package pkg

import "anubis/internal"

func LocalLinkFilter(link internal.DiscoveredLink) bool {
	return link.Current.Hostname() == link.Url.Hostname()
}
