package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/opalmer/check-go-version/api"
)

func main() {
	versions, err := api.GetVersions()
	if err != nil {
		log.Fatal(err)
	}

	candidates := api.Versions{}
	for _, version := range api.FilterVersionsToPlatform(versions) {
		// Make sure we skip any version where the version
		// does not equal the full version. This happens when
		// there's a qualifier in the version such as 'alpha'
		// or 'beta'.
		if version.Version.String() != version.FullVersion {
			continue
		}
		candidates = append(candidates, version)
	}

	sort.Sort(candidates)
	latest := candidates[len(candidates)-1]
	fmt.Println(latest)
	//runt

	//for _, value := range candidates {
	//	fmt.Println(value)
	//}
}
