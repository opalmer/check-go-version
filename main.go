package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opalmer/check-go-version/api"
)

func main() {
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
