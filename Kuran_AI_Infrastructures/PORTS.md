# Port Dağılımı

## Kubernetes Production

| Port | Servis | Domain |
|------|--------|--------|
| 7007 | React web NodePort | `kuranai.talhacan.com` |
| 7008 | Gateway NodePort | `kuranai-api.talhacan.com` |

Her iki servis `kuranai` namespace'indedir. PostgreSQL ve Redis yalnızca ClusterIP ile küme içinden erişilir.

## Docker Compose Development (9300–9400)

Yeni servis eklerken bu aralıktan uygun port seçin.

| Port | Servis | Açıklama |
|------|--------|----------|
| 9300 | Gateway | Go/Gin API, Scalar docs |
| 9301 | PostgreSQL | PgVector |
| 9302 | Redis | JWT cache |
| 9303 | React web | Vite production build |
| 9304 | - | Boş |
| ... | ... | ... |
| 9400 | - | Boş |
