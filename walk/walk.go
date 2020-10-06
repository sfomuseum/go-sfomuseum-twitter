package walk

import (
	"context"
	"encoding/json"
	"fmt"
	"gocloud.dev/blob"
	"path/filepath"
)

type WalkCallbackFunc func(context.Context, []byte) error

type WalkOptions struct {
	TweetChannel chan []byte
	ErrorChannel chan error
	DoneChannel  chan bool
	Callback     WalkCallbackFunc
}

func WalkTweets(ctx context.Context, opts *WalkOptions, tweets_uri string) {

	defer func() {
		opts.DoneChannel <- true
	}()

	tweets_fname := filepath.Base(tweets_uri)
	tweets_root := filepath.Dir(tweets_uri)

	tweets_bucket, err := blob.OpenBucket(ctx, tweets_root)

	if err != nil {
		opts.ErrorChannel <- fmt.Errorf("Failed to open %s, %v", tweets_root, err)
		return
	}

	tweets_fh, err := tweets_bucket.NewReader(ctx, tweets_fname, nil)

	if err != nil {
		opts.ErrorChannel <- fmt.Errorf("Failed to open %s, %v", tweets_fname, err)
		return
	}

	defer tweets_fh.Close()

	// Add hooks to trim leading JS stuff here (see also: cmd/trim)

	type post struct {
		Tweet interface{} `json:"tweet"`
	}

	var posts []post

	dec := json.NewDecoder(tweets_fh)
	err = dec.Decode(&posts)

	if err != nil {
		opts.ErrorChannel <- err
		return
	}

	for _, p := range posts {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		tw_body, err := json.Marshal(p.Tweet)

		if err != nil {
			opts.ErrorChannel <- err
			continue
		}

		opts.TweetChannel <- tw_body
	}

}
