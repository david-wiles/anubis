package pkg

import (
	"anubis/internal"
	"net/url"
)

// Returns a config object based on the current environment. Env variables are checked first, then args
func FromEnv() *Config {
	c, err := internal.FromEnv()
	if err != nil {
		// TODO helpful messages on config errors
		panic(err)
	}
	return (*Config)(c)
}

type Anubis struct {
	Config *Config
	s      *internal.Supervisor
}

func NewAnubis(config *Config) *Anubis {
	return &Anubis{
		Config: config,
		s:      internal.NewSupervisor((*internal.Config)(config)),
	}
}

func (a *Anubis) WithPipeline(pipeline Pipeline) *Anubis {
	a.s.Pipeline = pipeline
	return a
}

func (a *Anubis) WithLinkJudge(judge internal.LinkJudge) *Anubis {
	a.s.ShouldAddLink = judge
	return a
}

func (a *Anubis) Start() error {
	return a.s.Start()
}

func (a *Anubis) AddLink(link *url.URL) {
	_ = a.s.Urls.QueueLinks(link)
}
