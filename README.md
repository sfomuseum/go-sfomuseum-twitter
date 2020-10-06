# go-sfomuseum-twitter

Go package for working with Twitter archives.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/pointers cmd/pointers/main.go
go build -mod vendor -o bin/trim cmd/trim/main.go
go build -mod vendor -o bin/unshorten cmd/unshorten/main.go
```

### pointers

Export pointers (for users and hashtags) from a `tweets.json` file (produced by the `trim` tool).
	
```
$> ./bin/pointers -h
Usage of ./bin/pointers:
  -hashtags
    	Export hash tags in tweets. (default true)
  -mentions
    	Export users mentioned in tweets. (default true)
  -tweets string
    	The path your Twitter archive tweet.json file (produced by the sfomuseum/go-sfomuseum-twitter/cmd/trim tool, or equivalent)
```

For example:

```
$> ./bin/pointers -tweets /usr/local/data/tweet.json
property,value
user,genavnews
user,ladyfleur
user,JuanPDeAnda
user,sfo1977
user,787FirstClass
user,WNYCculture
user,Mairin_
user,ManuBarack
user,LMBernhard
user,OneBrownGirl
user,jtroll
user,spatrizi
user,patriciadenni20
...
tag,JimLund
tag,MoodLighting
tag,giants
tag,republicairlines
tag,T2
tag,eyecandy
tag,Wildenhain
```

### trim

Trim JavaScript boilerplate from a `tweets.js` file in order to make it valid JSON. Outputs to `STDOUT`.

```
$> ./bin/trim -h
Usage of ./bin/trim:
  -trim tweets.js
    	The leading string to remove from your tweets.js file. (default "window.YTD.tweet.part0 = ")
  -tweets tweets.js
    	The path to your tweets.js file.
```	

For example:

```
$> ./bin/trim -tweets /usr/local/twitter/data/tweet.js > /usr/local/twitter/data/tweet.json
```

### Unshorten

Expand all the URLs in a `tweets.js` file. Outputs a JSON dictionary to `STDOUT`.

```
$> ./bin/unshorten -h
Usage of ./bin/unshorten:
  -progress
    	Display progress information
  -qps int
    	Number of (unshortening) queries per second (default 10)
  -seed string
    	Pre-fill the unshortening cache with data in this file
  -timeout int
    	Maximum number of seconds of for an unshorterning request (default 30)
  -tweets string
    	The path your Twitter archive tweet.json file (produced by the sfomuseum/go-sfomuseum-twitter/cmd/trim tool, or equivalent)
```

For example:

```
$> ./bin/unshorten -progress /usr/local/twitter/data/tweet.json
2020/10/06 11:37:02 Head "http://gowal.la/p/hHG2": dial tcp: lookup gowal.la: no such host
2020/10/06 11:37:08 Head "http://danceonmarket.com/events/": dial tcp: lookup danceonmarket.com: no such host
2020/10/06 11:37:11 3431 of 3461 URLs left to unshorten (from 9759 tweets)
...and so on
```

## See also

* https://github.com/sfomuseum/go-url-unshortener