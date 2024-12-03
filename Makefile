GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/emit cmd/emit/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/pointers cmd/pointers/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/trim cmd/trim/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/unshorten cmd/unshorten/main.go
