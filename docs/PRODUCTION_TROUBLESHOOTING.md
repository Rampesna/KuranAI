# Production Troubleshooting – iOS "Network Connection Lost" (-1005)

## Hata
```
NSURLErrorDomain Code=-1005 "The network connection was lost."
_kCFStreamErrorCodeKey=57
URL: https://api.kuranaichat.rampesna.com/api/user/auth/login/google
```

Bu hata genellikle bağlantı kurulduktan sonra sunucu veya aradaki bileşen (proxy, load balancer) tarafından kesildiğinde oluşur.

---

## 1. SSL Sertifikası Kontrolü

En sık sebep **eksik sertifika zinciri** (intermediate certificates). iOS sertifika zincirini doğrulayamazsa bağlantıyı keser.

### Kontrol
```bash
openssl s_client -connect api.kuranaichat.rampesna.com:443 -servername api.kuranaichat.rampesna.com
```

Çıktıda `Verify return code: 0 (ok)` görmelisin. `verify error` varsa zincir eksik.

### SSL Labs ile test
https://www.ssllabs.com/ssltest/analyze.html?d=api.kuranaichat.rampesna.com

- **Chain issues** veya **Incomplete** uyarısı olmamalı.

### Çözüm (Nginx + Let's Encrypt)
`fullchain.pem` kullanıldığından emin ol (sadece `cert.pem` değil):
```nginx
ssl_certificate     /etc/letsencrypt/live/api.kuranaichat.rampesna.com/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/api.kuranaichat.rampesna.com/privkey.pem;
```

---

## 2. Nginx Reverse Proxy Ayarları

Timeout’lar kısa olursa (özellikle Google/Apple login gibi dış API çağrıları uzun sürebilir) bağlantı kesilir.

### Önerilen Nginx Config
```nginx
server {
    listen 443 ssl http2;
    server_name api.kuranaichat.rampesna.com;

    ssl_certificate     /path/to/fullchain.pem;
    ssl_certificate_key /path/to/privkey.pem;

    # SSL ayarları
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;

    location / {
        proxy_pass http://127.0.0.1:9300;  # veya gateway container IP
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Kritik: Uzun süren istekler (Google OAuth, chat SSE, TTS)
        proxy_connect_timeout 60s;
        proxy_send_timeout 120s;
        proxy_read_timeout 120s;

        # Buffer ayarları – büyük yanıtlar için
        proxy_buffering off;
        proxy_buffer_size 4k;

        # SSE için gerekli
        proxy_set_header Connection '';
        chunked_transfer_encoding off;
    }
}
```

> `proxy_read_timeout 120s` Google tokeninfo çağrısı için önemli. Varsayılan 60s bazen yetmeyebilir.

---

## 3. CORS (Mobil için gerekli değil ama web varsa)

iOS native uygulaması CORS’tan etkilenmez. Web sürümü varsa Gateway CORS ayarlarını kontrol et. `CORS_ORIGINS` env değişkeni ile domain eklenebilir.

---

## 4. Sunucu Tarafı Kontroller

### Gateway logları
İstek sunucuya ulaşıyor mu?
```bash
docker logs -f kuranai_gateway
# veya
journalctl -u kuranai-gateway -f
```

### cURL ile test
```bash
# Health check
curl -v https://api.kuranaichat.rampesna.com/health

# Google login (token olmadan 400 beklenir, ama bağlantı kopmamalı)
curl -v -X POST https://api.kuranaichat.rampesna.com/api/user/auth/login/google \
  -H "Content-Type: application/json" \
  -d '{"idToken":"test"}'
```

Bağlantı `curl` ile çalışıyorsa, sorun büyük ihtimalle:
- SSL zinciri (iOS daha katı)
- Mobil ağ / firewall

---

## 5. Cloudflare veya CDN Kullanıyorsan

- **SSL/TLS** modu: Full (strict) veya Full
- **Minimum TLS**: 1.2
- **Proxy timeout**: Arttır
- **Rocket Loader** gibi özellikler API için kapatılmalı

---

## 6. iOS Tarafı

### Timeout
`APIEnvironment.timeout` 30 saniye. Login için yeterli olmalı; sorun genelde sunucu/proxy tarafında.

### ATS (App Transport Security)
HTTPS kullanıyorsun, ek ayar gerekmez. `Info.plist`’te `NSAppTransportSecurity` ile HTTP’ye izin verme (güvenlik riski).

---

## 7. Öncelik Sırası

1. **SSL zinciri** – `fullchain.pem` kullan, SSL Labs’ta kontrol et
2. **Nginx timeouts** – `proxy_read_timeout 120s` (ve benzeri) ayarla
3. **Gateway logları** – İstek geliyor mu, hata var mı bak
4. **cURL test** – Sunucuya erişim ve SSL doğrulaması

---

## Hızlı Nginx Düzeltmesi

Mevcut `location /` bloğuna ekle:
```nginx
proxy_connect_timeout 60s;
proxy_send_timeout 120s;
proxy_read_timeout 120s;
proxy_buffering off;
```

Sonra:
```bash
sudo nginx -t
sudo systemctl reload nginx
```
