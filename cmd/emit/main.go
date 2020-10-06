package main

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-sfomuseum-twitter"
	"github.com/sfomuseum/go-sfomuseum-twitter/walk"
	"github.com/tidwall/pretty"
	_ "gocloud.dev/blob/fileblob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

func main() {

	tweets_uri := flag.String("tweets-uri", "", "A valid gocloud.dev/blob URI to your `tweets.js` file.")
	trim_prefix := flag.Bool("trim-prefix", true, "")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")
	as_json := flag.Bool("json", false, "Emit a JSON list.")
	format_json := flag.Bool("format-json", false, "Format JSON output for each record.")

	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	writers := make([]io.Writer, 0)

	if *to_stdout {
		writers = append(writers, os.Stdout)
	}

	if *to_devnull {
		writers = append(writers, ioutil.Discard)
	}

	if len(writers) == 0 {
		log.Fatal("Nothing to write to.")
	}

	wr := io.MultiWriter(writers...)

	open_opts := &twitter.OpenTweetsOptions{
		TrimPrefix: *trim_prefix,
	}

	tweets_fh, err := twitter.OpenTweets(ctx, *tweets_uri, open_opts)

	if err != nil {
		log.Fatalf("Failed to open %s, %v", *tweets_uri, err)
	}

	defer tweets_fh.Close()

	err_ch := make(chan error)
	tweet_ch := make(chan []byte)
	done_ch := make(chan bool)

	walk_opts := &walk.WalkOptions{
		DoneChannel:  done_ch,
		ErrorChannel: err_ch,
		TweetChannel: tweet_ch,
	}

	go walk.WalkTweets(ctx, walk_opts, tweets_fh)

	working := true
	count := uint32(0)

	if *as_json {
		wr.Write([]byte("["))
	}

	for {
		select {
		case <-done_ch:
			working = false
		case err := <-err_ch:
			log.Println(err)
			cancel()
		case body := <-tweet_ch:

			new_count := atomic.AddUint32(&count, 1)

			if new_count > 1 {

				if *as_json {
					wr.Write([]byte(","))
				}
			}

			if *as_json && *format_json {
				body = pretty.Pretty(body)
			}

			wr.Write(body)
			wr.Write([]byte("\n"))
		}

		if !working {
			break
		}
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
