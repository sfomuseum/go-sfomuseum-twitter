package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/sfomuseum/go-url-unshortener"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	qps := flag.Int("qps", 10, "Number of (unshortening) queries per second")
	to := flag.Int("timeout", 30, "Maximum number of seconds of for an unshorterning request")
	seed_file := flag.String("seed", "", "Pre-fill the unshortening cache with data in this file")

	tweets := flag.String("tweets", "", "...")

	flag.Parse()

	rate := time.Second / time.Duration(*qps)
	timeout := time.Second * time.Duration(*to)

	worker, err := unshortener.NewThrottledUnshortener(rate, timeout)

	if err != nil {
		log.Fatal(err)
	}

	seed := make(map[string]string)

	if *seed_file != "" {

		fh, err := os.Open(*seed_file)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(body, &seed)

		if err != nil {
			log.Fatal(err)
		}
	}

	cache, err := unshortener.NewCachedUnshortenerWithSeed(worker, seed)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal_ch := make(chan os.Signal)
	signal.Notify(signal_ch, os.Interrupt, syscall.SIGTERM)

	go func(c chan os.Signal) {
		<-c
		cancel()
		os.Exit(0)
	}(signal_ch)

	fh, err := os.Open(*tweets)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		log.Fatal(err)
	}

	rsp := gjson.ParseBytes(body)

	if !rsp.Exists() {
		log.Fatal("Nothing to export")
	}

	count := len(rsp.Array())
	remaining := count

	done_ch := make(chan bool)
	err_ch := make(chan error)

	rsp.ForEach(func(_, tw gjson.Result) bool {

		go func(tw gjson.Result) {

			defer func() {
				done_ch <- true
			}()

			short_url := "PLEASE WRITE ME"

			long_url, err := unshortener.UnshortenString(ctx, cache, short_url)

			if err != nil {
				err_ch <- err
				return
			}

			log.Println(short_url, long_url)

		}(tw)

		return true
	})

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
			// log.Printf("%d of %d remaining\n", remaining, count)
		case err := <-err_ch:
			log.Println(err)
		default:
			// pass
		}
	}

}
