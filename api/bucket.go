package api

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	// BucketTimeout represents the amount of time we're willing to spend
	// retrieving information from the golang bucket.
	BucketTimeout = time.Minute * 5

	// url is used for testing.
	url = ""
)

// GetBucketObjects queries the golang bucket in Google Object Store and
// returns the versions present as a list.
func GetBucketObjects() ([]*storage.ObjectAttrs, error) {
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
	defer client.Close()

	bucket := client.Bucket("golang")
	objects := bucket.Objects(ctx, nil)
	var entries []*storage.ObjectAttrs
	for {
		object, err := objects.Next()

		switch err {
		case nil:
			entries = append(entries, object)
		case iterator.Done:
			return entries, nil
		default:
			return entries, err
		}
	}
}
