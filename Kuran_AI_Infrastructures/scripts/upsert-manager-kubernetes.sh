#!/usr/bin/env bash
set -Eeuo pipefail

readonly NAMESPACE="kuranai"

: "${MANAGER_EMAIL:?MANAGER_EMAIL is required}"
: "${MANAGER_PASSWORD:?MANAGER_PASSWORD is required}"
MANAGER_NAME="${MANAGER_NAME:-Kuran AI Manager}"
MANAGER_USERNAME="${MANAGER_USERNAME:-}"

gateway_pod="$(kubectl get pods -n "$NAMESPACE" -l app=kuranai-gateway \
  --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}')"
[[ -n "$gateway_pod" ]] || { echo "No running gateway pod found" >&2; exit 1; }

kubectl exec -n "$NAMESPACE" "$gateway_pod" -- \
  env \
    "MANAGER_EMAIL=$MANAGER_EMAIL" \
    "MANAGER_PASSWORD=$MANAGER_PASSWORD" \
    "MANAGER_NAME=$MANAGER_NAME" \
    "MANAGER_USERNAME=$MANAGER_USERNAME" \
    /app/gateway manager:create

unset MANAGER_PASSWORD
