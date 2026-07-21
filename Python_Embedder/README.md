# Python Embedder

OpenAI embedding API ile surah ve ayah exegesis JSON'larını vektörleştirip PostgreSQL'e yazar.

## Kurulum

```bash
cd Python_Embedder
pip install -r requirements.txt
```

## Yapılandırma

`.env` dosyasını düzenleyin (veya `.env.example`'dan kopyalayın):

- `OPENAI_API_KEY` – OpenAI API anahtarı
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` – PostgreSQL bağlantı bilgileri

## Çalıştırma

```bash
python embedder.py
```

Önce `sure_tefsirleri_jsons` → `surah_exegesis_embeddings`, ardından `ayet_tefsirleri_jsons` → `ayah_exegesis_embeddings` işlenir.
