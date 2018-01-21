package api

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	// RegexVersion extracts version information for a posted release.
	RegexVersion = regexp.MustCompile(`^go(\d+\.\d+(\.\d+|)).+$`)

	// ErrFailedToMatchVersion is returned by Release.Version if
	// RegexVersion does not produce any results.
	ErrFailedToMatchVersion = errors.New(
		"failed to retrieve version information")
)

// Version returns specific information about a Go version from
// a Release struct.
type Version struct {
	Version      string
	Type         string
	Platform     string
	Architecture string
}

// String returns a human readable string representing this struct.
func (v *Version) String() string {
	return fmt.Sprintf(
		`Version{"%s", Platform: "%s", Architecture: "%s"}`,
		v.Version, v.Platform, v.Architecture)
}

// Release encapsulates information about a single version of Go from
// the remote API. T
type Release struct {
	Key          string    `xml:"Key"`
	Generation   int       `xml:"Generation"`
	LastModified time.Time `xml:"LastModified"`
	Size         int       `xml:"Size"`
	SHA256       string
}

func (r *Release) version() (string, error) {
	matches := RegexVersion.FindAllStringSubmatch(r.Key, -1)
	if matches == nil {
		return "", ErrFailedToMatchVersion
	}
	return matches[0][1], nil
}

// String returns a human readable string representing this struct.
func (r *Release) String() string {
	return fmt.Sprintf(
		`Release{Key: "%s", Modified: "%s"}`, r.Key, r.LastModified)
}

// Version parses and returns the version of Go from the Key field.
func (r *Release) Version() (*Version, error) {
	split := strings.Split(r.Key, "-")
	matches := RegexVersion.FindAllStringSubmatch(r.Key, 1)
	if matches == nil {
		return nil, ErrFailedToMatchVersion
	}
	version := matches[0][1]
	fmt.Println(version)
	if strings.Contains(r.Key, ".src.") {
		matches := RegexVersion.FindAllStringSubmatch(r.Key, -1)
		if matches == nil {
			return nil, ErrFailedToMatchVersion
		}
		return &Version{}, nil
	}

	if len(split) == 1 {

	}
	//fmt.Println(version)

	return nil, nil
	//matches := RegexVersion.FindAllStringSubmatch(r.Key, -1)
	//if matches == nil {
	//	return nil, ErrFailedToMatchVersion
	//}
	//
	//strings.Split(r.Key, "-")
	//
	//version, err := r.version()
	//if err != nil {
	//	return nil, err
	//}
	//_ = version
	//return nil, nil
	//return &Version{
	//	Version:      version,
	//	Type:         match[2],
	//	Platform:     match[3],
	//	Architecture: match[4],
	//}, nil
}
