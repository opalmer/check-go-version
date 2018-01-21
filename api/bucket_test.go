package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

type LookupTest struct{}

var _ = Suite(&LookupTest{})

func (s *LookupTest) TearDownTest(c *C) {
	url = ""
}

func (s *LookupTest) TestGetBucketVersions(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/b/golang/o?alt=json&delimiter=&pageToken=&prefix=&projection=full&versions=false":
			fmt.Fprint(w, `{"items": [{"name": "testing"}]}`)

		default:
			c.Fatalf("unhandled url %s", r.URL)
		}
	}))
	defer ts.Close()
	url = ts.URL
	objects, err := GetBucketObjects()
	c.Assert(err, IsNil)
	c.Assert(len(objects), Equals, 1)
	c.Assert(objects[0].Name, Equals, "testing")
}
