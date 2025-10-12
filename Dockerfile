ARG GO_VERSION=1.24-alpine

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS builder
WORKDIR /src
RUN apk --update add ca-certificates
ARG TARGETARCH
RUN adduser -D -u 1001 nonroot
COPY go.mod go.sum ./

# Use cache mounts to speed up the installation of existing dependencies
# Mount go.sum as a cache key to invalidate cache when dependencies change
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/tmp/go-sum,id=go-sum-${TARGETARCH} \
    cp go.sum /tmp/go-sum/go.sum.${TARGETARCH} && \
    go mod download -x && \
    go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build \
    -ldflags="-w -s" \
    -v -o /bin/server

FROM builder AS testing
RUN go vet -v ./...
RUN go test -v ./...

FROM alpine:latest as final
LABEL maintainer="Admir Trakic <atrakic@users.noreply.github.com>"
RUN apk --no-cache add curl ca-certificates
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/server /
USER nonroot
ENV GIN_MODE=release
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ping || exit 1
CMD ["/server"]
