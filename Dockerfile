FROM golang:1.9.3-alpine3.7

RUN apk --no-cache --update add git

RUN go get -u \
    cloud.google.com/go/pubsub \
    golang.org/x/net/context \
    github.com/go-redis/redis \
    github.com/go-pg/pg \
    github.com/logpacker/PayPal-Go-SDK \
    github.com/gorilla/mux \
    github.com/satori/go.uuid \
    github.com/stretchr/testify/assert

ADD . /go/src/billing-api

WORKDIR /go/src/billing-api

CMD go run app/main.go

