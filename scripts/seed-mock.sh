#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8001}"

SECURITY_PROFILE_COUNT=6
TRAFFIC_PROFILE_COUNT=4
RESOURCE_COUNT=8
SECURITY_RULE_COUNT=12
TRAFFIC_RULE_COUNT=10
ML_MODEL_COUNT=6

post_json() {
  local path="$1"
  local body="$2"

  curl -sS -X POST \
    -H "Content-Type: application/json" \
    -d "${body}" \
    "${BASE_URL}${path}"
}

extract_id() {
  python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])'
}

security_profile_ids=()
traffic_profile_ids=()

echo "Creating security profiles..."
for i in $(seq 1 "${SECURITY_PROFILE_COUNT}"); do
  base_action="allow"
  if [ $((i % 3)) -eq 1 ]; then
    base_action="block"
  elif [ $((i % 3)) -eq 2 ]; then
    base_action="challenge"
  fi
  log_enabled="true"
  if [ $((i % 2)) -eq 0 ]; then
    log_enabled="false"
  fi
  is_enabled="true"
  if [ $((i % 4)) -eq 0 ]; then
    is_enabled="false"
  fi
  payload=$(cat <<EOF
{
  "name": "Security Profile ${i}",
  "description": "Mock security profile ${i}",
  "base_action": "${base_action}",
  "log_enabled": ${log_enabled},
  "is_enabled": ${is_enabled}
}
EOF
)
  id=$(post_json "/api/v1/security-profiles" "${payload}" | extract_id)
  security_profile_ids+=("${id}")
done

echo "Creating traffic profiles..."
for i in $(seq 1 "${TRAFFIC_PROFILE_COUNT}"); do
  is_enabled="true"
  if [ $((i % 3)) -eq 0 ]; then
    is_enabled="false"
  fi
  payload=$(cat <<EOF
{
  "name": "Traffic Profile ${i}",
  "description": "Mock traffic profile ${i}",
  "is_enabled": ${is_enabled}
}
EOF
)
  id=$(post_json "/api/v1/traffic-profiles" "${payload}" | extract_id)
  traffic_profile_ids+=("${id}")
done

echo "Creating resources..."
for i in $(seq 1 "${RESOURCE_COUNT}"); do
  sec_id="${security_profile_ids[$(( (i - 1) % SECURITY_PROFILE_COUNT ))]}"
  traf_id="${traffic_profile_ids[$(( (i - 1) % TRAFFIC_PROFILE_COUNT ))]}"
  domain="example.com"
  if [ $((i % 3)) -eq 1 ]; then
    domain="corp.local"
  elif [ $((i % 3)) -eq 2 ]; then
    domain="api.test"
  fi
  payload=$(cat <<EOF
{
  "name": "Resource ${i}",
  "url_pattern": "https://app${i}.${domain}/*",
  "security_profile_id": ${sec_id},
  "traffic_profile_id": ${traf_id}
}
EOF
)
  post_json "/api/v1/resources" "${payload}" > /dev/null
done

echo "Creating security rules..."
for i in $(seq 1 "${SECURITY_RULE_COUNT}"); do
  sec_id="${security_profile_ids[$(( (i - 1) % SECURITY_PROFILE_COUNT ))]}"
  rule_type="ip"
  if [ $((i % 3)) -eq 1 ]; then
    rule_type="ua"
  elif [ $((i % 3)) -eq 2 ]; then
    rule_type="rate"
  fi
  action="block"
  if [ $((i % 4)) -eq 1 ]; then
    action="allow"
  elif [ $((i % 4)) -eq 2 ]; then
    action="challenge"
  fi
  dry_run="false"
  if [ $((i % 5)) -eq 0 ]; then
    dry_run="true"
  fi
  is_enabled="true"
  if [ $((i % 6)) -eq 0 ]; then
    is_enabled="false"
  fi
  conditions=$(cat <<EOF
{
  "cidr": "192.168.${i}.0/24",
  "user_agent_contains": "Bot/${i}",
  "path_prefix": "/api/v${i}"
}
EOF
)
  payload=$(cat <<EOF
{
  "profile_id": ${sec_id},
  "name": "Security Rule ${i}",
  "description": "Mock security rule ${i}",
  "priority": ${i},
  "rule_type": "${rule_type}",
  "action": "${action}",
  "conditions": ${conditions},
  "dry_run": ${dry_run},
  "is_enabled": ${is_enabled}
}
EOF
)
  post_json "/api/v1/security-rules" "${payload}" > /dev/null
done

echo "Creating traffic rules..."
for i in $(seq 1 "${TRAFFIC_RULE_COUNT}"); do
  traf_id="${traffic_profile_ids[$(( (i - 1) % TRAFFIC_PROFILE_COUNT ))]}"
  dry_run="false"
  if [ $((i % 4)) -eq 0 ]; then
    dry_run="true"
  fi
  match_all="true"
  if [ $((i % 3)) -eq 0 ]; then
    match_all="false"
  fi
  is_enabled="true"
  if [ $((i % 5)) -eq 0 ]; then
    is_enabled="false"
  fi
  period_seconds=$((30 + (i % 4) * 30))
  requests_limit=$((50 + i * 10))
  conditions=$(cat <<EOF
{
  "method": "GET",
  "path_prefix": "/public/${i}",
  "country": "US"
}
EOF
)
  payload=$(cat <<EOF
{
  "profile_id": ${traf_id},
  "name": "Traffic Rule ${i}",
  "description": "Mock traffic rule ${i}",
  "priority": ${i},
  "dry_run": ${dry_run},
  "match_all": ${match_all},
  "requests_limit": ${requests_limit},
  "period_seconds": ${period_seconds},
  "conditions": ${conditions},
  "is_enabled": ${is_enabled}
}
EOF
)
  post_json "/api/v1/traffic-rules" "${payload}" > /dev/null
done

echo "Creating ML models..."
for i in $(seq 1 "${ML_MODEL_COUNT}"); do
  status="active"
  if [ $((i % 3)) -eq 1 ]; then
    status="training"
  elif [ $((i % 3)) -eq 2 ]; then
    status="archived"
  fi
  config=$(cat <<EOF
{
  "threshold": 0.$((70 + i)),
  "window": $((10 + i))
}
EOF
)
  payload=$(cat <<EOF
{
  "name": "ML Model ${i}",
  "description": "Mock ML model ${i}",
  "version": "1.0.${i}",
  "status": "${status}",
  "config": ${config},
  "artifact_url": "https://example.com/models/model-${i}.bin"
}
EOF
)
  post_json "/api/v1/ml-models" "${payload}" > /dev/null
done

echo "Done. Base URL: ${BASE_URL}"
