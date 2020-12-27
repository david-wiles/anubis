package main

import (
	"anubis/pkg"
)

func main() {
	anubis := pkg.NewAnubis(pkg.ParseEnv()).WithPipeline(pkg.Pipeline{
		pkg.GetResources,
	})
	err := anubis.Start()
	if err != nil {
		panic(err)
	}
}
