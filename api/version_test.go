package api

import (
	"runtime"
	"sort"

	"cloud.google.com/go/storage"
	"github.com/blang/semver"
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

func (s *VersionTest) TestGetVersions(c *C) {
	versions, err := GetVersions()
	c.Assert(err, IsNil)
	c.Assert(len(versions), Not(Equals), 0)
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

func (s *VersionTest) TestVersionString(c *C) {
	v := &Version{Name: "foo", Version: semver.Version{Major: 1, Minor: 2, Patch: 2}, FullVersion: "1.2.3-foo", Architecture: "amd64"}
	c.Assert(v.String(), DeepEquals, "Version{foo, Version: 1.2.2, FullVersion: 1.2.3-foo, Architecture: amd64}")
}

func (s *VersionTest) TestVersionsSort(c *C) {
	found, err := GetVersions()
	c.Assert(err, IsNil)
	var versions Versions

	for _, version := range found {
		versions = append(versions, version)
	}
	c.Assert(sort.IsSorted(versions), Equals, false)
	sort.Sort(versions)
	c.Assert(sort.IsSorted(versions), Equals, true)
}

func (s *VersionTest) TestGetOfficialVersions(c *C) {
	versions, err := GetOfficialVersions()
	c.Assert(err, IsNil)

	for _, version := range versions {
		c.Assert(version.Version.String(), Equals, version.FullVersion)
	}
}
