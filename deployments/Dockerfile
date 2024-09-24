FROM golang:1.21.0-alpine AS builder
RUN go version
RUN apk add --no-cache bash make

COPY . /go/src/github.com/shubhamgptln/hotels-data-merge/
WORKDIR /go/src/github.com/shubhamgptln/hotels-data-merge/

RUN set -eux && \
    go mod download && \
    go mod verify

RUN make build

FROM alpine:3.18.3
RUN apk add --no-cache bash
RUN apt-get install -y qemu qemu-user-static

WORKDIR /
COPY --from=builder /go/src/github.com/shubhamgptln/hotels-data-merge/build/api /
EXPOSE 8090
CMD ["./api"]