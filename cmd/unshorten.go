package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/sfomuseum/go-url-unshortener"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
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

	lookup := new(sync.Map)

	rsp.ForEach(func(_, tw gjson.Result) bool {

		go func(tw gjson.Result) {

			defer func() {
				done_ch <- true
			}()

			urls_rsp := tw.Get("entities.urls")

			if !urls_rsp.Exists() {
				err_ch <- errors.New("Missing URLs")
				return
			}

			to_fetch := make([]string, 0)

			for _, u := range urls_rsp.Array() {

				short_rsp := u.Get("expanded_url")
				short_url := short_rsp.String()

				_, ok := lookup.Load(short_url)

				if ok {
					return
				}

				to_fetch = append(to_fetch, short_url)
			}

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			for _, short_url := range to_fetch {

				url, err := unshortener.UnshortenString(ctx, cache, short_url)

				if err != nil {
					lookup.Store(short_url, "?")
					err_ch <- err
					continue
				}

				long_url := url.String()

				if short_url == long_url {
					long_url = "-"
				}

				lookup.Store(short_url, long_url)
			}

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

	report := make(map[string]string)

	lookup.Range(func(k interface{}, v interface{}) bool {
		shortened_url := k.(string)
		unshortened_url := v.(string)
		report[shortened_url] = unshortened_url
		return true
	})

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stdout)

	out := io.MultiWriter(writers...)

	enc := json.NewEncoder(out)
	err = enc.Encode(report)

	if err != nil {
		log.Fatal(err)
	}

}
