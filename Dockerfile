FROM golang:latest as builder
WORKDIR /go/src/github.com/intelligentfish/
COPY go.* main.go ./
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.io,direct && \
    go build
FROM scratch
VOLUME [ "/data", "/logs" ]
WORKDIR /code
COPY --from=builder /go/src/github.com/intelligentfish/dcn .
STOPSIGNAL SIGINT
ENTRYPOINT [ "/code/dcn" ]
