FROM alpine

LABEL maintainer="atrakic" \
    description="Rolling Sqlite backup service."

RUN set -x \
    && apk add --no-cache bash sqlite py-pip curl

# https://gist.github.com/andyshinn/3ae01fa13cb64c9d36e7?permalink_comment_id=2044506#gistcomment-2044506
COPY --chmod=0755 ./backup.sh /usr/local/bin/db-backup

RUN echo '*/5  *  *  *  *   /usr/local/bin/db-backup' > /etc/crontabs/root

# Add healthcheck to verify crond is running
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep -f "crond" > /dev/null || exit 1

CMD ["sh", "-c", "crond -l 2 -f"]
