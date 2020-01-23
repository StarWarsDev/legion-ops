FROM golang:1.13-alpine as server
WORKDIR /src
COPY . .
RUN go build -v .

FROM node as client
WORKDIR /src
COPY --from=server /src/client .
RUN yarn
RUN yarn build

FROM alpine

ENV CLIENT_FILES_PATH /app/client/build

RUN apk --no-cache add ca-certificates
COPY --from=server /src/legion-ops /bin
COPY --from=client /src/build/ /app/client/build/
RUN apk update \
  && apk upgrade \
  && apk add --no-cache ca-certificates \
  && update-ca-certificates 2>/dev/null || true
EXPOSE 5000
CMD ["/bin/legion-ops"]
