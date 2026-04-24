#!/bin/bash

api_get() {
  local base_url="$1"
  local path="$2"
  curl -s "${base_url}${path}"
}

api_post() {
  local base_url="$1"
  local path="$2"
  local body="$3"
  local user_id="$4"

  if [ -n "$user_id" ]; then
    curl -s -X POST "${base_url}${path}" \
      -H "Content-Type: application/json" \
      -H "x-user-id: ${user_id}" \
      -d "$body"
    return
  fi

  curl -s -X POST "${base_url}${path}" \
    -H "Content-Type: application/json" \
    -d "$body"
}

assert_has_key() {
  local json="$1"
  local key_path="$2"
  echo "$json" | jq -e "$key_path" >/dev/null 2>&1
}
