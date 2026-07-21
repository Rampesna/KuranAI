# Kuran.com Veri Çekme (Scraper)

Bu klasör, kuran.com üzerinden **Tefsir** ve **Hadis** verilerini çekip JSON olarak kaydeder.

## Kurulum

```bash
cd Web_Scrapper/Scrapper
pip install -r requirements.txt
```

**Not:** Hadis scraper Selenium kullandığı için Chrome (veya Chromium) yüklü olmalıdır. `webdriver-manager` ilk çalıştırmada uygun ChromeDriver’ı indirir.

## Çalıştırma

**Her iki scraper’ı sırayla çalıştırmak (önerilen):**

```bash
cd /Users/rampesna/Projects/Kuran_AI/Web_Scrapper/Scrapper
python run_all.py
```

**Sadece Tefsir:**

```bash
python tefsir_scraper.py
```

**Sadece Hadis:**

```bash
python hadis_scraper.py
```

**Diyanet Tefsir (kuran.diyanet.gov.tr):**

```bash
python diyanet_tefsir_scraper.py
```

- **Sure tefsirleri:** `Web_Scrapper/sure_tefsirleri_jsons/` — her sure için `hakkinda`, `nuzul`, `konu`, `fazilet`.
- **Ayet tefsirleri:** `Web_Scrapper/ayet_tefsirleri_jsons/` — her ayet için ayrı JSON: `sure`, `ayet`, `meal`, `tefsir`.

Site bazen CAPTCHA gösterebilir; gerekirse script içinde `headless=False` yapıp tarayıcıda çözebilirsin.

## Çıktılar

- **Tefsir:** `Web_Scrapper/Tefsir_Jsons/`  
  - Dosya adları: `000_surelerin_inis_sebepleri_hakkinda.json`, `001_fatiha.json`, `002_bakara.json`, …  
  - Fatiha suresi `surah_number: 1` olarak kaydedilir (sitede 2. sırada gösterilse bile).

- **Hadis:** `Web_Scrapper/Hadis_Jsons/`  
  - Her bölüm bir JSON dosyası; her dosyada `hadisler` listesi (HADİS, FASIL, KONU, RAVİ, KAYNAK alanları).  
  - Sayfa içi “load more” / scroll ile tüm hadisler alınır.

Log’lar terminalde görünür; ilerlemeyi oradan takip edebilirsin.
