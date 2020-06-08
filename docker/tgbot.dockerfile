FROM golang:1.14 AS builder
ARG src_dir
RUN mkdir /src && mkdir /build
COPY ${src_dir} /src
RUN cd /src/app/tgbot && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a -o /build/tgbot


FROM golang:1.14
COPY --from=builder /build/tgbot /app/tgbot
ENTRYPOINT ["/app/tgbot"]
