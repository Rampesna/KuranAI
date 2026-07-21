#!/usr/bin/env python3
"""
App Store Connect API - JWT Token üretici
Çalıştırınca JWT token'ı konsola yazar.

Kurulum: pip install PyJWT cryptography

Kullanım:
  1. Aşağıdaki ISSUER_ID, KEY_ID ve PRIVATE_KEY_PATH değerlerini doldurun
  2. python appstore_api_key.py
"""
import os
import time

try:
    import jwt
except ImportError:
    print("Hata: PyJWT kurulu değil. 'pip install PyJWT cryptography' çalıştırın.")
    exit(1)

# ============= BURAYI DOLDURUN =============
# App Store Connect → Users and Access → Integrations → App Store Connect API
ISSUER_ID = os.environ.get("APPSTORE_ISSUER_ID", "buraya-issuer-id-yazin")
KEY_ID = os.environ.get("APPSTORE_KEY_ID", "TTZYVNJ39C")

# .p8 dosyası Kuran_AI_iOS/ altında
PRIVATE_KEY_PATH = os.environ.get(
    "APPSTORE_KEY_PATH",
    os.path.join(os.path.dirname(__file__), "Kuran_AI_iOS", "AuthKey_TTZYVNJ39C.p8"),
)
# ==========================================


def generate_token():
    if not os.path.exists(PRIVATE_KEY_PATH):
        print(f"Hata: .p8 dosyası bulunamadı: {PRIVATE_KEY_PATH}")
        print("  - App Store Connect'ten API Key oluşturup .p8 dosyasını indirin")
        print("  - Dosyayı bu script ile aynı dizine koyun veya PRIVATE_KEY_PATH'i ayarlayın")
        exit(1)

    with open(PRIVATE_KEY_PATH, "r") as f:
        private_key = f.read()

    payload = {
        "iss": ISSUER_ID,
        "iat": int(time.time()),
        "exp": int(time.time()) + 1200,  # 20 dakika geçerli
        "aud": "appstoreconnect-v1",
    }

    token = jwt.encode(
        payload,
        private_key,
        algorithm="ES256",
        headers={"kid": KEY_ID},
    )

    # PyJWT 2.x string döner, 1.x bytes döner
    if isinstance(token, bytes):
        token = token.decode("utf-8")

    print("=== App Store Connect API JWT Token ===")
    print(token)
    print("=======================================")
    return token


if __name__ == "__main__":
    generate_token()
