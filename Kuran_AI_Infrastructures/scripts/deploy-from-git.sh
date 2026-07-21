#!/usr/bin/env bash
set -Eeuo pipefail

readonly EXPECTED_ROOT="/www/wwwroot/kuranai"
readonly NAMESPACE="kuranai"
readonly LOCK_FILE="/var/lock/kuranai-deploy.lock"
readonly AAPANEL_WELL_KNOWN_DIR="/www/wwwroot/kuranai.talhacan.com/.well-known"
readonly AAPANEL_VHOST="/www/server/panel/vhost/nginx/kuranai.talhacan.com.conf"
readonly AAPANEL_NGINX="/www/server/nginx/sbin/nginx"
readonly AASA_INCLUDE="/www/wwwroot/kuranai/Kuran_AI_Infrastructures/aapanel/kuranai-well-known.conf"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$INFRA_DIR/.." && pwd)"
K8S_DIR="$INFRA_DIR/k8s"

if [[ "$PROJECT_ROOT" != "$EXPECTED_ROOT" ]]; then
  echo "Refusing deployment outside $EXPECTED_ROOT (found $PROJECT_ROOT)" >&2
  exit 1
fi

for command_name in docker git kubectl curl flock; do
  command -v "$command_name" >/dev/null || { echo "Missing command: $command_name" >&2; exit 1; }
done

exec 9>"$LOCK_FILE"
flock -n 9 || { echo "Another Kuran AI deployment is already running" >&2; exit 1; }

diagnostics() {
  local exit_code=$?
  if (( exit_code != 0 )); then
    echo "Deployment failed; current workload state:" >&2
    kubectl get pods,deployments,services -n "$NAMESPACE" >&2 || true
  fi
  exit "$exit_code"
}
trap diagnostics EXIT

cd "$PROJECT_ROOT"
if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
  echo "Tracked production files contain local changes; refusing to overwrite them" >&2
  exit 1
fi

git fetch --prune origin main
git checkout main
git pull --ff-only origin main

docker build -t kuranai-gateway:latest "$PROJECT_ROOT/Kuran_AI_Gateway"
docker build \
  --build-arg VITE_API_BASE=https://kuranai-api.talhacan.com \
  --build-arg VITE_GOOGLE_CLIENT_ID=157832066999-3tgb2daavmrd1d1vd10ptrato1iu382t.apps.googleusercontent.com \
  --build-arg VITE_APPLE_CLIENT_ID=com.talhacan.kuranai.web \
  --build-arg VITE_APPLE_REDIRECT_URI=https://kuranai.talhacan.com/auth/apple/callback \
  -t kuranai-web:latest "$PROJECT_ROOT/Kuran_AI_Web"

kubectl apply -f "$K8S_DIR/namespace.yaml"
kubectl apply -f "$K8S_DIR/config.yaml"
kubectl apply -f "$K8S_DIR/data.yaml"
kubectl rollout status statefulset/postgres -n "$NAMESPACE" --timeout=240s
kubectl rollout status statefulset/redis -n "$NAMESPACE" --timeout=240s

# Existing Kubernetes secrets and PVCs are intentionally preserved.
kubectl delete job gateway-migrate -n "$NAMESPACE" --ignore-not-found
kubectl apply -f "$K8S_DIR/migration.yaml"
kubectl wait --for=condition=complete job/gateway-migrate -n "$NAMESPACE" --timeout=240s

kubectl apply -f "$K8S_DIR/apps.yaml"
kubectl rollout restart deployment/gateway deployment/web -n "$NAMESPACE"
kubectl rollout status deployment/gateway -n "$NAMESPACE" --timeout=300s
kubectl rollout status deployment/web -n "$NAMESPACE" --timeout=180s

gateway_ready="$(kubectl get deployment/gateway -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')"
web_ready="$(kubectl get deployment/web -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')"
[[ "$gateway_ready" == "2" ]] || { echo "Expected 2 ready gateway pods, found $gateway_ready" >&2; exit 1; }
[[ "$web_ready" == "1" ]] || { echo "Expected 1 ready web pod, found $web_ready" >&2; exit 1; }

curl --fail --silent --show-error --retry 10 --retry-delay 2 http://127.0.0.1:7008/health >/dev/null
curl --fail --silent --show-error --retry 10 --retry-delay 2 http://127.0.0.1:7007/ >/dev/null

# aaPanel serves /.well-known from its own document root before the proxy project.
# Keep Apple's association file there in addition to the copy in the web image.
install -d -m 755 "$AAPANEL_WELL_KNOWN_DIR"
install -m 644 \
  "$PROJECT_ROOT/Kuran_AI_Web/public/.well-known/apple-app-site-association" \
  "$AAPANEL_WELL_KNOWN_DIR/apple-app-site-association"
if ! grep -Fq "include $AASA_INCLUDE;" "$AAPANEL_VHOST"; then
  sed -i "/#SERVER-BLOCK START/a\\    include $AASA_INCLUDE;" "$AAPANEL_VHOST"
fi
"$AAPANEL_NGINX" -t
"$AAPANEL_NGINX" -s reload
curl --fail --silent --show-error --retry 10 --retry-delay 2 \
  --resolve kuranai.talhacan.com:443:127.0.0.1 \
  https://kuranai.talhacan.com/.well-known/apple-app-site-association >/dev/null
aasa_content_type="$(curl --fail --silent --show-error \
  --resolve kuranai.talhacan.com:443:127.0.0.1 \
  --output /dev/null --write-out '%{content_type}' \
  https://kuranai.talhacan.com/.well-known/apple-app-site-association)"
[[ "$aasa_content_type" == application/json* ]] || {
  echo "Unexpected AASA content type: $aasa_content_type" >&2
  exit 1
}

echo "Kuran AI deployed successfully at commit $(git rev-parse --short HEAD)"
