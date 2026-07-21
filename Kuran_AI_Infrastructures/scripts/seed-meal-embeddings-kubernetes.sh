#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$INFRA_DIR/.." && pwd)"
SQL_FILE="$PROJECT_ROOT/DB_Backup/kuran_embeddings.sql"
NAMESPACE="kuranai"

if [[ ! -f "$SQL_FILE" ]]; then
  echo "Missing $SQL_FILE"
  exit 1
fi

POSTGRES_USER="$(kubectl get secret kuranai-secrets -n "$NAMESPACE" -o jsonpath='{.data.POSTGRES_USER}' | base64 -d)"
POSTGRES_DB="$(kubectl get secret kuranai-secrets -n "$NAMESPACE" -o jsonpath='{.data.POSTGRES_DB}' | base64 -d)"

ROW_COUNT="$(kubectl exec -n "$NAMESPACE" postgres-0 -- psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Atc 'SELECT COUNT(*) FROM meal_embeddings')"
if [[ "$ROW_COUNT" != "0" ]]; then
  echo "meal_embeddings already contains $ROW_COUNT rows; seed skipped"
  exit 0
fi

echo "Importing meal embeddings..."
sed 's/public\.embeddings/public.meal_embeddings/g' "$SQL_FILE" | \
  kubectl exec -i -n "$NAMESPACE" postgres-0 -- psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -q

kubectl exec -n "$NAMESPACE" postgres-0 -- psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -c \
  "SELECT setval(pg_get_serial_sequence('meal_embeddings','id'), COALESCE((SELECT MAX(id) FROM meal_embeddings), 1), true);"

kubectl exec -n "$NAMESPACE" postgres-0 -- psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Atc 'SELECT COUNT(*) FROM meal_embeddings'

