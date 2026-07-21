# Kuran AI Gateway

Go/Gin tabanlı API gateway. PostgreSQL (PgVector) + Redis ile çalışır.

## Özellikler

- **Auth**: Email/password, Google OAuth, Apple Sign in
- **JWT**: Access + Refresh token, Redis'te cache
- **Tablolar**: users, managers, jwt_tokens, meal_embeddings, surah_exegesis_embeddings, ayah_exegesis_embeddings
- **Scalar API Docs**: `/docs` endpoint'inde

## Lokal çalıştırma

```bash
# .env dosyası gerekli (DB_HOST, REDIS_HOST vb.)
go run main.go
```

## Docker ile

Kuran_AI_Infrastructures klasöründen `docker compose up -d --build` ile tüm sistem ayağa kalkar.

## API Endpoints

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/user/auth/register | Kayıt |
| POST | /api/user/auth/login | Giriş (email/pass) |
| POST | /api/user/auth/login/google | Google ile giriş |
| POST | /api/user/auth/login/apple | Apple ile giriş |
| POST | /api/user/auth/refresh | Token yenile |
| GET | /api/user/auth/me | Mevcut kullanıcı (auth gerekli) |
| POST | /api/user/auth/logout | Çıkış (auth gerekli) |

## Migration

```bash
./gateway migrate
```
