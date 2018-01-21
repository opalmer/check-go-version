package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

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

type LookupTest struct{}

var _ = Suite(&LookupTest{})

func (s *LookupTest) TearDownTest(c *C) {
	url = ""
}

func (s *LookupTest) TestGetBucketVersions(c *C) {
	ts := newTestServer(c, &items{Items: []*item{{Name: "testing"}}})
	defer ts.Close()
	objects, err := GetBucketObjects()
	c.Assert(err, IsNil)
	c.Assert(len(objects), Equals, 1)
	c.Assert(objects[0].Name, Equals, "testing")
}
