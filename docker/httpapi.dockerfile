FROM golang:1.14 AS builder
WORKDIR /src
ARG src_dir
RUN  mkdir /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY ${src_dir} /src
RUN cd /src/app/httpapi && go build -o /build/httpapi


FROM golang:1.14
COPY --from=builder /build/httpapi /app/httpapi
ENTRYPOINT ["/app/httpapi"]
