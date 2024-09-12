FROM alpine

RUN set -x \
    && apk add --no-cache bash sqlite

# https://gist.github.com/andyshinn/3ae01fa13cb64c9d36e7?permalink_comment_id=2044506#gistcomment-2044506
COPY --chmod=0755 ./backup.sh /usr/local/bin/db-backup

RUN echo '*/15  *  *  *  *   /usr/local/bin/db-backup' > /etc/crontabs/root

CMD ["sh", "-c", "crond -l 2 -f"]