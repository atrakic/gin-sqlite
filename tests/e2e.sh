#!/usr/bin/env bash

set -x

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
URL="${URL:-localhost:8080/api/v1/person}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-secret}"
AUTH_CREDS="$ADMIN_USER:$ADMIN_PASSWORD"

# Test vars
FIRST_NAME=$(shuf -n 1 -e Alice Bob Carol Dave Eve)
LAST_NAME=$(shuf -n 1 -e Smith Johnson Williams Brown Jones)
EMAIL=$(uuidgen | tr '[:upper:]' '[:lower:]' | cut -d'-' -f1)@example.com
JSON_PAYLOAD=$(jq -n --arg fn "$FIRST_NAME" --arg ln "$LAST_NAME" --arg em "$EMAIL" \
  '{first_name: $fn, last_name: $ln, email: $em}')

# functions
function test_create_person() {
  curl --fail --header "Content-Type: application/json" \
    --request POST \
    --include \
    --data "$JSON_PAYLOAD" \
    http://"$AUTH_CREDS"@"$URL"
}

function test_delete_person() {
  local id=$1

  curl -i -X PUT -H "Content-Type: application/json" \
    --data "$JSON_PAYLOAD" \
  http://"$AUTH_CREDS"@"$URL"/"$id"

  curl -i -X "DELETE" http://"$AUTH_CREDS"@"$URL"/"$id"
}

main () {
  test_create_person
  test_delete_person 9999
}

main "$@"
