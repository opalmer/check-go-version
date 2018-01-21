package main

import (
	"fmt"

	"github.com/opalmer/check-go-version/api"
)

func main() {
	versions, err := api.GetBucketVersions()
	if err != nil {
		panic(err)
	}

	for _, version := range versions {
		fmt.Println(version)
	}

}
