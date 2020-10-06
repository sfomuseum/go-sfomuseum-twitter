package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"github.com/sfomuseum/go-sfomuseum-twitter/document"
)

func main() {

	tweets := flag.String("tweets", "", "The path to your `tweets.js` file.")

	flag.Parse()

	ctx := context.Background()
	
	fh, err := os.Open(*tweets)

	if err != nil {
		log.Fatal(err)
	}

	trimmed, err := document.TrimJavaScriptPrefix(ctx, fh)

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, trimmed)

	if err != nil {
		log.Fatal(err)
	}
}
