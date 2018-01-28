package api

import (
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/blang/semver"
)

func getVersion(name string) (semver.Version, error) {
	matches := RegexSemanticVersion.FindAllStringSubmatch(name, 1)
	if matches == nil {
		return semver.Version{}, fmt.Errorf(
			`failed to retrieve version information from "%s"`, name)
	}
	version := matches[0][1]
	if strings.Count(version, ".") == 1 {
		version += ".0"
	}
	return semver.Make(version)
}

func getFullVersion(name string) (string, error) {
	matches := RegexFullVersion.FindAllStringSubmatch(name, 1)
	if matches == nil {
		return "", fmt.Errorf(
			`failed to retrieve full version information from "%s"`, name)
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
		return "", fmt.Errorf(`failed to retrieve architecture from "%s"`, name)
	}
}

func getPlatform(name string) (string, error) {
	switch strings.Count(name, "-") {
	case 1, 2:
		split := strings.Split(strings.Split(name, "-")[0], ".")
		return split[len(split)-1], nil
	default:
		return "", fmt.Errorf(`failed to retrieve platform from "%s"`, name)
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
