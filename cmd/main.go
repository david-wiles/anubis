package main

import (
	"anubis/pkg"
)

func main() {
	anubis := pkg.NewAnubis(pkg.FromEnv()).
		WithPipeline(pkg.Pipeline{pkg.FollowLinks}).
		WithLinkJudge(pkg.LocalLinkFilter)

	err := anubis.Start()
	if err != nil {
		panic(err)
	}
}
