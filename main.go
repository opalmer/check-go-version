package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/opalmer/check-go-version/api"
)

var (
	nocache = flag.Bool(
		"no-cache", false,
		"Disables caching of responses from Google object store.")
)

func main() {
	flag.Parse()
	if *nocache {
		api.BucketCache = false
	}

	running, err := api.GetRunningVersion()
	if err != nil {
		log.Fatal(err)
	}

	latest, err := api.GetLatestRelease()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" latest: %s\n", latest)
	fmt.Printf("running: %s\n", running)
	if !api.CheckLatest(running, latest) {
		os.Exit(1)
	}
}
