package main

import (
	"context"
	"github.com/sfomuseum/go-sfomuseum-twitter"	
	"github.com/sfomuseum/go-sfomuseum-twitter/walk"
	_ "gocloud.dev/blob/fileblob"
	"flag"
	"log"
)

func main() {

	tweets_uri := flag.String("tweets", "", "")
	trim_prefix := flag.Bool("trim-prefix", true, "")
	
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	tweets_fh, err := twitter.OpenTweets(ctx, *tweets_uri)

	if err != nil {
		log.Fatalf("Failed to open %s, %v", *tweets_uri, err)
	}

	defer tweets_fh.Close()
	
	err_ch := make(chan error)
	tweet_ch := make(chan []byte)
	done_ch := make(chan bool)

	walk_opts := &walk.WalkOptions{
		DoneChannel: done_ch,
		ErrorChannel: err_ch,
		TweetChannel: tweet_ch,
		TrimPrefix := *trim_prefix,
	}
	
	go walk.WalkTweets(ctx, walk_opts, tweets_fh)

	working := true
	
	for {
		select {
		case <- done_ch:
			working = false
		case err := <- err_ch:
			log.Println(err)
			cancel()
		case body := <- tweet_ch:

			log.Println(string(body))
		}
		
		if !working {
			break
		}
	}
	
}
