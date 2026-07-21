# KuranAI – App Store Yayınlama Kılavuzu

Bu dosya, KuranAI uygulamasının App Store Connect üzerinden yayınlanması için gerekli tüm metinleri, izin açıklamalarını ve inceleme notlarını içerir.

---

## 1. Temel Bilgiler

| Alan | Değer |
|------|-------|
| **Uygulama Adı** | KuranAI |
| **Bundle ID** | talhacan.KuranAIChat |
| **Sürüm** | 1.0 |
| **Telif** | © 2026 [Şirket/İsim] |

---

## 2. Promotional Text (Tanıtım Metni) – Max 170 karakter

Güncellemelerde sürüm değişikliği gerektirmez. En üstte görünür.

### Türkçe
```
Kur'an meal, tefsir ve hadis hakkında yapay zeka ile sohbet edin. Sesli soru sorun, cevapları dinleyin. Namaz vakitleri ve günün ayeti ile gününüzü anlamlı kılın.
```

### English
```
Chat with AI about Quran, tafsir and hadith. Ask questions by voice, listen to answers. Enrich your day with prayer times and the verse of the day.
```

---

## 3. Description (Açıklama) – Max 4.000 karakter

SEO ve kullanıcı bilgilendirmesi için kritik. Anahtar kelimeleri doğal biçimde kullanın.

### Türkçe
```
KuranAI, Kur'an-ı Kerim meal, tefsir ve hadis konularında size yardımcı olmak için yapay zeka destekli bir sohbet uygulamasıdır.

ÖZELLİKLER

• Meal, Tefsir ve Hadis Sohbeti
Kur'an ayetlerinin meal ve tefsirini, sahih hadisleri sorgulayabilirsiniz. Yapay zeka, Kur'an ve Sünnet bilgisine dayalı cevaplar sunar.

• Sesli Soru Sorma
Soru sormak istediğinizde mikrofonu kullanarak doğrudan konuşabilirsiniz. Sesiniz metne çevrilir ve cevap alırsınız.

• Sesli Okuma (TTS)
Uzun cevapları okumak yerine dinleyebilirsiniz. Özel sesler ile metinleri sesli olarak dinleyin.

• Namaz Vakitleri
Konumunuza göre günlük namaz vakitlerini görüntüleyin. Vakitleri kaçırmayın.

• Günün Ayeti
Her gün farklı bir ayet ile gününüze anlam katın.

• Kişisel Profil
Google veya Apple ile güvenli giriş yapın, profil fotoğrafınızı ekleyin.

ABONELİK PAKETLERİ

Ücretsiz kullanımda günde 3 soru hakkı vardır. Sesli okuma için abonelik gereklidir. Bronze, Silver, Gold ve Diamond paketleri ile kullanım limitlerinizi artırın.

KuranAI, Kur'an öğrenmek, meal ve tefsir okumak, hadis araştırmak isteyen herkes için tasarlanmıştır.
```

### English
```
KuranAI is an AI-powered chat app to help you learn and explore the Quran, tafsir, and hadith.

FEATURES

• Quran, Tafsir & Hadith Chat
Ask about Quranic verses, their meanings and interpretations, and authentic hadith. AI provides answers based on Quran and Sunnah knowledge.

• Voice Questions
Use the microphone to ask questions by speaking. Your voice is transcribed and you get an answer.

• Text-to-Speech
Listen to long answers instead of reading. Hear responses in natural voices.

• Prayer Times
View daily prayer times based on your location. Never miss a prayer.

• Verse of the Day
Start each day with a different verse from the Quran.

• Personal Profile
Sign in securely with Google or Apple, add your profile photo.

SUBSCRIPTION PACKAGES

Free users get 3 questions per day. Text-to-speech requires a subscription. Bronze, Silver, Gold and Diamond packages increase your daily limits.

KuranAI is designed for anyone who wants to study the Quran, read translations and tafsir, and explore hadith.
```

---

## 4. Keywords – Max 100 karakter

Virgülle ayrılmış, boşluksuz. Tekrarlar kabul edilmez. En etkili kelimeler önce yazılmalı.

### Türkçe (virgül, boşluksuz)
```
Kuran,Kur'an,meal,tefsir,hadis,yapay zeka,sesli soru,namaz vakti,ayet,İslam,Kur'an meal
```

### English (virgül, boşluksuz)
```
Quran,Koran,tafsir,hadith,AI,voice,Muslim,prayer times,ayah,Islam,Quran translation
```

> **ASO İpucu:** "Kuran", "meal", "tefsir", "hadis", "namaz" gibi Türkçe aramalarda çıkması için Türkçe lokalizasyonda kullanın. İngilizce için "Quran", "tafsir", "hadith", "prayer" önemli.

---

## 5. URL’ler

Gateway’de aşağıdaki sayfalar yayında (API base URL’iniz üzerinden erişilir):

| Endpoint | Açıklama | App Store Connect’te kullanım |
|----------|----------|-------------------------------|
| **GET /support** | Destek (Türkçe) – Turkish support URL şartı için | **Support URL** → `https://[domain]/support` |
| **GET /privacy-policy** | Gizlilik politikası | **Privacy Policy URL** → `https://[domain]/privacy-policy` |
| **GET /terms** | Kullanım koşulları | İsteğe bağlı / mağaza açıklaması |
| **GET /contact** | İletişim | İsteğe bağlı |

| Alan | Önerilen Değer |
|------|----------------|
| **Support URL** | `https://[domain]/support` (Türkçe destek sayfası – kesin şart) |
| **Privacy Policy URL** | `https://[domain]/privacy-policy` |
| **Marketing URL** | `https://[domain]` (opsiyonel) |

> `[domain]` = API’nin yayında olduğu alan adı. Örnek: `https://api.kuranai.app/support`

---

## 6. Copyright

```
© 2026 [Şirket veya Geliştirici Adı]
```

---

## 7. İzinler (Permissions) ve Review Notları

### 7.1 Mikrofon (NSMicrophoneUsageDescription)

**Info.plist metni:**
> Kuran AI sohbet için sesli soru sormak üzere mikrofonu kullanır.

**Review ekibi için not:**
> The microphone is used to record voice questions. When the user taps the microphone button in the chat input, they can speak their question. The audio is sent to our backend, transcribed via Whisper API, and the result is used as the chat message. No audio is stored permanently.

### 7.2 Konum (NSLocationWhenInUseUsageDescription)

**Info.plist metni:**
> Namaz vakitlerini göstermek için konumunuz kullanılır.

**Review ekibi için not:**
> Location is used only to calculate daily prayer times based on the user's geographic position. If location is denied or unavailable, we fall back to Istanbul coordinates. We do not store or share location data.

### 7.3 Fotoğraf Kütüphanesi

Profil fotoğrafı için `PhotosPicker` kullanılıyor. iOS 16+ ile sınırlı erişim varsayılandır; ek bir Info.plist anahtarı gerekmez. Tam erişim kullanıyorsanız:

**NSPhotoLibraryUsageDescription:**
> Profil fotoğrafı eklemek için fotoğraf kütüphanesine erişim gereklidir.

### 7.4 Routing App Coverage File

Uygulama navigasyon/harita sunmuyor. **Dosya yüklemeniz gerekmez.**

---

## 8. App Store Connect – App Review Information

**Contact Information**
- First Name, Last Name, Phone, Email: İnceleme ekibinin size ulaşacağı bilgiler

**Demo Account (Test hesabı)**
- E-posta ve şifre (veya “Sign in with Apple/Google” ile giriş yapılabildiğini belirtin)
- Opsiyonel: Sandbox test kullanıcısı bilgisi

**Notes**
```
DEMO ACCOUNT:
- Sign in with Apple or Google to access the app.
- Free users: 3 questions/day, no TTS. To test TTS and higher limits, use a Sandbox test account to purchase a subscription (Settings > Sandbox Account on device).

PERMISSIONS:
- Microphone: Used only for voice input when user taps the mic button in chat. Transcribed to text and sent as message.
- Location: Used only for prayer times on the home screen. Fallback to Istanbul if denied.

IN-APP PURCHASES:
- Subscriptions (Bronze, Silver, Gold, Diamond) are available. Use Sandbox account to test purchase flow.
```

---

## 9. App Privacy (Gizlilik)

App Store Connect → App Privacy bölümünde belirtmeniz gerekenler:

| Veri Türü | Toplanıyor mu? | Amaç |
|-----------|----------------|------|
| Email | Evet | Hesap oluşturma (Google/Apple) |
| Ad, Soyad | Evet | Profil (Google/Apple’dan) |
| Profil fotoğrafı | Evet | Profil (kullanıcı seçimi) |
| Konum (yaklaşık) | Evet | Namaz vakitleri hesaplama |
| Ses kaydı | Geçici | Sesli soru → metin dönüşümü, saklanmaz |
| Sohbet mesajları | Evet | Sohbet geçmişi, hizmet sunumu |

**Veri Kullanımı:** Hesap yönetimi, hizmet sunumu, kişiselleştirme (namaz vakitleri). Üçüncü taraflarla paylaşım: Ödeme işlemleri (Apple), kimlik doğrulama (Google/Apple).

---

## 10. Ekran Görüntüleri (Screenshots)

**iPhone 6.5" (iPhone 14 Pro Max vb.):** 1284 x 2778 px  
**İlk 3 ekran** App Store kurulum sayfasında gösterilir.

Önerilen sıra:
1. Ana ekran / Vahy ile konuş
2. Sohbet ekranı (örnek soru–cevap)
3. Namaz vakitleri veya Günün ayeti kartı
4. Profil / Abonelik paketleri

---

## 11. In-App Purchase Açıklamaları

Her ürün için App Store Connect’te gösterilecek metin:

| Product ID | Açıklama (TR) | Açıklama (EN) |
|------------|---------------|---------------|
| kuranai.bronze.monthly | Bronze Paket – Aylık. 100 soru/gün, 5 sesli okuma/gün | Bronze Package – Monthly. 100 questions/day, 5 TTS/day |
| kuranai.bronze.yearly | Bronze Paket – Yıllık. 100 soru/gün, 5 sesli okuma/gün | Bronze Package – Yearly. 100 questions/day, 5 TTS/day |
| kuranai.silver.monthly | Silver Paket – Aylık. 250 soru/gün, 50 sesli okuma/gün | Silver Package – Monthly. 250 questions/day, 50 TTS/day |
| kuranai.silver.yearly | Silver Paket – Yıllık. 250 soru/gün, 50 sesli okuma/gün | Silver Package – Yearly. 250 questions/day, 50 TTS/day |
| kuranai.gold.monthly | Gold Paket – Aylık. Sınırsız soru, 100 sesli okuma/gün | Gold Package – Monthly. Unlimited questions, 100 TTS/day |
| kuranai.gold.yearly | Gold Paket – Yıllık. Sınırsız soru, 100 sesli okuma/gün | Gold Package – Yearly. Unlimited questions, 100 TTS/day |
| kuranai.diamond.monthly | Diamond Paket – Aylık. Sınırsız soru ve sesli okuma | Diamond Package – Monthly. Unlimited questions and TTS |
| kuranai.diamond.yearly | Diamond Paket – Yıllık. Sınırsız soru ve sesli okuma | Diamond Package – Yearly. Unlimited questions and TTS |

---

## 12. Kontrol Listesi (Checklist)

- [ ] Agreements, Tax, and Banking tamamlandı
- [ ] Paid Applications sözleşmesi imzalı
- [ ] In-App Purchase ürünleri oluşturuldu (Product ID’ler eşleşiyor)
- [ ] Support URL gerçek ve çalışır
- [ ] App Privacy dolduruldu
- [ ] Ekran görüntüleri yüklendi (en az 3)
- [ ] Promotional Text ve Description (TR/EN) girildi
- [ ] Keywords girildi
- [ ] App Review Information ve Demo Account eklendi
- [ ] Info.plist izin metinleri güncel

---

## 13. Sık Karşılaşılan Red Sebepleri ve Önlemler

1. **Incomplete Metadata:** Tüm zorunlu alanlar doldurulmuş olmalı.
2. **Privacy:** App Privacy ile uygulama davranışı tutarlı olmalı.
3. **Demo Account:** İnceleme ekibi uygulamaya girebilmeli.
4. **Permission Clarification:** Mikrofon ve konum için açıklamalar net olmalı.
5. **Guideline 4.2 (Minimum Functionality):** Uygulamanın işlevi belirgin olmalı; ekran görüntüleri bunu göstermeli.

---

*Son güncelleme: Şubat 2026*
