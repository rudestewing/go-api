#!/bin/bash

# Quick API test script
echo "🚀 Testing Go API endpoints..."

# Check if server is running
if ! curl -s http://localhost:8000/api/v1/ > /dev/null; then
    echo "❌ Server is not running. Please start with: make dev"
    exit 1
fi

echo "✅ Server is running"

# Test health endpoint
echo "🔍 Testing health endpoint..."
response=$(curl -s http://localhost:8000/api/v1/)
echo "Response: $response"

# Test user registration
echo "🔍 Testing user registration..."
register_response=$(curl -s -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "testpassword123"
  }')
echo "Registration response: $register_response"

# Test user login
echo "🔍 Testing user login..."
login_response=$(curl -s -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123"
  }')
echo "Login response: $login_response"

# Extract token (basic extraction, might need adjustment based on response format)
token=$(echo $login_response | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$token" ]; then
    echo "🔍 Testing protected endpoint with token..."
    profile_response=$(curl -s -X GET http://localhost:8000/api/v1/user/profile \
      -H "Authorization: Bearer $token")
    echo "Profile response: $profile_response"
else
    echo "⚠️  No token received, skipping protected endpoint test"
fi

echo "✅ API testing completed!"
