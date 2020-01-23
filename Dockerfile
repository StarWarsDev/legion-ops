FROM golang as server
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

FROM node as client
COPY --from=server /app/client ./
RUN yarn
RUN yarn build

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=server /main ./
COPY --from=client /build ./client/build
RUN chmod +x ./main
EXPOSE 5000
CMD ./main