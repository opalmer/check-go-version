package checkversion

import (
	"encoding/xml"
	"net/http"
	"strings"
)

var (
	// APIUrl is the url returning information about Go versions.
	APIUrl = "https://storage.googleapis.com/golang/"

	// FilterGetGoVersions will cause GetReleaes to filter out
	// any release with Key field starting with getgo/*.
	FilterGetGoVersions = true
)

// ListBucketResult is used to parse the top level
type ListBucketResult struct {
	IsTruncated bool       `xml:"IsTruncated"`
	Contents    []*Release `xml:"Contents"`
}

// GetReleases will call the remote API and return information
// all versions present in the response.
func GetReleases() ([]*Release, error) {
	response, err := http.Get(APIUrl)
	if err != nil {
		return nil, err
	}

	marshaller := xml.NewDecoder(response.Body)
	result := &ListBucketResult{}

	if err := marshaller.Decode(result); err != nil {
		return nil, err
	}

	if !FilterGetGoVersions {
		return result.Contents, nil
	}

	var releases []*Release
	for _, release := range result.Contents {
		if strings.HasPrefix(release.Key, "getgo/") {
			continue
		}
		releases = append(releases, release)
	}

	return releases, nil
}
