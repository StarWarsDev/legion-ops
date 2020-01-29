all: clean build

build: build-react build-go

build-react:
	cd client && yarn && yarn build

build-go:
	go build .

clean:
	rm -rf client/build

regenerate:
	go run github.com/99designs/gqlgen -v