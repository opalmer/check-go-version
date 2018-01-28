package api

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"strconv"

	"cloud.google.com/go/storage"
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
	Version string

	// VersionIn is `Version` converted to integers. Note, if 1.9 is the
	// version then the result of VersionInt is [3]int{1, 9, 0}. Golang's
	// versioning, historically speaking, has never had a .0 micro release.
	VersionInt [3]int

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
		`Version{"%s", Version: "%s", FullVersion: "%s", Architecture: "%s"}`,
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
	if vs[a].VersionInt[0] <= vs[b].VersionInt[0] {
		return true
	}

	if vs[a].VersionInt[1] <= vs[b].VersionInt[1] {
		return true
	}
	if vs[a].VersionInt[0] > vs[b].VersionInt[0] {
		return false
	}
	return false

	//return vs[a].VersionInt[2] >= vs[b].VersionInt[2]
}

func getVersion(name string) (string, error) {
	matches := RegexSemanticVersion.FindAllStringSubmatch(name, 1)
	if matches == nil {
		return "", fmt.Errorf(
			"failed to retrieve semantic version information from '%s'", name)
	}
	return matches[0][1], nil
}

func getFullVersion(name string) (string, error) {
	matches := RegexFullVersion.FindAllStringSubmatch(name, 1)
	if matches == nil {
		return "", fmt.Errorf(
			"failed to retrieve full version information from '%s'", name)
	}
	return matches[0][1], nil
}

func stripSuffix(name string) string {
	if strings.Contains(name, ".tar.gz") {
		return strings.TrimRight(name, ".tar.gz")
	}

	split := strings.Split(name, ".")
	return strings.TrimRight(name, "."+split[len(split)-1])
}

func getArchitecture(name string) (string, error) {
	switch strings.Count(name, "-") {
	case 1, 2:
		split := strings.Split(strings.Split(name, "-")[1], ".")
		return split[len(split)-1], nil
	default:
		return "", fmt.Errorf(`failed to extract architecture from "%s"`, name)
	}
}

func getPlatform(name string) (string, error) {
	switch strings.Count(name, "-") {
	case 1, 2:
		split := strings.Split(strings.Split(name, "-")[0], ".")
		return split[len(split)-1], nil
	default:
		return "", fmt.Errorf(`failed to extract platform from "%s"`, name)
	}
}

func skip(object *storage.ObjectAttrs) bool {
	switch object.ContentType {
	case "text/plain", "text/plain; charset=utf-8":
		return true
	}
	if strings.HasPrefix(object.Name, "getgo/") {
		return true
	}
	if strings.HasSuffix(object.Name, ".asc") {
		return true
	}
	if strings.Contains(object.Name, ".src.") {
		return true
	}
	return false
}

func versionToIntegers(version string) ([3]int, error) {
	output := [3]int{}
	for i, value := range strings.Split(version, ".") {
		converted, err := strconv.Atoi(value)
		if err != nil {
			return output, err
		}
		output[i] = converted
	}
	return output, nil
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

		versionInts, err := versionToIntegers(version)
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
			VersionInt:   versionInts,
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
