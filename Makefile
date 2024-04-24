PKG=github.com/wonwooseo/panawa
APP=batch
VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -o $(APP) -ldflags="-X $(PKG)/build.Version=$(VERSION) -X $(PKG)/build.BuildTime=$(shell date -Iseconds)" ./cmd

image:
	GOARCH=amd64 GOOS=linux go build -o $(APP)-deploy -ldflags="-X $(PKG)/build.Version=$(VERSION) -X $(PKG)/build.BuildTime=$(shell date -Iseconds) -w -s" ./cmd
	docker build -t wonwooseo/panawa:latest .
	docker scout cves wonwooseo/panawa:latest

clean:
	@-rm $(APP)

.PHONY: build
