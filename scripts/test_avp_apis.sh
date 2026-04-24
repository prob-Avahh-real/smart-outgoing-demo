#!/bin/bash

# AVP API Test Script
# Tests the underground automated valet parking endpoints

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/test_lib.sh"

BASE_URL="${BASE_URL:-http://localhost:8080}"
USER_ID="${USER_ID:-demo_user}"
VEHICLE_ID="${VEHICLE_ID:-avp_vehicle_001}"
LOT_ID="${LOT_ID:-lot_1}"
DROPOFF_ZONE="${DROPOFF_ZONE:-B1-DROP-01}"
PICKUP_ZONE="${PICKUP_ZONE:-B1-PICKUP-01}"
TARGET_SPACE_ID="${TARGET_SPACE_ID:-space_1}"

echo "=== AVP API Tests ==="
echo "Base URL: $BASE_URL"
echo "User ID: $USER_ID"
echo "Vehicle ID: $VEHICLE_ID"
echo ""

echo "1. Start AVP auto-park task..."
START_RESPONSE=$(api_post "$BASE_URL" "/api/parking/avp/start" "{
    \"vehicle_id\": \"$VEHICLE_ID\",
    \"parking_lot_id\": \"$LOT_ID\",
    \"dropoff_zone\": \"$DROPOFF_ZONE\",
    \"target_space_id\": \"$TARGET_SPACE_ID\"
  }" "$USER_ID")

if assert_has_key "$START_RESPONSE" '.task.id'; then
  TASK_ID=$(echo "$START_RESPONSE" | jq -r '.task.id')
  SESSION_ID=$(echo "$START_RESPONSE" | jq -r '.parking_session.id')
  echo "  SUCCESS: Started auto-park task: $TASK_ID"
  echo "  Linked session: $SESSION_ID"
else
  echo "  FAILED: could not start AVP task"
  echo "  Response: $START_RESPONSE"
  exit 1
fi
echo ""

echo "2. Query AVP task status..."
QUERY_RESPONSE=$(api_get "$BASE_URL" "/api/parking/avp/tasks/$TASK_ID")
if assert_has_key "$QUERY_RESPONSE" '.task.status'; then
  STATUS=$(echo "$QUERY_RESPONSE" | jq -r '.task.status')
  PROGRESS=$(echo "$QUERY_RESPONSE" | jq -r '.task.progress')
  CHECKPOINT=$(echo "$QUERY_RESPONSE" | jq -r '.task.last_checkpoint')
  echo "  SUCCESS: status=$STATUS progress=$PROGRESS checkpoint=$CHECKPOINT"
else
  echo "  FAILED: could not query AVP task"
  echo "  Response: $QUERY_RESPONSE"
  exit 1
fi
echo ""

echo "3. Cancel AVP auto-park task..."
CANCEL_RESPONSE=$(api_post "$BASE_URL" "/api/parking/avp/tasks/$TASK_ID/cancel" '{}')
if assert_has_key "$CANCEL_RESPONSE" '.task.status'; then
  CANCEL_STATUS=$(echo "$CANCEL_RESPONSE" | jq -r '.task.status')
  CANCEL_CHECKPOINT=$(echo "$CANCEL_RESPONSE" | jq -r '.task.last_checkpoint')
  echo "  SUCCESS: status=$CANCEL_STATUS checkpoint=$CANCEL_CHECKPOINT"
else
  echo "  FAILED: could not cancel AVP task"
  echo "  Response: $CANCEL_RESPONSE"
  exit 1
fi
echo ""

echo "4. Start AVP summon task..."
SUMMON_RESPONSE=$(api_post "$BASE_URL" "/api/parking/avp/summon" "{
    \"vehicle_id\": \"$VEHICLE_ID\",
    \"parking_lot_id\": \"$LOT_ID\",
    \"pickup_zone\": \"$PICKUP_ZONE\"
  }" "$USER_ID")

if assert_has_key "$SUMMON_RESPONSE" '.task.id'; then
  SUMMON_TASK_ID=$(echo "$SUMMON_RESPONSE" | jq -r '.task.id')
  echo "  SUCCESS: Started summon task: $SUMMON_TASK_ID"
else
  echo "  FAILED: could not start summon task"
  echo "  Response: $SUMMON_RESPONSE"
  exit 1
fi
echo ""

echo "5. Query summon task status..."
SUMMON_QUERY_RESPONSE=$(api_get "$BASE_URL" "/api/parking/avp/tasks/$SUMMON_TASK_ID")
if assert_has_key "$SUMMON_QUERY_RESPONSE" '.task.status'; then
  SUMMON_STATUS=$(echo "$SUMMON_QUERY_RESPONSE" | jq -r '.task.status')
  SUMMON_PROGRESS=$(echo "$SUMMON_QUERY_RESPONSE" | jq -r '.task.progress')
  echo "  SUCCESS: summon status=$SUMMON_STATUS progress=$SUMMON_PROGRESS"
else
  echo "  FAILED: could not query summon task"
  echo "  Response: $SUMMON_QUERY_RESPONSE"
  exit 1
fi
echo ""

echo "=== AVP Test Summary ==="
echo "AVP API smoke tests completed successfully."
echo "Tip: wait 30-90s and query task again to observe status auto-advance."
