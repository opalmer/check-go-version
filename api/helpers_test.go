package api

import (
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	. "gopkg.in/check.v1"
)

var _ = Suite(&HelpersTest{})

type HelpersTest struct {
	objects []*storage.ObjectAttrs
}

func (s *HelpersTest) SetUpSuite(c *C) {
	objects, err := GetBucketObjects()
	c.Assert(err, IsNil)
	s.objects = objects
}

func (s *HelpersTest) TearDownTest(c *C) {
	url = ""
}

func (s *HelpersTest) TestSkip(c *C) {
	c.Assert(skip(&storage.ObjectAttrs{ContentType: "text/plain"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{ContentType: "text/plain; charset=utf-8"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{Name: "getgo/"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{Name: "foobar.asc"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{Name: "foo.src.tar.gz"}), Equals, true)
}

func (s *HelpersTest) TestStripSuffixSynthetic(c *C) {
	facets := []*strfacet{
		{"go1.9.2rc2.windows-386.zip", "go1.9.2rc2.windows-386"},
		{"go1.9.2rc2.windows-386.msi", "go1.9.2rc2.windows-386"},
		{"go1.9.freebsd-amd64.tar", "go1.9.freebsd-amd64"},
		{"go1.9.freebsd-amd64.tar.gz", "go1.9.freebsd-amd64"},
	}
	for _, facet := range facets {
		out := stripSuffix(facet.in)
		if out != facet.out {
			c.Fatalf(`"%s" -> "%s" != "%s"`, facet.in, out, facet.out)
		}
	}
}

func (s *HelpersTest) TestStripSuffix(c *C) {
	knownSuffixes := []string{".tar.gz", ".tar", ".asc", "msi", ".zip", "pkg"}
	for _, object := range s.objects {
		if skip(object) {
			continue
		}
		name := stripSuffix(object.Name)
		for _, suffix := range knownSuffixes {
			if strings.HasSuffix(name, suffix) {
				c.Fatalf("Failed to strip suffix from %s", object.Name)
			}
		}
	}
}

func (s *HelpersTest) TestGetPlatformError(c *C) {
	_, err := getPlatform("go1.4rc2.darwin-amd64-osx10.8-1")
	c.Assert(err, ErrorMatches, `failed to retrieve platform from "go1.4rc2.darwin-amd64-osx10.8-1"`)
	_, err = getPlatform("go1.4rc2.darwin")
	c.Assert(err, ErrorMatches, `failed to retrieve platform from "go1.4rc2.darwin"`)
}

func (s *HelpersTest) TestGetPlatformSynthetic(c *C) {
	facets := []*strfacet{
		{"go1.9.2rc2.windows-386", "windows"},
		{"go1.9.2rc2.windows-amd64", "windows"},
		{"go1.9.freebsd-386", "freebsd"},
		{"go1.9.freebsd-amd64", "freebsd"},
		{"go1.4rc2.darwin-amd64-osx10.8", "darwin"},
	}
	for _, facet := range facets {
		out, err := getPlatform(facet.in)
		c.Assert(err, IsNil)
		if out != facet.out {
			c.Fatalf(`"%s" -> "%s" != "%s"`, facet.in, out, facet.out)
		}
	}
}

func (s *HelpersTest) TestGetPlatform(c *C) {
	for _, object := range s.objects {
		if skip(object) {
			continue
		}
		_, err := getPlatform(stripSuffix(object.Name))
		c.Assert(err, IsNil)
	}
}

func (s *HelpersTest) TestGetArchitectureSynthetic(c *C) {
	facets := []*strfacet{
		{"go1.9.2rc2.windows-386", "386"},
		{"go1.9.2rc2.windows-amd64", "amd64"},
		{"go1.9.freebsd-386", "386"},
		{"go1.9.freebsd-amd64", "amd64"},
		{"go1.4rc2.darwin-amd64-osx10.8", "amd64"},
	}
	for _, facet := range facets {
		out, err := getArchitecture(facet.in)
		c.Assert(err, IsNil)
		if out != facet.out {
			c.Fatalf(`"%s" -> "%s" != "%s"`, facet.in, out, facet.out)
		}
	}
}
func (s *HelpersTest) TestGetArchitectureError(c *C) {
	for _, input := range []string{"", "1-2-3-4"} {
		_, err := getArchitecture(input)
		c.Assert(err, ErrorMatches, fmt.Sprintf(`failed to retrieve architecture from "%s"`, input))
	}
}

func (s *HelpersTest) TestGetVersionSynthetic(c *C) {
	facets := []*strfacet{
		{"go1.9.2rc2.windows-386", "1.9.2"},
		{"go1.9.2rc2.windows-amd64", "1.9.2"},
		{"go1.9.freebsd-386", "1.9.0"},
		{"go1.9.freebsd-amd64", "1.9.0"},
		{"go1.4rc2.darwin-amd64-osx10.8", "1.4.0"},
	}
	for _, facet := range facets {
		out, err := getVersion(facet.in)
		c.Assert(err, IsNil)
		if out.String() != facet.out {
			c.Fatalf(`"%s" -> "%s" != "%s"`, facet.in, out, facet.out)
		}
	}
}

func (s *HelpersTest) TestGetVersionError(c *C) {
	for _, input := range []string{"x.", "", "1", "1."} {
		_, err := getVersion(input)
		c.Assert(err, ErrorMatches, fmt.Sprintf(`failed to retrieve version information from "%s"`, input))
	}
}
func (s *HelpersTest) TestGetFullVersionSynthetic(c *C) {
	facets := []*strfacet{
		{"go1.9.2rc2.windows-386", "1.9.2rc2"},
		{"go1.9.2rc2.windows-amd64", "1.9.2rc2"},
		{"go1.9.freebsd-386", "1.9"},
		{"go1.9.freebsd-amd64", "1.9"},
		{"go1.4rc2.darwin-amd64-osx10.8", "1.4rc2"},
	}
	for _, facet := range facets {
		out, err := getFullVersion(facet.in)
		c.Assert(err, IsNil)
		if out != facet.out {
			c.Fatalf(`"%s" -> "%s" != "%s"`, facet.in, out, facet.out)
		}
	}
}

func (s *HelpersTest) TestGetFullVersionError(c *C) {
	for _, input := range []string{"x.", "", "1", "1."} {
		_, err := getFullVersion(input)
		c.Assert(err, ErrorMatches, fmt.Sprintf(`failed to retrieve full version information from "%s"`, input))
	}
}
