GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/go-getting-started

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

all: build

build: build-react build-go

build-react:
	cd client && yarn && yarn build

build-go:
	go build .

clean:
	rm -rf $(DOCKER_BUILD)
	rm -rf client/build

heroku: $(DOCKER_CMD)
	heroku container:push web