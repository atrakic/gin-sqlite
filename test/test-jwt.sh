#!/bin/bash

# JWT Authentication Test Script
set -e

echo "🚀 Starting JWT Authentication Test..."

# Start server in background
echo "📝 Starting server on port 8083..."
PORT=8083 go run cmd/server/main.go &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo "🔐 Testing login endpoint..."
# Test login and capture token
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8083/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secret"}')

echo "📋 Login response: $LOGIN_RESPONSE"

# Extract token (assuming JSON response)
if command -v jq &> /dev/null; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
    echo "🎯 Extracted token: ${TOKEN:0:50}..."

    # Test authenticated endpoint
    echo "🔒 Testing authenticated endpoint..."
    curl -s -X POST http://localhost:8083/api/v1/person \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"first_name":"John","last_name":"Doe","email":"john.doe@example.com"}'

    echo ""
    echo "✅ JWT authentication test completed!"
else
    echo "⚠️  jq not available, showing raw response"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
echo "🧹 Server stopped"
