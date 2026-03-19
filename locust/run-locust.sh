#!/bin/sh
set -eu

LOCUST_FILES="${LOCUSTFILES:-$(find /mnt/locust -maxdepth 1 -type f -name '*.py' ! -name 'run-locust.sh' | sort | paste -sd, -)}"

if [ -z "${LOCUST_FILES}" ]; then
  echo "No locust scripts found in /mnt/locust"
  exit 1
fi

HOST="${LOCUST_HOST:-http://envoy:10000}"

if [ -n "${LOCUST_OPTS:-}" ]; then
  exec locust -f "${LOCUST_FILES}" --host "${HOST}" --class-picker ${LOCUST_OPTS}
fi

exec locust -f "${LOCUST_FILES}" --host "${HOST}" --class-picker
