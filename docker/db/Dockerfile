FROM alpine:latest

RUN set -x \
    && apk add --no-cache bash sqlite

COPY --chmod=0755 entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
