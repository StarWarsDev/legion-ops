FROM golang as server
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build .

FROM node as client
COPY --from=server /app/client ./
RUN yarn
RUN yarn build

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=server /app/legion-ops ./
COPY --from=client /build ./client/build
RUN chmod +x ./legion-ops
EXPOSE 5000
CMD ./legion-ops