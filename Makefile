VERSION = "master"

build:
	go build -ldflags="-s -w -X example/version.Version=${VERSION}"
