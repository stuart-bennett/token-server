FROM golang:1.11-alpine
ADD . /go/src/github.com/stuart-bennett/token-server
WORKDIR /go/src/github.com/stuart-bennett/token-server/server
RUN apk add git && \
    go get github.com/gomodule/redigo/redis && \
    go install
CMD ["/go/bin/server"]
EXPOSE 8000
