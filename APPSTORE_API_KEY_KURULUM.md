# App Store Connect API Key - Issuer ID ve Key ID Nereden Alınır?

## 1. App Store Connect'e girin

- [appstoreconnect.apple.com](https://appstoreconnect.apple.com) adresine gidin
- Apple Developer hesabınızla giriş yapın

## 2. Issuer ID'yi bulun

1. **Users and Access** sekmesine tıklayın (üst menüde)
2. Sol menüden **Integrations** → **App Store Connect API** seçin
3. **App Store Connect API Keys** bölümünde, sağ üstte **Issuer ID** gösterilir
   - Örnek format: `57246542-96fe-1a63-e053-0824d011072a`
   - Bu UUID formatındaki değeri kopyalayın

![Issuer ID konumu: Users and Access → Integrations → App Store Connect API]

## 3. Key ID'yi bulun (veya yeni key oluşturun)

1. Aynı sayfada (**Integrations** → **App Store Connect API**) aşağı kaydırın
2. **Keys** tablosunda mevcut key'ler listelenir
3. **Yeni key yoksa:**
   - **Generate API Key** butonuna tıklayın
   - **Key Name** girin (örn: "Kuran AI Backend")
   - **Access** için **App Manager** veya ihtiyacınız olan rolü seçin
   - **Generate** tıklayın
4. Oluşan key satırında **Key ID** sütunu görünür
   - Örnek: `ABC123XYZ9` (10 karakterlik alfanumerik)
   - **Download API Key** ile `.p8` dosyasını indirin (sadece 1 kez indirilebilir!)

## 4. .p8 dosyasını saklayın

- İndirilen dosya: `AuthKey_KEYID.p8` (örn: `AuthKey_ABC123XYZ9.p8`)
- Bu dosyayı proje dizinine veya güvenli bir yere koyun
- **`.gitignore`'a mutlaka ekleyin** (gizli anahtar!)

## 5. Script'i çalıştırın

```bash
# Gerekli paket
pip install PyJWT cryptography

# appstore_api_key.py içinde ISSUER_ID, KEY_ID ve PRIVATE_KEY_PATH değerlerini doldurun
# veya environment variable ile:
export APPSTORE_ISSUER_ID="57246542-96fe-1a63-e053-0824d011072a"
export APPSTORE_KEY_ID="ABC123XYZ9"
export APPSTORE_KEY_PATH="./AuthKey_ABC123XYZ9.p8"

python appstore_api_key.py
```

Çıktıda JWT token görünecektir. Bu token'ı App Store Connect API isteklerinde `Authorization: Bearer <token>` header'ı ile kullanırsınız.
