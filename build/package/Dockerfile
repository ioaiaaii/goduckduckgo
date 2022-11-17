# Multi-stage docker build
##    Build base, parse LD_FLAGS, add user and build the service bin to /tmp dir.
FROM golang@sha256:c6e6813d9027c18f80fa59ab132269c47d8d845a510895ed5a74a3e677fe78f6 as base
ARG LD_FLAGS
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    goduckduckgo
WORKDIR /$GOPATH/src/goduckduckgo/
COPY . .
RUN go mod download &&\
    go mod verify &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -mod=vendor -ldflags="${LD_FLAGS}" -o /tmp/goduckduckgo ./cmd/goduckduckgo/*

#   Build image
##    Copy files, document the default port, and run the bin as non-root user
FROM gcr.io/distroless/base-debian11:nonroot@sha256:21a86d92771aa0ee8fc238e94f1a3a5fbb2c0ad55adda8ad626a6dd1e7faaa8f
LABEL org.opencontainers.image.url="https://github.com/ioaiaaii/goduckduckgo"
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /tmp/goduckduckgo /usr/local/bin/goduckduckgo
EXPOSE 8080/tcp
USER goduckduckgo
ENTRYPOINT ["/usr/local/bin/goduckduckgo"]
