package checkversion

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GetBucketVersions queries the golang bucket in Google Object Store and returns
// the versions as a list.
func GetBucketVersions() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
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
