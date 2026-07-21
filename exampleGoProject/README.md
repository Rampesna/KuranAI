# Gateway (Go / Gin)

Laravel benzeri mimaride backend API gateway. PostgreSQL, pgvector, Redis, JWT (access + refresh) kullanır.

## Yapı

- **routes/** — `user.go`, `manager.go` (user ve manager API rotaları)
- **handlers/user/** — `auth_handler.go`, `user_handler.go`, `todo_handler.go`
- **handlers/manager/** — `auth_handler.go`, `manager_handler.go`, `todo_handler.go`
- **models/** — `user.go`, `manager.go`, `todo.go`, `personal_access_token.go`
- **core/** — `service_response.go` (servis cevap standardı)
- **services/** — iş mantığı; handler request → DTO → service → ServiceResponse
- **middleware/** — JWT auth; token doğrulama, user/manager context
- **config/** — ortam değişkenleri
- **database/** — GORM bağlantı ve AutoMigrate (users, managers, todos, personal_access_tokens)
- **redis/** — JWT token cache (doğrulama için)

## Gereksinimler

- Go 1.21+
- PostgreSQL (pgvector extension)
- Redis

## Çalıştırma

```bash
cp .env.example .env
# .env içinde DB ve Redis bilgilerini düzenle

go run main.go
```

API: `http://localhost:8080`

- Health: `GET /health`
- User login: `POST /api/user/auth/login` — body: `{"credential":"email|username","password":"..."}`
- User refresh: `POST /api/user/auth/refresh` — body: `{"refreshToken":"..."}`
- User me: `GET /api/user/auth/me` — Header: `Authorization: Bearer <accessToken>`
- Manager login: `POST /api/manager/auth/login`
- Manager me: `GET /api/manager/auth/me` — Header: `Authorization: Bearer <accessToken>`

## JWT ve Cihaz Bilgisi

- Login’de access + refresh token üretilir; hem `personal_access_tokens` tablosuna hem Redis’e yazılır.
- Tabloda: `user_agent`, `ip_address`, `device_info`, `device_type` alanları cihaz/session takibi için kullanılır.
- API isteklerinde token Redis'ten kontrol edilir; geçerliyse context'e user/manager eklenir.

---

## Veritabanını sıfırlama / baştan oluşturma

Tüm tabloları silip migration'ları tekrar çalıştırmak (temiz DB) için:

### Docker ile (volume dahil tam sıfır)

```bash
# Proje kökünden (talhacan/)
docker compose down -v
docker compose up -d --build
```

`-v` volume'leri de sildiği için PostgreSQL verisi tamamen silinir; gateway ayağa kalkınca AutoMigrate tabloları yeniden oluşturur.

### Sadece tabloları drop edip yeniden migrate (volume aynen kalır)

**Docker:**

```bash
docker compose down
docker compose up -d postgres redis
# Postgres hazır olana kadar birkaç saniye bekle, sonra:
DB_RESET=1 docker compose run --rm gateway
docker compose up -d
```

**Yerel (Postgres/Redis zaten çalışıyorsa):**

```bash
DB_RESET=1 go run main.go
```

Bir kez çalıştır; uygulama tabloları drop edip tekrar oluşturur ve normal şekilde ayağa kalkar. Bir sonraki çalıştırmada `DB_RESET=1` vermezsen mevcut veri korunur.
