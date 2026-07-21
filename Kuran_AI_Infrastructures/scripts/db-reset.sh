#!/bin/bash
# Veritabanını sıfırlar: schema drop, extension tekrar, sonra migration + seed.
set -e
echo "Resetting database..."
docker exec -it kuranai_postgres psql -U kuranai_user -d kuranai_pgsql -c "
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO kuranai_user;
CREATE EXTENSION IF NOT EXISTS vector;
"
echo "Database reset. Restart gateway to run migrations, then run seed-meal-embeddings.sh"
cd "$(cd "$(dirname "$0")/.." && pwd)" && docker compose restart gateway
