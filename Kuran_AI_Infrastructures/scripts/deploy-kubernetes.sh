#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$INFRA_DIR/.." && pwd)"
K8S_DIR="$INFRA_DIR/k8s"
ENV_FILE="$INFRA_DIR/.env"
NAMESPACE="kuranai"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Missing $ENV_FILE"
  exit 1
fi

for command_name in docker kubectl; do
  command -v "$command_name" >/dev/null || { echo "Missing command: $command_name"; exit 1; }
done

if [[ -z "${OPENAI_API_KEY:-}" ]]; then
  read -rsp "OpenAI API key: " OPENAI_API_KEY
  echo
fi
if [[ -z "$OPENAI_API_KEY" ]]; then
  echo "OpenAI API key cannot be empty"
  exit 1
fi

kubectl apply -f "$K8S_DIR/namespace.yaml"

kubectl create secret generic kuranai-secrets \
  --namespace "$NAMESPACE" \
  --from-env-file "$ENV_FILE" \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic kuranai-openai \
  --namespace "$NAMESPACE" \
  --from-literal "OPENAI_API_KEY=$OPENAI_API_KEY" \
  --dry-run=client -o yaml | kubectl apply -f -
unset OPENAI_API_KEY

if [[ -n "${GOOGLE_PLAY_SERVICE_ACCOUNT_FILE:-}" ]]; then
  if [[ ! -f "$GOOGLE_PLAY_SERVICE_ACCOUNT_FILE" ]]; then
    echo "Google Play service-account file not found: $GOOGLE_PLAY_SERVICE_ACCOUNT_FILE" >&2
    exit 1
  fi
  kubectl create secret generic kuranai-google-play \
    --namespace "$NAMESPACE" \
    --from-file "service-account.json=$GOOGLE_PLAY_SERVICE_ACCOUNT_FILE" \
    --dry-run=client -o yaml | kubectl apply -f -
fi

docker build -t kuranai-gateway:latest "$PROJECT_ROOT/Kuran_AI_Gateway"
docker build \
  --build-arg VITE_API_BASE=https://kuranai-api.talhacan.com \
  --build-arg VITE_GOOGLE_CLIENT_ID=157832066999-3tgb2daavmrd1d1vd10ptrato1iu382t.apps.googleusercontent.com \
  --build-arg VITE_APPLE_CLIENT_ID=com.talhacan.kuranai.web \
  --build-arg VITE_APPLE_REDIRECT_URI=https://kuranai.talhacan.com/auth/apple/callback \
  -t kuranai-web:latest "$PROJECT_ROOT/Kuran_AI_Web"
docker build \
  --build-arg VITE_API_BASE=https://kuranai-api.talhacan.com \
  -t kuranai-admin:latest "$PROJECT_ROOT/Kuran_AI_Admin"

kubectl apply -f "$K8S_DIR/config.yaml"
kubectl apply -f "$K8S_DIR/data.yaml"
kubectl rollout status statefulset/postgres -n "$NAMESPACE" --timeout=240s
kubectl rollout status statefulset/redis -n "$NAMESPACE" --timeout=240s

kubectl delete job gateway-migrate -n "$NAMESPACE" --ignore-not-found
kubectl apply -f "$K8S_DIR/migration.yaml"
kubectl wait --for=condition=complete job/gateway-migrate -n "$NAMESPACE" --timeout=240s

kubectl apply -f "$K8S_DIR/apps.yaml"
kubectl rollout restart deployment/gateway deployment/web deployment/admin -n "$NAMESPACE"
kubectl rollout status deployment/gateway -n "$NAMESPACE" --timeout=300s
kubectl rollout status deployment/web -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/admin -n "$NAMESPACE" --timeout=180s

kubectl get pods,svc,pvc -n "$NAMESPACE"
