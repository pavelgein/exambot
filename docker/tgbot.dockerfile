FROM golang:1.14 AS builder
ARG src_dir
RUN mkdir /src && mkdir /build
COPY ${src_dir} /src
RUN cd /src/app/tgbot && go build -o /build/tgbot
RUN cd /src/app/uploader && go build -o /build/uploader


FROM golang:1.14
COPY --from=builder /build/tgbot /app/tgbot
COPY --from=builder /build/uploader /app/uploader
ENTRYPOINT ["/app/tgbot"]
