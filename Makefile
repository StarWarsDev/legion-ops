all: build

build: build-react build-go

build-react:
	cd client && yarn && yarn build

build-go:
	go build .