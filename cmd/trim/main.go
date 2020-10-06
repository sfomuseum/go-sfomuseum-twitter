package main

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-sfomuseum-twitter"
	_ "gocloud.dev/blob/fileblob"
	"io"
	"log"
	"os"
)

func main() {

	tweets_uri := flag.String("tweets-uri", "", "A valid gocloud.dev/blob URI to your `tweets.js` file.")

	flag.Parse()

	ctx := context.Background()

	opts := &twitter.OpenTweetsOptions{
		TrimPrefix: true,
	}

	fh, err := twitter.OpenTweets(ctx, *tweets_uri, opts)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	_, err = io.Copy(os.Stdout, fh)

	if err != nil {
		log.Fatal(err)
	}
}
