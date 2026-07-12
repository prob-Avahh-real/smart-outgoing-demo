#!/bin/bash

# AI Car Parking API Test Script
# Tests all parking-related endpoints

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/test_lib.sh"

BASE_URL="${BASE_URL:-http://localhost:8080}"
USER_ID="${USER_ID:-demo_user}"

echo "=== AI Car Parking API Tests ==="
echo "Base URL: $BASE_URL"
echo "User ID: $USER_ID"
echo ""

# Test 1: Find Parking Spots
echo "1. Testing Find Parking API..."
RESPONSE=$(api_post "$BASE_URL" "/api/parking/find" '{
    "latitude": 22.6913,
    "longitude": 114.0448,
    "max_price": 20,
    "max_distance": 5,
    "limit": 5
  }')

if assert_has_key "$RESPONSE" '.recommendations'; then
    COUNT=$(echo "$RESPONSE" | jq '.recommendations | length')
    echo "  SUCCESS: Found $COUNT parking recommendations"
    echo "  First recommendation: $(echo "$RESPONSE" | jq -r '.recommendations[0].parking_lot.name')"
else
    echo "  FAILED: Invalid response"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 2: Get All Parking Lots
echo "2. Testing Get Parking Lots API..."
RESPONSE=$(api_get "$BASE_URL" "/api/parking/lots")

if assert_has_key "$RESPONSE" '.parking_lots'; then
    COUNT=$(echo "$RESPONSE" | jq '.parking_lots | length')
    echo "  SUCCESS: Found $COUNT parking lots"
else
    echo "  FAILED: Invalid response"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 3: Get Specific Parking Lot
echo "3. Testing Get Parking Lot Details API..."
RESPONSE=$(api_get "$BASE_URL" "/api/parking/lots/lot_1")

if assert_has_key "$RESPONSE" '.parking_lot'; then
    NAME=$(echo "$RESPONSE" | jq -r '.parking_lot.name')
    echo "  SUCCESS: Retrieved parking lot: $NAME"
else
    echo "  FAILED: Invalid response"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 4: Get Parking Spaces
echo "4. Testing Get Parking Spaces API..."
RESPONSE=$(api_get "$BASE_URL" "/api/parking/lots/lot_1/spaces")

if assert_has_key "$RESPONSE" '.parking_spaces'; then
    COUNT=$(echo "$RESPONSE" | jq '.parking_spaces | length')
    echo "  SUCCESS: Found $COUNT parking spaces"
else
    echo "  FAILED: Invalid response"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 5: Reserve Parking Space
echo "5. Testing Reserve Space API..."
RESPONSE=$(api_post "$BASE_URL" "/api/parking/reserve" '{
    "parking_lot_id": "lot_1",
    "space_id": "space_1",
    "start_time": "2026-04-21T12:00:00Z",
    "end_time": "2026-04-21T14:00:00Z"
  }' "$USER_ID")

if assert_has_key "$RESPONSE" '.reservation'; then
    RESERVATION_ID=$(echo "$RESPONSE" | jq -r '.reservation.id')
    echo "  SUCCESS: Created reservation: $RESERVATION_ID"
else
    echo "  FAILED: Invalid response"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 6: Start Parking Session
echo "6. Testing Start Parking Session API..."
RESPONSE=$(api_post "$BASE_URL" "/api/parking/session/start" '{
    "parking_lot_id": "lot_1",
    "space_id": "space_1"
  }' "$USER_ID")

if assert_has_key "$RESPONSE" '.session'; then
    SESSION_ID=$(echo "$RESPONSE" | jq -r '.session.id')
    echo "  SUCCESS: Started parking session: $SESSION_ID"
else
    echo "  FAILED: Invalid response"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 7: Health Check
echo "7. Testing Server Health..."
RESPONSE=$(api_get "$BASE_URL" "/api/config")

if assert_has_key "$RESPONSE" '.amap_js_key'; then
    echo "  SUCCESS: Server is running and responding"
else
    echo "  FAILED: Server not responding properly"
    echo "  Response: $RESPONSE"
fi
echo ""

# Test 8: Parking UI Access
echo "8. Testing Parking UI Access..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/parking")

if [ "$HTTP_CODE" = "200" ]; then
    echo "  SUCCESS: Parking UI accessible (HTTP 200)"
else
    echo "  FAILED: Parking UI not accessible (HTTP $HTTP_CODE)"
fi
echo ""

echo "=== Test Summary ==="
echo "All API tests completed."
echo "Check the results above for any failures."
echo ""
echo "To test the full user experience:"
echo "1. Open $BASE_URL/parking in your browser"
echo "2. Click 'GPS' to get your location"
echo "3. Set search criteria and click 'Find Parking Spots'"
echo "4. Try reserving a space and starting a session"
echo ""
