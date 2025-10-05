FROM golang:1.21 AS builder
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

FROM alpine AS final
LABEL maintainer="Admir Trakic <atrakic@users.noreply.github.com>"

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/server /
COPY --chmod=0755 entrypoint.sh /entrypoint.sh

# Install Litestream for SQLite replication
COPY --from=litestream/litestream /usr/local/bin/litestream /bin/litestream
COPY ./litestream.yml /etc/litestream.yml

RUN apk add --no-cache bash

USER nonroot

ENV GIN_MODE=release

# Frequency that database snapshots are replicated.
ENV DB_SYNC_INTERVAL="10s"
ENV DATABASE_FILE="/var/tmp/database.db"
ENV LITESTREAM_CONFIG_FILE="/etc/litestream.yml"

# Expose port 8080 to the outside world
EXPOSE 8080

# CMD ["/server"]
ENTRYPOINT ["/entrypoint.sh" ]
CMD []
