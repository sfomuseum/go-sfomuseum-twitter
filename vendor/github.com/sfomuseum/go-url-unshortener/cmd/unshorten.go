package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	_ "fmt"
	"github.com/sfomuseum/go-url-unshortener"
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

	progress := flag.Bool("progress", false, "Display progress information")
	verbose := flag.Bool("verbose", false, "Be chatty about what's going on")
	stdin := flag.Bool("stdin", false, "Read URLs from STDIN")
	qps := flag.Int("qps", 10, "Number of (unshortening) queries per second")
	to := flag.Int("timeout", 30, "Maximum number of seconds of for an unshorterning request")
	seed_file := flag.String("seed", "", "Pre-fill the unshortening cache with data in this file")

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

	remaining := 0

	type UnshortenResponse struct {
		ShortenedURL   string
		UnshortenedURL string
		Error          error
	}

	done_ch := make(chan bool)
	rsp_ch := make(chan *UnshortenResponse)

	unshorten := func(ctx context.Context, str_url string) {

		defer func() {
			done_ch <- true
		}()

		u, err := unshortener.UnshortenString(ctx, cache, str_url)

		var rsp UnshortenResponse

		if err != nil {

			rsp = UnshortenResponse{
				ShortenedURL: str_url,
				Error:        err,
			}

			rsp_ch <- &rsp
		}

		if u != nil {

			rsp = UnshortenResponse{
				ShortenedURL:   str_url,
				UnshortenedURL: u.String(),
			}

			rsp_ch <- &rsp
		}

		// assume that ctx.Done() has been invoked
	}

	if *stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			remaining += 1
			str_url := scanner.Text()
			go unshorten(ctx, str_url)
		}

	} else {

		for _, str_url := range flag.Args() {
			remaining += 1
			go unshorten(ctx, str_url)
		}
	}

	total := remaining

	completed_ch := make(chan bool)

	if *progress {

		go func() {

			for {
				select {
				case <-completed_ch:
					break
				case <-time.After(10 * time.Second):
					log.Printf("%d of %d URLs left to unshorten\n", remaining, total)
				}
			}
		}()
	}

	lookup := new(sync.Map)

	for remaining > 0 {

		select {
		case <-done_ch:
			remaining -= 1
		case rsp := <-rsp_ch:

			if rsp.Error != nil {

				lookup.Store(rsp.ShortenedURL, "?")
				log.Printf("Failed to unshorted '%s' %s", rsp.ShortenedURL, rsp.Error)

			} else {

				if rsp.ShortenedURL == rsp.UnshortenedURL {
					lookup.Store(rsp.ShortenedURL, "-")
				} else {
					lookup.Store(rsp.ShortenedURL, rsp.UnshortenedURL)
				}

				if *verbose {
					log.Printf("%s becomes %s\n", rsp.ShortenedURL, rsp.UnshortenedURL)
				}
			}

		default:
			// log.Println(remaining)
		}
	}

	completed_ch <- true

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
