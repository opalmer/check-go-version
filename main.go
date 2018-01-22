package main

import (
	"log"
	//"sort"
	"fmt"

	"github.com/opalmer/check-go-version/api"
)

func main() {
	versions, err := api.GetVersions()
	if err != nil {
		log.Fatal(err)
	}

	candidates := api.Versions{}
	for _, version := range api.GetVersionsMatchingPlatform(versions) {
		// Make sure we skip any version where the version
		// does not equal the full version. This happens when
		// there's a qualifier in the version such as 'alpha'
		// or 'beta'.
		if version.Version != version.FullVersion {
			continue
		}
		candidates = append(candidates, version)
	}

	//sort.Sort(candidates)

	for _, value := range candidates {
		fmt.Println(value)
	}
}
