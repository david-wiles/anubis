package internal

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Url       string        // --url, URL
	Sitemap   string        // --sitemap, SITEMAP
	Seed      []string      // --seed, SEED
	Output    string        // --out, OUT
	Stats     string        // --stats, STATS
	Proxy     string        // --proxy, PROXY
	Delay     time.Duration // --delay, DELAY
	UserAgent string        // --user-agent, USER_AGENT
	Auth      Auth
	Workers   int // --workers, WORKERS
}

func FromEnv() (c *Config, err error) {
	url := flag.String("url", "", "Url to preserve")
	sitemap := flag.String("sitemap", "", "Sitemap url to generate seeds")
	seed := flag.String("seed", "", "Optional comma-separated list of seed urls to start with")
	output := flag.String("out", "", "Location to put preserved files. If this is a url, then the program will "+
		"POST the files to the location specified. If empty, no files will be saved")
	stats := flag.String("stats", "", "If specified, stats about processed pages will be written here")
	proxy := flag.String("proxy", "", "If the program should run behind a proxy, then this is the url of the proxy")
	delay := flag.String("delay", "0", "Delay between successive requests. Not really necessary when using a proxy")
	userAgent := flag.String("user-agent", "", "Custom user-agent to use in requests")

	workers := flag.Int("workers", 1, "Number of workers making requests. Each will use their own delay, but pull from "+
		"the same queue of urls")

	flag.Parse()

	c = &Config{}

	// Assign variables from args, overwrite with environment variables if they exist
	c.Url = *url
	c.Sitemap = *sitemap

	if *seed != "" {
		c.Seed = strings.Split(*seed, ",")
	}

	c.Output = *output
	c.Stats = *stats
	c.Proxy = *proxy

	if dur, err := time.ParseDuration(*delay); err == nil {
		c.Delay = dur
	} else {
		return nil, errors.New("Could not parse delay string")
	}

	c.UserAgent = *userAgent
	c.Workers = *workers

	// Overwrite with environment variables
	if env := os.Getenv("URL"); env != "" {
		c.Url = env
	}

	if env := os.Getenv("SITEMAP"); env != "" {
		c.Sitemap = env
	}

	if env := os.Getenv("SEED"); env != "" {
		c.Seed = strings.Split(env, ",")
	}

	if env := os.Getenv("OUT"); env != "" {
		c.Output = env
	}

	if env := os.Getenv("STATS"); env != "" {
		c.Stats = env
	}

	if env := os.Getenv("PROXY"); env != "" {
		c.Proxy = env
	}

	if env := os.Getenv("DELAY"); env != "" {
		if d, err := time.ParseDuration(env); err == nil {
			c.Delay = d
		} else {
			return nil, errors.New("Could not parse delay string")
		}
	}

	if env := os.Getenv("WORKERS"); env != "" {
		if w, err := strconv.Atoi(env); err == nil && w > 0 {
			c.Workers = w
		} else {
			return nil, errors.New("Could not get number of workers")
		}
	}

	if c.Url != "" && c.Url[len(c.Url)-1] != '/' {
		c.Url += "/"
	}

	if c.Output[len(c.Output)-1] != '/' {
		c.Output += "/"
	}

	return c, nil
}
