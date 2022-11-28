#! /usr/bin/env sh

set -e

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)

# run the build
$SCRIPT_DIR/run.sh

curl --header "Content-Type: application/json" \
  -H "X-HTTP-Method-Override: PUT" \
  --request POST \
  --data '{"first_name":"xyz","last_name":"xyz","email":"aaaa@bar.com"}' \
  localhost:8080/api/v1/person

curl -X PUT --url localhost:8080/api/v1/person/2 --data 'email=bbbb\@bar.com'

#sqlite3 person.db "select * from people where id = 2"


curl -X "DELETE" localhost:8080/api/v1/person/2

sqlite3 person.db "select * from people"
