# go-sfomuseum-twitter

Go package for working with Twitter archives.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/emit cmd/emit/main.go
go build -mod vendor -o bin/pointers cmd/pointers/main.go
go build -mod vendor -o bin/trim cmd/trim/main.go
go build -mod vendor -o bin/unshorten cmd/unshorten/main.go
```

### emit

```
$> ./bin/emit -h
Usage of ./bin/emit:
  -format-json
    	Format JSON output for each record.
  -json
    	Emit a JSON list.
  -null
    	Emit to /dev/null
  -stdout
    	Emit to STDOUT (default true)
  -trim-prefix
    	Trim default tweet.js JavaScript prefix. (default true)
  -tweets-uri tweets.js
    	A valid gocloud.dev/blob URI to your tweets.js file.
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
  -trim-prefix
    	Trim default tweet.js JavaScript prefix. (default true)
  -tweets-uri tweets.js
    	A valid gocloud.dev/blob URI to your tweets.js file.
```

For example:

```
$> ./bin/pointers -tweets-uri file:///usr/local/twitter/data/tweet.js
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
./bin/trim -h
Usage of ./bin/trim:
  -tweets-uri tweets.js
    	A valid gocloud.dev/blob URI to your tweets.js file.
```	

For example:

```
$> ./bin/trim -tweets /usr/local/twitter/data/tweet.js > /usr/local/twitter/data/tweet.json
```

### Unshorten

Expand all the URLs in a `tweets.js` file. Outputs a JSON dictionary to `STDOUT`.

```
./bin/unshorten -h
Usage of ./bin/unshorten:
  -progress
    	Display progress information
  -qps int
    	Number of (unshortening) queries per second (default 10)
  -seed string
    	Pre-fill the unshortening cache with data in this file
  -timeout int
    	Maximum number of seconds of for an unshorterning request (default 30)
  -trim-prefix
    	Trim default tweet.js JavaScript prefix. (default true)
  -tweets-uri tweets.js
    	A valid gocloud.dev/blob URI to your tweets.js file.
```

For example:

```
$> ./bin/unshorten -progress -tweets-uri file:///usr/local/twitter/data/tweet.js
2020/10/06 11:37:02 Head "http://gowal.la/p/hHG2": dial tcp: lookup gowal.la: no such host
2020/10/06 11:37:08 Head "http://danceonmarket.com/events/": dial tcp: lookup danceonmarket.com: no such host
2020/10/06 11:37:11 3431 of 3461 URLs left to unshorten (from 9759 tweets)
...and so on
```

## See also

* https://github.com/sfomuseum/go-url-unshortener