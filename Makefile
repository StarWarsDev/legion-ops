#All the things
all: clean build

build:
	go build .

clean:
	rm -f internal/gql/generated.go \
		  internal/gql/models/generated.go

start:
	go run .

gql-regenerate:
	go run -v github.com/99designs/gqlgen $1