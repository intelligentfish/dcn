FROM golang:latest as builder
WORKDIR /go/src/github.com/intelligentfish/dcn
COPY app ./app
COPY config ./config
COPY crypto ./crypto
COPY define ./define
COPY io ./io
COPY log ./log
COPY mqMessageHandler ./mqMessageHandler
COPY proto ./proto
COPY serviceSorter ./serviceSorter
COPY signalHandler ./signalHandler
COPY srv ./srv
COPY srvGroup ./srvGroup
COPY types ./types
COPY go.sum go.mod ./
COPY main.go .
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.io,direct && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"'
FROM scratch
VOLUME [ "/data", "/logs" ]
WORKDIR /code
COPY --from=builder /go/src/github.com/intelligentfish/dcn/dcn .
STOPSIGNAL SIGINT
ENTRYPOINT [ "/code/dcn" ]
