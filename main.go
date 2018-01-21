package main

import (
	//"fmt"
	"log"

	"fmt"

	"github.com/opalmer/check-go-version/api"
)

func main() {
	versions, err := api.GetVersions()
	if err != nil {
		log.Fatal(err)
	}

	for _, version := range versions {
		fmt.Println(version)
	}

}
