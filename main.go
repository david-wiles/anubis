package main

import (
	"anubis/pkg"
	"flag"
	"fmt"
	"os"
)

func main() {

	output := flag.String("output", ".", "The output directory. Note that if you are preserving only a single page, the full path to the file will be created.")
	proxy := flag.String("proxy", "", "Specifies the proxy to use during program execution")
	nWorkers := flag.Int("workers", 4, "Maximum number of concurrent requests")

	flag.Parse()

	a := anubis.NewAnubis(anubis.OutputOpt(*output), anubis.ProxyOpt(*proxy), anubis.NWorkerOpt(*nWorkers))

	startURLs := flag.Args()

	// Print error if no start URLs were provided
	if len(startURLs) == 0 {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s [ options... ] [ urls... ]\n", os.Args[0])
		flag.Usage()
		return
	}

	for _, url := range flag.Args() {
		a.AddURL(url)

		// Initialize start urls in the instance's handler
		a.Handler.(anubis.DefaultResponseHandler).NeededLinks[url] = true
	}

	a.Start()
	a.Wait()

	if err := a.Commit(); err != nil {
		panic(err)
	}
}
