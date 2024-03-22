PKG=github.com/wonwooseo/panawa
APP=batch
VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -o $(APP) -ldflags="-X $(PKG)/build.Version=$(VERSION) -X $(PKG)/build.BuildTime=$(shell date -Iseconds)" ./cmd

clean:
	@-rm $(APP)

.PHONY: build
