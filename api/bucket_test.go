package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"path/filepath"

	"cloud.google.com/go/storage"
	. "gopkg.in/check.v1"
)

type item struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
}

type items struct {
	Items []*item `json:"items"`
}

func newTestServer(c *C, i *items) *httptest.Server {
	data, err := json.Marshal(i)
	c.Assert(err, IsNil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/b/golang/o?alt=json&delimiter=&pageToken=&prefix=&projection=full&versions=false":
			_, err := w.Write(data)
			c.Assert(err, IsNil)

		default:
			c.Fatalf("unhandled url %s", r.URL)
		}
	}))
	url = ts.URL
	return ts
}

type BucketTest struct{}

var _ = Suite(&BucketTest{})

func (s *BucketTest) SetUpTest(c *C) {
	BucketCache = false
}

func (s *BucketTest) TearDownTest(c *C) {
	url = ""
	BucketCache = true
	BucketCacheFile = ""
}

func (s *BucketTest) TestGetBucketVersions(c *C) {
	ts := newTestServer(c, &items{Items: []*item{{Name: "testing"}}})
	defer ts.Close()
	objects, err := GetBucketObjects()
	c.Assert(err, IsNil)
	c.Assert(len(objects), Equals, 1)
	c.Assert(objects[0].Name, Equals, "testing")
}

func (s *BucketTest) TestGetBucketVersionsCached(c *C) {
	BucketCache = true
	BucketCacheFile = filepath.Join(c.MkDir(), "foo", "cache.json")
	expected, err := GetBucketObjects()
	c.Assert(err, IsNil)
	objects, err := readCache(BucketCacheFile)
	c.Assert(err, IsNil)
	c.Assert(expected, DeepEquals, objects)
}

func (s *BucketTest) TestCachePathCustomFile(c *C) {
	expected := filepath.Join(c.MkDir(), "foo", "cache.json")
	BucketCacheFile = expected
	path, err := cachepath()
	c.Assert(err, IsNil)
	c.Assert(path, Equals, expected)
}

func (s *BucketTest) TestCachePathDefaultPath(c *C) {
	BucketCacheFile = ""
	path, err := cachepath()
	c.Assert(err, IsNil)
	usr, err := user.Current()
	c.Assert(err, IsNil)
	c.Assert(path, Equals, filepath.Join(usr.HomeDir, ".cache", "check-go-version", "bucket.json"))
}

func (s *BucketTest) TestWriteCache(c *C) {
	BucketCacheFile = filepath.Join(c.MkDir(), "foo", "cache.json")
	objects := []*storage.ObjectAttrs{{Name: "1"}, {Name: "2"}}
	writeCache(objects)
	content, err := ioutil.ReadFile(BucketCacheFile)
	c.Assert(err, IsNil)
	expected, err := json.Marshal(cache{Objects: objects})
	c.Assert(err, IsNil)
	c.Assert(content, DeepEquals, expected)
}

func (s *BucketTest) TestReadCacheInvalid(c *C) {
	BucketCacheFile = filepath.Join(c.MkDir(), "foo", "cache.json")
	c.Assert(os.MkdirAll(filepath.Dir(BucketCacheFile), 0700), IsNil)
	c.Assert(ioutil.WriteFile(BucketCacheFile, []byte("{"), 0600), IsNil)
	_, err := readCache(BucketCacheFile)
	c.Assert(err, NotNil)
	_, stat := os.Stat(BucketCacheFile)
	c.Assert(os.IsNotExist(stat), Equals, true)
}

func (s *BucketTest) TestHasExpiredMissingFile(c *C) {
	BucketCacheFile = filepath.Join(c.MkDir(), "foo", "cache.json")
	expired, err := hasExpired(BucketCacheFile)
	c.Assert(err, IsNil)
	c.Assert(expired, Equals, true)
}

