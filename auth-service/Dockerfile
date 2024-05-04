FROM golang:1.22-alpine

RUN apk add make

COPY . /app

COPY ./temp/proto /proto

WORKDIR /app

RUN make build-linux

CMD ["./bin/app"]
