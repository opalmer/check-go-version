package api

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/blang/semver"
)

var (
	// RegexSemanticVersion extracts version information for a posted
	// release. This does not include the qualifier such as 'alpha', 'beta',
	// etc.
	RegexSemanticVersion = regexp.MustCompile(`^go(\d+\.\d+(\.\d+|)).+$`)

	// RegexFullVersion extracts the full version information. Unlike
	// RegexSemanticVersion this will include the qualifiers such as
	// 'alpha', 'beta', etc.
	RegexFullVersion = regexp.MustCompile(`^go(\d+\.\d+[a-z0-9]*(\.\d[a-z0-9]*|)).+$`)
)

// Version returns specific information about a Go version.
type Version struct {
	// Name is the complete name of the version excluding the original
	// suffix (ex. go1.9.2rc2.windows-386 instead of go1.9.2rc2.windows-386.msi)
	Name string

	// Version is the version excluding qualifiers such as 'beta1'. For
	// example if FullVersion is 'go1.4rc2` then FullVersion is '1.4'.
	Version semver.Version

	// FullVersion contains the complete version information include the semantic
	// version.
	FullVersion string

	// Platform is the generic platform (ex. windows, darwin, linux)
	Platform string

	// Architecture is the architecture of the binary. (ex. amd64, 386)
	Architecture string
}

// String returns a human readable string representing this struct.
func (v *Version) String() string {
	return fmt.Sprintf(
		`Version{%v, Version: %v, FullVersion: %v, Architecture: %v}`,
		v.Name, v.Version, v.FullVersion, v.Architecture)
}

// Versions is a list of Version structs that have the added benefit of
// being sortable.
type Versions []*Version

func (vs Versions) Len() int {
	return len(vs)
}

func (vs Versions) Swap(a int, b int) {
	vs[a], vs[b] = vs[b], vs[a]
}

func (vs Versions) Less(a int, b int) bool {
	return vs[a].Version.LT(vs[b].Version)
}

// GetVersions returns a list of golang releases. This function will ignore
// any object returned from GetBucketObjects() that matches one or more
// of the following conditions:
//   * Name of object starts with getgo/ - Ignored because this tool assumes
//     you already have go installed. These objects also don't contain any
//     version information.
//   * Content type of objects starts with text/plain - Ignored because they
//     are not actually releases (typically sha256 sums).
//   * Is a signature instead of a release.
//   * Appears to be a source rather than a release.
func GetVersions() ([]*Version, error) {
	objects, err := GetBucketObjects()
	if err != nil {
		return nil, err
	}

	unique := map[string]bool{}
	var versions []*Version

	for _, object := range objects {
		if skip(object) {
			continue
		}

		name := stripSuffix(object.Name)
		if _, ok := unique[name]; ok {
			continue
		}
		unique[name] = true

		platform, err := getPlatform(name)
		if err != nil {
			return nil, err
		}

		version, err := getVersion(name)
		if err != nil {
			return nil, err
		}

		fullVersion, err := getFullVersion(name)
		if err != nil {
			return nil, err
		}

		architecture, err := getArchitecture(name)
		if err != nil {
			return nil, err
		}

		versions = append(versions, &Version{
			Name:         name,
			Platform:     platform,
			Version:      version,
			FullVersion:  fullVersion,
			Architecture: architecture,
		})
	}

	return versions, nil
}

// GetVersionsMatchingPlatform acts as a filter and returns all versions
// that are matching the current platform.
func GetVersionsMatchingPlatform(versions []*Version) []*Version {
	var output []*Version

	for _, version := range versions {
		if version.Platform != runtime.GOOS {
			continue
		}
		if version.Architecture != runtime.GOARCH {
			continue
		}
		output = append(output, version)
	}

	return output
}

// GetOfficialVersions calls GetVersions and removes any version that
// are alpha/beta/test/etc.
func GetOfficialVersions() ([]*Version, error) {
	versions, err := GetVersions()
	if err != nil {
		return nil, err
	}
	var output []*Version
	for _, version := range versions {
		if version.Version.String() != version.FullVersion {
			continue
		}
		output = append(output, version)
	}
	return output, nil
}

// GetLatestVersionMatchingPlatform will return the latest release matching
// the current operating system and platform. Note, this excludes any release
// that's alpha/beta/etc.
func GetLatestVersionMatchingPlatform() (*Version, error) {
	return nil, nil
}
