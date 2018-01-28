package api

import (
	"fmt"
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
	BucketCache = false
	objects, err := GetBucketObjects()
	c.Assert(err, IsNil)
	s.objects = objects
}

func (s *VersionTest) SetUpTest(c *C) {
	BucketCache = false
}

func (s *VersionTest) TearDownTest(c *C) {
	url = ""
	BucketCache = true
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

func (s *VersionTest) TestGetVersionsMatchingPlatformMock(c *C) {
	versions := []*Version{
		{Platform: "foobar", Architecture: "foobar"},
		{Platform: runtime.GOOS, Architecture: "foobar"},
		{Platform: runtime.GOOS, Architecture: runtime.GOARCH},
	}

	c.Assert(
		FilterVersionsToPlatform(versions), DeepEquals,
		Versions{{Platform: runtime.GOOS, Architecture: runtime.GOARCH}})
}

func (s *VersionTest) TestFilterVersionsToPlatform(c *C) {
	versions, err := GetVersions()
	c.Assert(err, IsNil)

	for _, version := range FilterVersionsToPlatform(versions) {
		c.Assert(version.Platform, Equals, runtime.GOOS)
		c.Assert(version.Architecture, Equals, runtime.GOARCH)
	}
}

func (s *VersionTest) TestVersionString(c *C) {
	v := &Version{Name: "foo", Version: semver.Version{Major: 1, Minor: 2, Patch: 2}, Platform: "test", Architecture: "amd64"}
	c.Assert(v.String(), DeepEquals, "Version{Name: foo, Version: 1.2.2, Platform: test, Architecture: amd64}")
}

func (s *VersionTest) TestVersionsSort(c *C) {
	versions, err := GetVersions()
	c.Assert(err, IsNil)
	c.Assert(sort.IsSorted(versions), Equals, false)
	sort.Sort(versions)
	c.Assert(sort.IsSorted(versions), Equals, true)
}

func (s *VersionTest) TestGetReleaseVersions(c *C) {
	versions, err := GetReleaseVersions()
	c.Assert(err, IsNil)

	for _, version := range versions {
		c.Assert(version.Version.String(), Equals, version.FullVersion)
	}
}

func (s *VersionTest) TestGetReleaseVersionsForPlatform(c *C) {
	versions, err := GetReleaseVersionsForPlatform()
	c.Assert(err, IsNil)
	for _, version := range versions {
		c.Assert(version.Version.String(), Equals, version.FullVersion)
		c.Assert(version.Platform, Equals, runtime.GOOS)
		c.Assert(version.Architecture, Equals, runtime.GOARCH)
	}
}

func (s *VersionTest) TestGetRunningVersion(c *C) {
	v1, err := GetRunningVersion()
	c.Assert(err, IsNil)

	v2, err := getVersionFromName(
		fmt.Sprintf("%s.%s-%s", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	c.Assert(err, IsNil)
	c.Assert(v1, DeepEquals, v2)
}

func (s *VersionTest) TestGetLatestRelease(c *C) {
	versions, err := GetReleaseVersionsForPlatform()
	c.Assert(err, IsNil)
	sort.Sort(versions)
	expected := versions[len(versions)-1]
	version, err := GetLatestRelease()
	c.Assert(err, IsNil)
	c.Assert(expected, DeepEquals, version)
}

func (s *VersionTest) TestCheckLatest(c *C) {
	c.Assert(CheckLatest(&Version{Version: semver.MustParse("1.0.0")}, &Version{Version: semver.MustParse("1.9.3")}), Equals, false)
	c.Assert(CheckLatest(&Version{Version: semver.MustParse("1.9.2")}, &Version{Version: semver.MustParse("1.9.3")}), Equals, false)
	c.Assert(CheckLatest(&Version{Version: semver.MustParse("1.9.3")}, &Version{Version: semver.MustParse("1.9.3")}), Equals, true)
	c.Assert(CheckLatest(&Version{Version: semver.MustParse("1.9.4")}, &Version{Version: semver.MustParse("1.9.3")}), Equals, true)
	c.Assert(CheckLatest(&Version{Version: semver.MustParse("1.50.0")}, &Version{Version: semver.MustParse("1.9.3")}), Equals, true)
}
