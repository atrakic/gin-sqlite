#!/bin/bash -x

set -e

database_file="${DATABASE_FILE:-/var/tmp/database.db}"

## https://gcollazo.com/optimal-sqlite-settings-for-django/

sqlite3 "$database_file" <<EOF
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
CREATE TABLE IF NOT EXISTS people (
    id INTEGER PRIMARY KEY AUTOINCREMENT unique,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL
);
-- Data:
INSERT or IGNORE INTO people VALUES (1, 'Foo','Bar','foo@bar.com');
EOF

while true; do
    sqlite3 "$database_file" "select count(*) from people;"
    sleep 1
done
