package pkg

import "anubis/internal"

// Returns a config object based on the current environment. Env variables are checked first, then args
func ParseEnv() *Config {
	c, err := internal.ParseArgs()
	if err != nil {
		// TODO helpful messages on config errors
		panic(err)
	}
	return (*Config)(c)
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

func (a *Anubis) Start() error {
	return a.s.Start()
}
