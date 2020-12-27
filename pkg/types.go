package pkg

import "anubis/internal"

type Config internal.Config
type Pipeline []internal.PipelineFunc

type Anubis struct {
	Config *Config
	s      *internal.Supervisor
}
