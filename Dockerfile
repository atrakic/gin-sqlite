FROM golang:1.19 AS builder
WORKDIR /src

RUN useradd -u 1001 nonroot

COPY go.mod ./


# Use cache mounts to speed up the installation of existing dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN go build \
    #-ldflags="-w -s" \
    -ldflags="-linkmode external -extldflags -static" \
    -v -o /bin/server

FROM scratch
LABEL maintainer="Admir Trakic <atrakic@users.noreply.github.com>"
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/server /
USER nonroot
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/server"]
