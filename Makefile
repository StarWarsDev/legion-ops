#All the things
all: clean build

build: client-build server-build

clean: server-clean client-clean

# Run this with the -j flag ex: `make -j start`
start: server-start client-start

# Client
client-build:
	cd client && yarn && yarn build

client-clean:
	rm -rf client/build

client-start:
	cd client && yarn start

#Server
server-build: server-gql-regenerate
	go build .

server-clean:
	rm -f internal/gql/generated.go \
		  internal/gql/models/generated.go

server-start:
	go run .

server-gql-regenerate:
	go run -v github.com/99designs/gqlgen $1