#!/bin/bash

# Script untuk testing authentication dan authorization
BASE_URL="https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/wechat"

if [ $# -ne 2 ]; then
    echo "Usage: $0 <username> <password>"
    exit 1
fi

USERNAME=$1
PASSWORD=$2

echo "=== Testing Login ==="
# Login untuk mendapatkan token
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "Login Response: $LOGIN_RESPONSE"

# Extract token dari response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Login failed - no token received"
    exit 1
fi

echo "Token: ${TOKEN:0:50}..."

echo ""
echo "=== Testing Debug Token ==="
# Test debug endpoint
DEBUG_RESPONSE=$(curl -s -X GET "$BASE_URL/debug/token" \
    -H "Authorization: Bearer $TOKEN")

echo "Debug Response: $DEBUG_RESPONSE"

# Extract user ID dari debug response
USER_ID=$(echo $DEBUG_RESPONSE | grep -o '"token_user_id":"[^"]*' | cut -d'"' -f4)

if [ -z "$USER_ID" ]; then
    echo "Debug failed - no user ID received"
    exit 1
fi

echo "User ID from token: $USER_ID"

echo ""
echo "=== Testing Update Profile ==="
# Test update profile dengan ID yang sama
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/user/$USER_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"test@example.com\",\"fullname\":\"Test User Updated\"}")

echo "Update Response: $UPDATE_RESPONSE"

echo ""
echo "=== Testing Get Profile ==="
# Test get profile
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/user/$USER_ID" \
    -H "Authorization: Bearer $TOKEN")

echo "Profile Response: $PROFILE_RESPONSE"
