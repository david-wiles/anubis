package internal

import (
	"flag"
	"strings"
	"time"
)

type Config struct {
	Url       string
	Sitemap   string
	Seed      []string
	Output    string
	Proxy     string
	Delay     time.Duration
	UserAgent string
	Auth      Auth
	Workers   int
}

func ParseArgs() (*Config, error) {
	url := flag.String("url", "", "Url to preserve")
	sitemap := flag.String("sitemap", "", "Sitemap url to generate seeds")
	seed := flag.String("seed", "", "Optional comma-separated list of seed urls to start with")
	output := flag.String("out", "./", "Location to put preserved files. If this is a url, then the program will "+
		"POST the files to the location specified")
	proxy := flag.String("proxy", "", "If the program should run behind a proxy, then this is the url of the proxy")
	delay := flag.String("delay", "", "Delay between successive requests. Not really necessary when using a proxy")
	userAgent := flag.String("user-agent", "", "Custom user-agent to use in requests")

	authType := flag.String("auth", "", "Type of authentication to use, if necessary")

	username := flag.String("username", "", "Username for basic auth")
	password := flag.String("password", "", "Password for basic auth")

	workers := flag.Int("workers", 1, "Number of workers making requests. Each will use their own delay, but pull from "+
		"the same queue of urls")

	flag.Parse()

	var (
		normalizedUrl                  = *url
		normalizedOutput               = *output
		seedUrls                       = []string{}
		auth             Auth          = nil
		delayTime        time.Duration = 0
		err              error
	)

	if normalizedUrl != "" && normalizedUrl[len(normalizedUrl)-1] != '/' {
		normalizedUrl += "/"
	}

	if normalizedOutput[len(normalizedOutput)-1] != '/' {
		normalizedOutput += "/"
	}

	if *seed != "" {
		seedUrls = strings.Split(*seed, ",")
	}

	if *delay != "" {
		delayTime, err = time.ParseDuration(*delay)
		if err != nil {
			return nil, err
		}
	}

	switch *authType {
	case "basic":
		auth = &BasicAuth{
			username: *username,
			password: *password,
		}
	default:
		auth = &NoAuth{}
	}

	return &Config{
		Url:       normalizedUrl,
		Sitemap:   *sitemap,
		Seed:      seedUrls,
		Auth:      auth,
		Output:    *output,
		Proxy:     *proxy,
		Delay:     delayTime,
		UserAgent: *userAgent,
		Workers:   *workers,
	}, nil
}
