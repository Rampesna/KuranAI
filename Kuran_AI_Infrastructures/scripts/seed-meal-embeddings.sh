#!/bin/bash
# meal_embeddings tablosuna kuran_embeddings.sql verilerini aktarır.
# Kuran_embeddings.sql 'public.embeddings' kullanıyor; biz 'meal_embeddings' kullanıyoruz.
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$INFRA_DIR/.." && pwd)"
SQL_FILE="$PROJECT_ROOT/DB_Backup/kuran_embeddings.sql"
if [ ! -f "$SQL_FILE" ]; then
  echo "Error: $SQL_FILE not found"
  exit 1
fi
echo "Seeding meal_embeddings from kuran_embeddings.sql (this may take a while)..."
sed 's/public\.embeddings/public.meal_embeddings/g' "$SQL_FILE" | docker exec -i kuranai_postgres psql -U kuranai_user -d kuranai_pgsql -q
echo "Done."
