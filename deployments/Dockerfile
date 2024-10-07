FROM golang:1.22.0-alpine AS builder
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
RUN apk add --no-cache qemu

WORKDIR /
COPY --from=builder /go/src/github.com/shubhamgptln/hotels-data-merge/build/api /
COPY --from=builder /go/src/github.com/shubhamgptln/hotels-data-merge/config.yml /
EXPOSE 8090
CMD ["./api"]
