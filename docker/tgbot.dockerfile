FROM golang:1.14 AS builder
WORKDIR /src
ARG src_dir
RUN  mkdir /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY ${src_dir} /src
RUN cd /src/app/tgbot && go build -o /build/tgbot
RUN cd /src/app/uploader && go build -o /build/uploader


FROM golang:1.14
COPY --from=builder /build/tgbot /app/tgbot
COPY --from=builder /build/uploader /app/uploader
ENTRYPOINT ["/app/tgbot"]
