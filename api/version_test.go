package api

import (
	"cloud.google.com/go/storage"
	. "gopkg.in/check.v1"
)

type VersionTest struct{}

var _ = Suite(&VersionTest{})

type versionfacet struct {
	key      string
	expected *Version
	err      error
}

//func (s *VersionTest) TestReleaseString(c *C) {
//	v := Release{Key: "foo"}
//	c.Assert(v.String(), Not(Equals), "")
//}
//
//func (s *VersionTest) TestReleaseVersion(c *C) {
//	facets := []*versionfacet{
//		{"", nil, ErrFailedToMatchVersion},
//		{
//			"go1.7.2.src.tar.gz",
//			&Version{
//				Version:      "1.7.2",
//				Type:         "",
//				Platform:     "src",
//				Architecture: "tar",
//			},
//			nil,
//		},
//		{
//			"go1.3.3.windows-amd64.zip",
//			&Version{
//				Version:      "1.3.3",
//				Platform:     "windows",
//				Architecture: "amd64",
//			},
//			nil,
//		},
//		{
//			"go1.10beta1.linux-amd64.tar.gz",
//			&Version{
//				Version:      "1.10beta1",
//				Platform:     "linux",
//				Architecture: "amd64",
//			},
//			nil,
//		},
//		{
//			"go1.3.2.freebsd-386.tar.gz",
//			&Version{
//				Version:      "1.3.2",
//				Platform:     "freebsd",
//				Architecture: "386",
//			},
//			nil,
//		},
//		{
//			"go1.10beta1.linux-ppc64le.tar.gz.sha256",
//			&Version{
//				Version:      "1.10beta1",
//				Platform:     "linux",
//				Architecture: "ppc64le",
//			},
//			nil,
//		},
//	}
//
//	for _, input := range facets {
//		release := &Release{Key: input.key}
//		version, err := release.Version()
//		if input.err != nil {
//			c.Assert(err, ErrorMatches, err.Error())
//			c.Assert(version, IsNil)
//			continue
//		}
//
//		c.Assert(err, IsNil)
//		if version.Version != input.expected.Version {
//			c.Errorf(`%s "%s" != "%s"`, input.key, version.Version, input.expected.Version)
//			continue
//		}
//
//		if version.Platform != input.expected.Platform {
//			c.Errorf(`%s "%s" != "%s"`, input.key, version.Platform, input.expected.Platform)
//			continue
//		}
//
//		if version.Architecture != input.expected.Architecture {
//			c.Errorf(`%s "%s" != "%s"`, input.key, version.Architecture, input.expected.Architecture)
//			continue
//		}
//	}
//}

func (s *VersionTest) GetGetName(c *C) {
	c.Assert(
		getName(&storage.ObjectAttrs{Name: "foobar.tar"}), Equals, "foobar")
	c.Assert(
		getName(&storage.ObjectAttrs{Name: "foobar.tar.gz", ContentType: "application/x-gzip"}),
		Equals, "foobar")
}
