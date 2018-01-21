package main

import (
	"fmt"

	"github.com/opalmer/check-go-version/checkversion"
)

func main() {
	versions, err := checkversion.GetBucketVersions()
	if err != nil {
		panic(err)
	}

	for _, version := range versions {
		fmt.Println(version)
	}

}
