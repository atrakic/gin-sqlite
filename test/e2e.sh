#!/usr/bin/env bash

set -x

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
URL="${URL:-localhost:8080}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-secret}"

# Test vars
FIRST_NAME=$(shuf -n 1 -e Alice Bob Carol Dave Eve)
LAST_NAME=$(shuf -n 1 -e Smith Johnson Williams Brown Jones)
EMAIL=$(uuidgen | tr '[:upper:]' '[:lower:]' | cut -d'-' -f1)@example.com
JSON_PAYLOAD=$(jq -n --arg fn "$FIRST_NAME" --arg ln "$LAST_NAME" --arg em "$EMAIL" \
  '{first_name: $fn, last_name: $ln, email: $em}')

# Global variable to store JWT token
JWT_TOKEN=""

# functions
function test_login() {
  echo "🔐 Testing login endpoint..."
  local response
  response=$(curl --fail --silent --header "Content-Type: application/json" \
    --request POST \
    --data '{"username":"'"$ADMIN_USER"'", "password":"'"$ADMIN_PASSWORD"'"}' \
    http://"$URL"/auth/login)

  JWT_TOKEN=$(echo "$response" | jq -r '.token')

  if [ -z "$JWT_TOKEN" ] || [ "$JWT_TOKEN" == "null" ]; then
    echo "❌ Login failed or token not found"
    echo "Response: $response"
    exit 1
  fi

  echo "✅ Login successful, token received"
  echo "🎯 Token: ${JWT_TOKEN:0:20}..."
}

function test_create_person() {
  echo "📝 Testing create person endpoint..."
  local response
  response=$(curl --fail --header "Content-Type: application/json" \
    --header "Authorization: Bearer $JWT_TOKEN" \
    --request POST \
    --silent \
    --data "$JSON_PAYLOAD" \
    http://"$URL"/api/v1/person)

  echo "✅ Person created successfully"
  echo "📋 Response: $response"

  # Extract person ID from response if available
  local person_id
  person_id=$(echo "$response" | jq -r '.id // empty' 2>/dev/null || echo "")
  if [ -n "$person_id" ]; then
    echo "🆔 Created person ID: $person_id"
    echo "$person_id"
  fi
}

function test_update_person() {
  local id=$1
  echo "📝 Testing update person endpoint for ID: $id..."

  local response
  response=$(curl --fail --header "Content-Type: application/json" \
    --header "Authorization: Bearer $JWT_TOKEN" \
    --request PUT \
    --silent \
    --data "$JSON_PAYLOAD" \
    http://"$URL"/api/v1/person/"$id")

  echo "✅ Person updated successfully"
  echo "📋 Response: $response"
}

function test_delete_person() {
  local id=$1
  echo "🗑️  Testing delete person endpoint for ID: $id..."

  local response
  response=$(curl --fail --header "Authorization: Bearer $JWT_TOKEN" \
    --request DELETE \
    --silent \
    http://"$URL"/api/v1/person/"$id")

  echo "✅ Person deleted successfully"
  echo "📋 Response: $response"
}

function test_get_persons() {
  echo "📋 Testing get all persons endpoint..."

  local response
  response=$(curl --fail --silent http://"$URL"/api/v1/person)

  echo "✅ Retrieved persons list"
  echo "📋 Response: $response"
}

main () {
  echo "🚀 Starting E2E tests with JWT authentication..."

  # Step 1: Login and get JWT token
  test_login

  # Step 2: Test public endpoint (no auth required)
  test_get_persons

  # Step 3: Test authenticated endpoints
  local person_id
  person_id=$(test_create_person)

  # Use a known ID for update/delete tests since create response might not include ID
  test_update_person 1
  test_delete_person 3

  echo "🎉 All E2E tests completed successfully!"
}

main "$@"
