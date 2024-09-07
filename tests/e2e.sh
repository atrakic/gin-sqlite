#!/usr/bin/env bash

set -x

set -o errexit
set -o nounset
set -o pipefail

#SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
#$SCRIPT_ROOT/run.sh

URL="localhost:8080/api/v1/person"
echo "$URL"

curl --fail --header "Content-Type: application/json" \
  --request POST \
  --include \
  --data '{"first_name":"xyz","last_name":"xyz","email":"xyz@bar.com"}' \
  http://admin:secret@localhost:8080/api/v1/person

curl -i -X PUT -H "Content-Type: application/json" --data '{ "first_name": "Test", "last_name":"Test", "email":"xyz@bar.com"}' \
  http://admin:secret@localhost:8080/api/v1/person/2

curl -i -X "DELETE" http://admin:secret@localhost:8080/api/v1/person/2
