package internal

import "net/http"

type Auth interface {
	AddAuth(r *http.Request)
}

type BasicAuth struct {
	username string
	password string
}

func (auth BasicAuth) AddAuth(r *http.Request) {
	r.SetBasicAuth(auth.username, auth.password)
}

type NoAuth struct{}

func (auth NoAuth) AddAuth(r *http.Request) {}
