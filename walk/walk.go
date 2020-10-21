package walk

import (
	"context"
	"encoding/json"
	"io"
	"sync"
)

type WalkOptions struct {
	TweetChannel chan []byte
	ErrorChannel chan error
	DoneChannel  chan bool
}

type WalkTweetsCallbackFunc func(ctx context.Context, tweet []byte) error

func WalkTweetsWithCallback(ctx context.Context, tweets_fh io.Reader, cb WalkTweetsCallbackFunc) error {

	err_ch := make(chan error)
	tweet_ch := make(chan []byte)
	done_ch := make(chan bool)

	walk_opts := &WalkOptions{
		DoneChannel:  done_ch,
		ErrorChannel: err_ch,
		TweetChannel: tweet_ch,
	}

	go WalkTweets(ctx, walk_opts, tweets_fh)

	working := true
	wg := new(sync.WaitGroup)

	for {
		select {
		case <-done_ch:
			working = false
		case err := <-err_ch:
			return err
		case body := <-tweet_ch:

			wg.Add(1)

			go func(body []byte) {

				defer wg.Done()

				err := cb(ctx, body)

				if err != nil {
					err_ch <- err
				}

			}(body)

		}

		if !working {
			break
		}
	}

	wg.Wait()
	return nil
}

func WalkTweets(ctx context.Context, opts *WalkOptions, tweets_fh io.Reader) {

	defer func() {
		opts.DoneChannel <- true
	}()

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
