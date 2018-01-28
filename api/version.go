package api

import (
	"fmt"
	"regexp"
	"runtime"
	"sort"

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
		`Version{Name: %v, Version: %v, Platform: %v, Architecture: %v}`,
		v.Name, v.Version, v.Platform, v.Architecture)
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
func GetVersions() (Versions, error) {
	objects, err := GetBucketObjects()
	if err != nil {
		return nil, err
	}

	unique := map[string]bool{}
	var versions Versions

	for _, object := range objects {
		if skip(object) {
			continue
		}

		name := stripSuffix(object.Name)
		if _, ok := unique[name]; ok {
			continue
		}
		unique[name] = true

		version, err := getVersionFromName(name)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// GetReleaseVersions calls GetVersions and removes any version that is a
// 'non-release' version. Basically this returns all versions except those that
// are alpha/beta/etc. This does not provide any platform specific filtering.
func GetReleaseVersions() ([]*Version, error) {
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

// GetReleaseVersionsForPlatform calls GetReleaseVersions() and filters the
// results so only releases matching the current platform are returned.
func GetReleaseVersionsForPlatform() (Versions, error) {
	versions, err := GetReleaseVersions()
	return FilterVersionsToPlatform(versions), err
}

// FilterVersionsToPlatform acts as a filter and returns all versions
// that are matching the current platform.
func FilterVersionsToPlatform(versions Versions) Versions {
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

// GetRunningVersion constructs and returns a *Version struct for the currently
// running instance of Go.
func GetRunningVersion() (*Version, error) {
	return getVersionFromName(
		fmt.Sprintf("%s.%s-%s", runtime.Version(), runtime.GOOS, runtime.GOARCH))
}

// GetLatestRelease returns a *Version struct matching the latest release
// for the currently running platform.
func GetLatestRelease() (*Version, error) {
	versions, err := GetReleaseVersionsForPlatform()
	if err != nil {
		return nil, err
	}
	sort.Sort(versions)
	return versions[len(versions)-1], nil
}

// CheckLatest will return true if the latest version is the same as the
// currently running version.
func CheckLatest(running *Version, latest *Version) bool {
	return running.Version.GTE(latest.Version)
}
