FROM golang:1.19 as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/dov-id/publisher-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/publisher-svc /go/src/github.com/dov-id/publisher-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/publisher-svc /usr/local/bin/publisher-svc
RUN apk --update add --no-cache musl libstdc++ gcompat libgomp ca-certificates

ENTRYPOINT ["publisher-svc"]
