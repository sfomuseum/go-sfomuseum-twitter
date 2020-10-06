package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {

	dump_hashtags := flag.Bool("hashtags", true, "Export hash tags in tweets.")
	dump_mentions := flag.Bool("mentions", true, "Export users mentioned in tweets.")

	tweets := flag.String("tweets", "", "The path your Twitter archive tweet.json file (produced by the sfomuseum/go-sfomuseum-twitter/cmd/trim tool, or equivalent)")

	flag.Parse()

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

	mentions := new(sync.Map)
	hashtags := new(sync.Map)

	done_ch := make(chan bool)
	err_ch := make(chan error)

	rsp.ForEach(func(_, tw gjson.Result) bool {

		go func(tw gjson.Result) {

			defer func() {
				done_ch <- true
			}()

			// log.Println(tw.String())

			if *dump_mentions {
				mentions_rsp := tw.Get("tweet.entities.user_mentions")

				if !mentions_rsp.Exists() {
					err_ch <- errors.New("Missing mentions")
					return
				}

				for _, m := range mentions_rsp.Array() {
					user_rsp := m.Get("screen_name")
					user := user_rsp.String()
					mentions.Store(user, true)
				}
			}

			if *dump_hashtags {
				hashtags_rsp := tw.Get("tweet.entities.hashtags")

				if !hashtags_rsp.Exists() {
					err_ch <- errors.New("Missing hashtags")
					return
				}

				for _, h := range hashtags_rsp.Array() {
					tag_rsp := h.Get("text")
					tag := tag_rsp.String()
					hashtags.Store(tag, true)
				}
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

	writers := []io.Writer{
		os.Stdout,
	}

	multi := io.MultiWriter(writers...)
	wr := csv.NewWriter(multi)

	write_row := func(row []string) bool {

		err := wr.Write(row)

		if err != nil {
			log.Println(err)
			return false
		}

		return true
	}

	write_row([]string{"property", "value"})

	mentions.Range(func(key interface{}, value interface{}) bool {

		user := key.(string)

		out := []string{
			"user",
			user,
		}

		return write_row(out)
	})

	hashtags.Range(func(key interface{}, value interface{}) bool {

		tag := key.(string)

		out := []string{
			"tag",
			tag,
		}

		return write_row(out)
	})

}
