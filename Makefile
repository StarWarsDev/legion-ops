all: clean build

build: build-react build-go

build-react:
	cd client && yarn && yarn build

build-go: regenerate
	go build .

start-server:
	go run .

clean: clean-go clean-react

clean-go:
	rm -f internal/gql/generated.go \
		  internal/gql/models/generated.go \
		  internal/gql/resolvers/generated.go

clean-react:
	rm -rf client/build

regenerate:
	go run -v github.com/99designs/gqlgen $1