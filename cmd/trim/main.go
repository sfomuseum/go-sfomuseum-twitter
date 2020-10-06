package main

import (
	"flag"
	"io"
	"log"
	"os"
)

func main() {

	tweets := flag.String("tweets", "", "The path to your `tweets.js` file.")
	trim := flag.String("trim", "window.YTD.tweet.part0 = ", "The leading string to remove from your `tweets.js` file.")

	flag.Parse()

	fh, err := os.Open(*tweets)

	if err != nil {
		log.Fatal(err)
	}

	offset := int64(len(*trim))
	whence := 0

	_, err = fh.Seek(offset, whence)

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, fh)

	if err != nil {
		log.Fatal(err)
	}
}
