VERSION = "master"

build:
	go build -ldflags="-s -w -X example/version.Version=${VERSION}"

build-docker:
	docker build . -t example --build-arg VERSION=${VERSION}
