cli:
	go build -mod vendor -o bin/pointers cmd/pointers/main.go
	go build -mod vendor -o bin/trim cmd/trim/main.go
	go build -mod vendor -o bin/unshorten cmd/unshorten/main.go
