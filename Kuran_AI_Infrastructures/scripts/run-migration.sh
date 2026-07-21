#!/bin/bash
# Gateway container'da migration çalıştırır (tabloları oluşturur/günceller).
set -e
docker exec -it kuranai_gateway /app/gateway migrate
echo "Migration done."
