# Use specific version for reproducible builds
ARG LITESTREAM_VERSION=0.5
FROM litestream/litestream:${LITESTREAM_VERSION} AS final

LABEL maintainer="atrakic" \
    description="LiteStream database replication service"

RUN apk add --no-cache bash sqlite

# Set default environment variables (non-sensitive only)
ENV DB_SYNC_INTERVAL="10s" \
    DB_PATH="/data/app.db" \
    LITESTREAM_REGION="" \
    DB_REPLICA_URL=""

# Note: Sensitive variables (LITESTREAM_ACCESS_KEY_ID, LITESTREAM_SECRET_ACCESS_KEY)
# should be provided at runtime via docker run -e, docker-compose, or secrets management

WORKDIR /data
# $(dirname "${DB_PATH}")

# Create litestream config directory
RUN mkdir -p /etc/litestream && \
    chmod 755 /etc/litestream

COPY --chown=root:root --chmod=644 ./litestream.yml /etc/litestream.yml
COPY --chown=root:root --chmod=755 ./entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

# Default command (can be overridden)
CMD ["sleep", "infinity"]
