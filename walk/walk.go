package walk

import (
	"context"
	"encoding/json"
	"github.com/sfomuseum/go-sfomuseum-twitter/document"
	"io"
)

type WalkOptions struct {
	TweetChannel chan []byte
	ErrorChannel chan error
	DoneChannel  chan bool
	TrimPrefix bool
}

func WalkTweets(ctx context.Context, opts *WalkOptions, tweets_fh io.Reader) {

	defer func() {
		opts.DoneChannel <- true
	}()

	if opts.TrimPrefix {

		fh, err := document.TrimJavaScriptPrefix(ctx, tweets_fh)

		if err != nil {
			opts.ErrorChannel <- err
			return
		}

		tweets_fh = fh

	}
	
	type post struct {
		Tweet interface{} `json:"tweet"`
	}

	var posts []post

	dec := json.NewDecoder(tweets_fh)
	err := dec.Decode(&posts)

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
