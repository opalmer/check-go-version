package api

import (
	"net/http/httptest"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
	. "gopkg.in/check.v1"
)

var _ = Suite(&VersionTest{})

type VersionTest struct {
	objects []*storage.ObjectAttrs
}

type strfacet struct {
	in  string
	out string
}

func (s *VersionTest) SetUpSuite(c *C) {
	objects, err := GetBucketObjects()
	c.Assert(err, IsNil)
	s.objects = objects
}

func (s *VersionTest) TearDownTest(c *C) {
	url = ""
}

func (s *VersionTest) newServer(c *C) *httptest.Server {
	items := &items{}
	for _, object := range s.objects {
		items.Items = append(items.Items, &item{
			Name:        object.Name,
			ContentType: object.ContentType,
		})
	}
	return newTestServer(c, items)
}

func (s *VersionTest) TestGetVersions(c *C) {
	versions, err := GetVersions()
	c.Assert(err, IsNil)
	c.Assert(len(versions), Not(Equals), 0)
}

func (s *VersionTest) TestSkip(c *C) {
	c.Assert(skip(&storage.ObjectAttrs{ContentType: "text/plain"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{ContentType: "text/plain; charset=utf-8"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{Name: "getgo/"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{Name: "foobar.asc"}), Equals, true)
	c.Assert(skip(&storage.ObjectAttrs{Name: "foo.src.tar.gz"}), Equals, true)
}

func (s *VersionTest) TestStripSuffixSynthetic(c *C) {
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

func (s *VersionTest) TestStripSuffix(c *C) {
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

func (s *VersionTest) TestExtractPlatformError(c *C) {
	_, err := getPlatform("go1.4rc2.darwin-amd64-osx10.8-1")
	c.Assert(err, ErrorMatches, `failed to extract platform from "go1.4rc2.darwin-amd64-osx10.8-1"`)
	_, err = getPlatform("go1.4rc2.darwin")
	c.Assert(err, ErrorMatches, `failed to extract platform from "go1.4rc2.darwin"`)
}

func (s *VersionTest) TestExtractPlatformSynthetic(c *C) {
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

func (s *VersionTest) TestExtractPlatform(c *C) {
	for _, object := range s.objects {
		if skip(object) {
			continue
		}
		_, err := getPlatform(stripSuffix(object.Name))
		c.Assert(err, IsNil)
	}
}

func (s *VersionTest) TestExtractArchitectureSynthetic(c *C) {
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

func (s *VersionTest) TestGetSemanticVersionSynthetic(c *C) {
	facets := []*strfacet{
		{"go1.9.2rc2.windows-386", "1.9.2"},
		{"go1.9.2rc2.windows-amd64", "1.9.2"},
		{"go1.9.freebsd-386", "1.9"},
		{"go1.9.freebsd-amd64", "1.9"},
		{"go1.4rc2.darwin-amd64-osx10.8", "1.4"},
	}
	for _, facet := range facets {
		out, err := getVersion(facet.in)
		c.Assert(err, IsNil)
		if out != facet.out {
			c.Fatalf(`"%s" -> "%s" != "%s"`, facet.in, out, facet.out)
		}
	}
}

func (s *VersionTest) TestGerFullVersionSynthetic(c *C) {
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

func (s *VersionTest) TestVersionToIntegers(c *C) {
	v, err := versionToIntegers("1.0.0")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, [3]int{1, 0, 0})

	v, err = versionToIntegers("1.1.1")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, [3]int{1, 1, 1})

	v, err = versionToIntegers("1.1")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, [3]int{1, 1, 0})
}

func (s *VersionTest) TestGetVersionsIgnoresPlainText(c *C) {
	ts := newTestServer(c, &items{Items: []*item{
		{Name: "1", ContentType: "text/plain"},
		{Name: "2", ContentType: "text/plain"},
		{Name: "3", ContentType: "text/plain; charset=utf-8"},
		{Name: "4", ContentType: "text/plain; charset=utf-8"},
	}})
	defer ts.Close()
	versions, err := GetVersions()
	c.Assert(err, IsNil)
	c.Assert(len(versions), Equals, 0)
}

func (s *VersionTest) TestGetVersionsIgnoresGetGoInstallers(c *C) {
	ts := newTestServer(c, &items{Items: []*item{
		{Name: "getgo/1"},
		{Name: "getgo/2"},
	}})
	defer ts.Close()
	versions, err := GetVersions()
	c.Assert(err, IsNil)
	c.Assert(len(versions), Equals, 0)
}

func (s *VersionTest) TestGetVersionsMatchingPlatform(c *C) {
	versions := []*Version{
		{Platform: "foobar", Architecture: "foobar"},
		{Platform: runtime.GOOS, Architecture: "foobar"},
		{Platform: runtime.GOOS, Architecture: runtime.GOARCH},
	}

	c.Assert(
		GetVersionsMatchingPlatform(versions), DeepEquals,
		[]*Version{{Platform: runtime.GOOS, Architecture: runtime.GOARCH}})
}
