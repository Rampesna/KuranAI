# Kuran AI Infrastructure

## Production: Kubernetes (k3s)

Production kurulumu diğer projelerden ayrılmış `kuranai` namespace'inde çalışır:

- Gateway: 2 pod, NodePort `7008`, `https://kuranai-api.talhacan.com`
- Web: 1 pod, NodePort `7007`, `https://kuranai.talhacan.com`
- PostgreSQL + pgvector: 1 StatefulSet ve kalıcı disk
- Redis: 1 StatefulSet ve kalıcı disk
- Avatar yüklemeleri: kalıcı disk

Sunucuda proje kökünden dağıtmak için:

```bash
cd /www/wwwroot/kuranai
Kuran_AI_Infrastructures/scripts/deploy-kubernetes.sh
```

Script OpenAI anahtarını gizli olarak ister, `.env` içindeki diğer değerlerden Kubernetes Secret üretir, image'ları yerelde build eder, migration Job'unu çalıştırır ve rollout tamamlanana kadar bekler. Gerçek anahtarlar `k8s/` manifestlerine yazılmaz.

İlk kurulumda meal embedding yedeğini bir kez yüklemek için:

```bash
Kuran_AI_Infrastructures/scripts/seed-meal-embeddings-kubernetes.sh
```

Script tablo doluysa tekrar import etmez. Durumu kontrol etmek için:

```bash
kubectl get pods,svc,pvc -n kuranai
kubectl rollout status deployment/gateway -n kuranai
kubectl rollout status deployment/web -n kuranai
```

## Local Development: Docker Compose

```bash
cd Kuran_AI_Infrastructures
docker compose up -d --build
```

PostgreSQL (PgVector), Redis, Gateway ve React web servisleri ayağa kalkar. **Veritabanı otomatik sıfırlanmaz** — migration gateway başlarken çalışır, seed işlemini istediğiniz zaman manuel yaparsınız.

## Migration ve Seed Komutları

### 1. Migration (Tabloları oluştur / güncelle)

Migration, gateway başlarken otomatik çalışır. İsterseniz manuel tetiklemek için:

```bash
docker exec -it kuranai_gateway /app/gateway migrate
```

### 2. Seed (meal_embeddings verilerini yükle)

`kuran_embeddings.sql` dosyasındaki verileri `meal_embeddings` tablosuna aktarır:

```bash
cd Kuran_AI_Infrastructures
chmod +x scripts/seed-meal-embeddings.sh
./scripts/seed-meal-embeddings.sh
```

veya direkt:

```bash
sed 's/public\.embeddings/public.meal_embeddings/g' ../DB_Backup/kuran_embeddings.sql | docker exec -i kuranai_postgres psql -U kuranai_user -d kuranai_pgsql -q
```

### 3. Veritabanını sıfırlama (refresh)

Tüm verileri silip migration’ı yeniden çalıştırmak için:

```bash
cd Kuran_AI_Infrastructures

# Schema'yı drop et ve vector extension'ı yeniden oluştur
docker exec -it kuranai_postgres psql -U kuranai_user -d kuranai_pgsql -c "
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO kuranai_user;
CREATE EXTENSION IF NOT EXISTS vector;
"

# Gateway'i yeniden başlat (migration otomatik çalışır)
docker compose restart gateway

# Ardından seed çalıştır (isterseniz)
./scripts/seed-meal-embeddings.sh
```

## Gateway’i yeniden build etme

Go kodunda değişiklik yaptıktan sonra imajı yeniden build etmek için:

```bash
cd Kuran_AI_Infrastructures
docker compose up -d --build gateway
```

## iOS uygulaması – fiziksel cihazda backend bağlantısı

Uygulama **simulator**da çalışırken `localhost:9300` Mac’e gider; **fiziksel iPhone**da çalışırken `localhost` telefona işaret eder, bu yüzden bağlantı reddedilir.

Fiziksel cihazda test için:

1. Mac’in yerel ağ IP’sini öğren: **Sistem Ayarları → Ağ** veya terminalde `ipconfig getifaddr en0`.
2. iOS projesinde `KuranAIChat/Networking/APIEnvironment.swift` içinde `developmentBaseURL` değerini bu IP ile güncelle, örn: `"http://192.168.1.5:9300"`.
3. iPhone ve Mac aynı Wi‑Fi ağında olmalı.

## Scalar API Dokümantasyonu

Gateway çalışırken: **http://localhost:9300/docs** adresinden Scalar API Reference'a erişebilirsiniz.

## Port Dağılımı

| Ortam | Port | Servis | Açıklama |
|---|---:|---|---|
| Kubernetes | 7007 | Web | aaPanel proxy hedefi |
| Kubernetes | 7008 | Gateway | aaPanel API proxy hedefi |
| Docker Compose | 9300 | Gateway | Yerel API ve Scalar docs |
| Docker Compose | 9301 | PostgreSQL | Yerel pgvector |
| Docker Compose | 9302 | Redis | Yerel JWT cache |
| Docker Compose | 9303 | Web | Yerel React site |
