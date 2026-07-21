# Kuran AI

Kuran AI; Kur'an meali, tefsir ve hadis kaynakları üzerinde yapay zekâ destekli sohbet sunan iOS ve web uygulamasıdır.

## Bileşenler

- `Kuran_AI_Gateway`: Go API, kimlik doğrulama, sohbet, paketler ve abonelik doğrulama.
- `Kuran_AI_Web`: Türkçe varsayılan, Türkçe/İngilizce seçilebilir React web uygulaması.
- `Kuran_AI_iOS/KuranAIChat`: SwiftUI iOS uygulaması.
- `Kuran_AI_Infrastructures`: Docker ve Kubernetes üretim tanımları.
- `Python_Embedder` ve `Web_Scrapper`: Kaynak verisi hazırlama araçları.

## Yerel geliştirme

Web uygulaması:

```bash
cd Kuran_AI_Web
npm ci
npm run dev
```

API:

```bash
cd Kuran_AI_Gateway
go test ./...
go run .
```

Gerekli gizli değerleri sürüm kontrolüne eklemeyin. Örnek değişkenler `Kuran_AI_Infrastructures/.env.example` dosyasındadır.

## Üretim

Üretim kaynakları `kuranai` Kubernetes namespace'inde birbirinden izole çalışır:

- API: 2 pod, NodePort `7008`, `kuranai-api.talhacan.com`
- Web: 1 pod, NodePort `7007`, `kuranai.talhacan.com`
- PostgreSQL ve Redis: kalıcı diskli StatefulSet

`main` dalına gönderilen her değişiklik GitHub Actions üzerinde test edilir. Başarılı sonuçtan sonra sunucudaki salt-okunur Git checkout güncellenir, imajlar yeniden oluşturulur, migration çalışır ve rolling rollout doğrulanır. Kubernetes secret'ları ile kalıcı veriler dağıtım sırasında değiştirilmez.

