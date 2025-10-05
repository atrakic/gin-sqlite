#!/bin/bash

## set -ex -o pipefail

set -x

# Restore database from S3 if replica exists
litestream restore -if-replica-exists -v "${DB_PATH}"

# Start replication in the background
litestream replicate "${DB_PATH}" "${DB_REPLICA_URL}" &

# Start the application
exec "$@"
