package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	// BucketCache enables and disables caching of the bucket's
	// response. This ensures that if check-go-version is run
	// multiple times by a build we're not hitting Google's
	// API for each invocation.
	BucketCache = true

	// BucketCacheTime is the amount of time to keep bucket information
	// cached.
	BucketCacheTime = time.Minute * 30

	// BucketCacheFile is the location to cache bucket information. By default
	// this will be ~/.cache/check-go-version/bucket.json unless it's
	// overridden.
	BucketCacheFile = ""

	// BucketTimeout represents the amount of time we're willing to spend
	// retrieving information from the golang bucket.
	BucketTimeout = time.Minute * 5

	url = ""
)

type cache struct {
	Objects []*storage.ObjectAttrs `json:"objects"`
}

func cachepath() (string, error) {
	if BucketCacheFile != "" {
		return BucketCacheFile, os.MkdirAll(filepath.Dir(BucketCacheFile), 0700)
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	path := filepath.Join(usr.HomeDir, ".cache", "check-go-version", "bucket.json")
	return path, os.MkdirAll(filepath.Dir(path), 0700)
}

func hasExpired(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	now := time.Now()
	modified := stat.ModTime()
	expires := modified.Add(BucketCacheTime)
	return now.After(expires), nil
}

func readCache(path string) ([]*storage.ObjectAttrs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	cache := &cache{}

	// Decode the data but discard the cache if we can't.
	if err := decoder.Decode(cache); err != nil {
		os.Remove(path) // nolint: errcheck
		return nil, err
	}
	return cache.Objects, nil
}

// writes the objects to the cache. We ignore any failures here because
// we can always reach out to get the data again without the cache.
func writeCache(objects []*storage.ObjectAttrs) {
	if path, err := cachepath(); err == nil {
		cache := &cache{Objects: objects}
		if data, err := json.Marshal(cache); err == nil {
			ioutil.WriteFile(path, data, 0600) // nolint: errcheck
		}
	}
}

// GetBucketObjects queries the golang bucket in Google Object Store and
// returns the versions present as a list.
func GetBucketObjects() ([]*storage.ObjectAttrs, error) { // nolint: gocyclo
	if BucketCache {
		path, err := cachepath()
		if err != nil {
			return nil, err
		}

		expired, err := hasExpired(path)
		if err != nil {
			return nil, err
		}
		if !expired {
			if cached, err := readCache(path); err == nil {
				return cached, nil
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), BucketTimeout)
	defer cancel()

	options := []option.ClientOption{option.WithoutAuthentication()}
	if url != "" {
		options = append(options, option.WithEndpoint(url))
	}

	client, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}
	defer client.Close() // nolint: errcheck

	bucket := client.Bucket("golang")
	objects := bucket.Objects(ctx, nil)
	var entries []*storage.ObjectAttrs
	for {
		object, err := objects.Next()

		switch err {
		case nil:
			entries = append(entries, object)
		case iterator.Done:
			if BucketCache {
				writeCache(entries)
			}
			return entries, nil
		default:
			return entries, err
		}
	}
}
