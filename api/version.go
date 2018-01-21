package api

import (
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
)

var (
	// RegexSemanticVersion extracts version information for a posted release.
	RegexSemanticVersion = regexp.MustCompile(`^go(\d+\.\d+(\.\d+|)).+$`)
)

// Version returns specific information about a Go version from
// a Release struct.
type Version struct {
	// Name is the complete name of the version excluding the original
	// suffix (ex. go1.9.2rc2.windows-386 instead of go1.9.2rc2.windows-386.msi)
	Name string

	// SemanticVersion is the version excluding qualifiers such as
	// 'beta1'. For example if Version is 'go1.4rc2` then Version
	// becomes '1.4'.
	SemanticVersion string

	// Version contains the complete version information include the semantic
	// version.
	Version string

	// Platform is the generic platform (ex. windows, darwin, linux)
	Platform string

	// PlatformVersion is the version of the specific
	// platform (ex. osx10.8). Note, this is not set on all
	// platforms. Only those that have had specific releases published.
	PlatformVersion string

	// Architecture is the architecture of the binary. (ex. amd64, 386)
	Architecture string
}

// String returns a human readable string representing this struct.
func (v *Version) String() string {
	return fmt.Sprintf(
		`Version{"%s", SemanticVersion: "%s", Architecture: "%s"}`,
		v.Name, v.SemanticVersion, v.Architecture)
}

func getSemanticVersion(name string) (string, error) {
	matches := RegexSemanticVersion.FindAllStringSubmatch(name, 1)
	if matches == nil {
		return "", fmt.Errorf(
			"failed to retrieve version information from '%s'", name)
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

		architecture, err := getArchitecture(name)
		if err != nil {
			return nil, err
		}

		versions = append(versions, &Version{
			Name:         name,
			Platform:     platform,
			Architecture: architecture,
		})
	}

	return versions, nil
}
