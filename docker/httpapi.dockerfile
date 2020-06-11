# based on https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i

FROM golang:1.14 AS builder
WORKDIR /src
ARG src_dir
RUN  mkdir /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY ${src_dir} /src
RUN cd /src/app/httpapi && go build -o /build/httpapi
RUN cd /src/app/granter && go build -o /build/granter

WORKDIR /dist
RUN cp /build/httpapi ./httpapi
RUN cp /build/granter ./granter

# Optional: in case your application uses dynamic linking (often the case with CGO),
# this will collect dependent libraries so they're later copied to the final image
# NOTE: make sure you honor the license terms of the libraries you copy and distribute
RUN ls
RUN for p in "httpapi granter" ; do ldd $p | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;' ; done
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/


FROM scratch
COPY --chown=0:0 --from=builder /dist /
ENTRYPOINT ["/httpapi"]
