FROM golang:1.13-alpine as server
WORKDIR /src
COPY . .
RUN go build -v .

FROM alpine

ENV CLIENT_FILES_PATH /app/client/build

RUN apk --no-cache add ca-certificates
COPY --from=server /src/legion-ops /bin
RUN apk update \
  && apk upgrade \
  && apk add --no-cache ca-certificates \
  && update-ca-certificates 2>/dev/null || true
EXPOSE 5000
CMD ["/bin/legion-ops"]
