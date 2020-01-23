FROM golang as server
ADD . /src
WORKDIR /src
RUN go mod download
RUN go build .

FROM node as client
WORKDIR /src
COPY --from=server /src/client/ ./
RUN yarn
RUN yarn build

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=server /src/legion-ops ./legion-ops
COPY --from=client /src/build/ ./client/build/
RUN chmod +x /app/legion-ops
RUN apk update && apk upgrade && apk add --no-cache ca-certificates && update-ca-certificates 2>/dev/null || true
EXPOSE 5000
CMD ["/bin/sh", "-c", "./legion-ops"]
