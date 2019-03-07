package main

import (
	"flag"
	"io"
	"log"
	"os"
)

func main() {

	tweets := flag.String("tweets", "", "...")
	trim := flag.String("trim", "window.YTD.tweet.part0 = ", "...")

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
