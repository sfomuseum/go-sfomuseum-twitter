package twitter

import (
	"context"
	"path/filepath"
	"gocloud.dev/blob"
	"io"
)

func OpenTweets(ctx context.Context, tweets_uri string) (io.ReadCloser, error) {

	tweets_fname := filepath.Base(tweets_uri)
	tweets_root := filepath.Dir(tweets_uri)

	tweets_bucket, err := blob.OpenBucket(ctx, tweets_root)

	if err != nil {
		return nil, err
	}

	tweets_fh, err := tweets_bucket.NewReader(ctx, tweets_fname, nil)

	if err != nil {
		return nil, err
	}

	return tweets_fh, nil
}
