package main

import (
	"github.com/opalmer/check-go-version/checkversion"
)

func main() {
	releases, err := checkversion.GetReleases()
	if err != nil {
		panic(err)
	}

	for _, release := range releases {
		_ = release

	}

}
