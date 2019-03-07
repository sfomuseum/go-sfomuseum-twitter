CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/sfomuseum/go-sfomuseum-twitter
	# cp *.go src/github.com/sfomuseum/go-url-unshortener/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/sfomuseum/go-url-unshortener"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	# go fmt *.go
	go fmt cmd/*.go

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/pointers cmd/pointers.go
	@GOPATH=$(GOPATH) go build -o bin/trim cmd/trim.go
	@GOPATH=$(GOPATH) go build -o bin/unshorten cmd/unshorten.go
