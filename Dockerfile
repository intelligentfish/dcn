FROM golang:latest as builder
WORKDIR /go/src/github.com/intelligentfish/dcn
COPY agent ./agent
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
COPY service ./service
COPY serviceGroup ./serviceGroup
COPY types ./types
COPY go.sum go.mod ./
COPY main.go .
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.io,direct && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"'
FROM scratch
VOLUME [ "/data", "/var/log" ]
WORKDIR /code
COPY --from=builder /go/src/github.com/intelligentfish/dcn/dcn .
HEALTHCHECK --interval=5s --timeout=30s --start-period=5s --retries=3 CMD [ "curl", "http://localhost/api/v1/health" ]
STOPSIGNAL SIGINT
ENTRYPOINT [ "/code/dcn" ]
