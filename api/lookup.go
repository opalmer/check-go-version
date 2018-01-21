package api

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// BucketTimeout represents the amount of time we're willing to spend
// retrieving information from the golang bucket.
var BucketTimeout = time.Minute * 5

// GetBucketVersions queries the golang bucket in Google Object Store and
// returns the versions present as a list.
func GetBucketVersions() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), BucketTimeout)
	defer cancel()

	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bucket := client.Bucket("golang")
	objects := bucket.Objects(ctx, nil)
	var releases []string
	for {
		object, err := objects.Next()

		switch err {
		case nil:
			releases = append(releases, object.Name)
		case iterator.Done:
			return releases, nil
		default:
			return releases, err
		}
	}
}
