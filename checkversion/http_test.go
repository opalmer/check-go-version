package checkversion

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

const (
	// getGoOnly is used to ensure that entries containing getgo/* are filtered
	// by default.
	getGoOnly = `
<ListBucketResult xmlns="http://doc.s3.amazonaws.com/2006-03-01">
<Name>golang</Name>
<Prefix/>
<Marker/>
<NextMarker>go1.8.1.linux-amd64.tar.gz.sha256</NextMarker>
<IsTruncated>true</IsTruncated>
<Contents>
<Key>getgo/installer.exe</Key>
<Generation>1501264926677062</Generation>
<MetaGeneration>1</MetaGeneration>
<LastModified>2017-07-28T18:02:06.630Z</LastModified>
<ETag>"919c6c41be11fbf8a27405e42d033d18"</ETag>
<Size>5131776</Size>
</Contents>
<Contents>
<Key>getgo/installer_darwin</Key>
<Generation>1501264926674843</Generation>
<MetaGeneration>1</MetaGeneration>
<LastModified>2017-07-28T18:02:06.610Z</LastModified>
<ETag>"923241b8b23403d5d5c762afadc9844a"</ETag>
<Size>5157928</Size>
</Contents>
<Contents>
<Key>getgo/installer_linux</Key>
<Generation>1501264926685855</Generation>
<MetaGeneration>1</MetaGeneration>
<LastModified>2017-07-28T18:02:06.632Z</LastModified>
<ETag>"a38297226e0736950f9ee7b96f76ce81"</ETag>
<Size>5179246</Size>
</Contents>
</ListBucketResult>
`
)

type HTTPTest struct {
	url         string
	filtergetgo bool
}

var _ = Suite(&HTTPTest{url: APIUrl, filtergetgo: FilterGetGoVersions})

func (s *HTTPTest) TearDownTest(c *C) {
	APIUrl = s.url
	FilterGetGoVersions = s.filtergetgo
}

func (s *HTTPTest) TestGetReleases(c *C) {
	releases, err := GetReleases()
	c.Assert(err, IsNil)
	if len(releases) == 0 {
		c.Fatal()
	}
}

func (s *HTTPTest) TestGetReleasesVersion(c *C) {
	releases, err := GetReleases()
	c.Assert(err, IsNil)

	for _, release := range releases {
		if _, err := release.Version(); err != nil {
			c.Errorf("Version failed %s", release)
		}
	}
}

func (s *HTTPTest) TestGetReleasesFilterGetGoTrue(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, getGoOnly)
	}))
	defer ts.Close()
	APIUrl = ts.URL
	FilterGetGoVersions = true

	releases, err := GetReleases()
	c.Assert(err, IsNil)
	c.Assert(len(releases), Equals, 0)
}

func (s *HTTPTest) TestGetReleasesFilterGetGoFalse(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, getGoOnly)
	}))
	defer ts.Close()
	APIUrl = ts.URL
	FilterGetGoVersions = false

	releases, err := GetReleases()
	c.Assert(err, IsNil)
	c.Assert(len(releases), Equals, 3)
}
