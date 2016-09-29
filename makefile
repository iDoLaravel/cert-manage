.PHONY: build test

vet:
	go tool vet .
linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/cert-manage-linux .
osx:
	GOOS=darwin GOARCH=386 go build -o bin/cert-manage-osx .
win:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/cert-manage-win .

build: vet osx linux win

test: build
	go test -v ./...
