FROM golang:1.12.0-alpine3.9
RUN apk update \
    && apk upgrade \
    && apk add --no-cache git bash make
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -v -a -installsuffix cgo -o cmd/server.go
CMD ["/app/main"]