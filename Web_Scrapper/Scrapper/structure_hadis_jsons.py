# -*- coding: utf-8 -*-
"""
Hadis JSON'larını temizleyip tek bir yapı altında toplar:
- Kaynak/hadis vb. alanlardaki site footer/junk metni kaldırılır
- Nav/header kaynaklı sahte hadis kayıtları elenir
- Çıktı: { title, slug, hadis_count, hadisler: [ { hadis, fasil, konu, ravi, kaynak } ] }
- Klasör: structured_hadis_jsons
"""

import json
import logging
import os
import re
import glob

from config import HADIS_JSON_DIR, STRUCTURED_HADIS_JSON_DIR

logger = logging.getLogger(__name__)

# Site footer / junk: bu ifadeden sonrası atılacak (kaynak ve diğer alanlarda)
FOOTER_JUNK_START = "Kuran.com ile Kur'an-ı Kerim'i okuyun"
FOOTER_JUNK_PATTERN = re.compile(
    r"\n?Kuran\.com ile Kur'an-ı Kerim'i okuyun.*$",
    re.IGNORECASE | re.DOTALL,
)

# Sayfa başlığı/nav'dan gelen sahte hadis (bu içerik tek başına gelirse kayıt atlanır)
NOISE_HADIS_PATTERNS = [
    "Hadis\nAra\nMeal\nTürkçe\nİngilizce\nHadis",
    "Hadis\nAra\nMeal\nTürkçe\nİngilizce",
    "Hadis\nAra\nMeal",
]


def clean_text(s):
    """Metinden site footer ve gereksiz boşlukları temizler."""
    if not s or not isinstance(s, str):
        return ""
    s = FOOTER_JUNK_PATTERN.sub("", s)
    if FOOTER_JUNK_START.lower() in s.lower():
        s = s.split(FOOTER_JUNK_START)[0].strip()
    s = s.strip()
    return s


def is_noise_entry(entry):
    """Nav/header kaynaklı sahte hadis kaydı mı?"""
    hadis = (entry.get("hadis") or "").strip()
    fasil = (entry.get("fasil") or "").strip()
    konu = (entry.get("konu") or "").strip()
    ravi = (entry.get("ravi") or "").strip()
    kaynak = (entry.get("kaynak") or "").strip()

    if not hadis and not fasil and not konu and not ravi and not kaynak:
        return True
    for noise in NOISE_HADIS_PATTERNS:
        if hadis == noise or hadis.startswith(noise):
            if not fasil and not konu and not ravi and not kaynak:
                return True
    if len(hadis) < 20 and not fasil and not konu and not ravi and not kaynak:
        return True
    return False


def clean_entry(entry):
    """Tek bir hadis objesinin tüm string alanlarını temizler."""
    return {
        "hadis": clean_text(entry.get("hadis") or ""),
        "fasil": clean_text(entry.get("fasil") or ""),
        "konu": clean_text(entry.get("konu") or ""),
        "ravi": clean_text(entry.get("ravi") or ""),
        "kaynak": clean_text(entry.get("kaynak") or ""),
    }


def process_file(path):
    """Tek bir hadis JSON dosyasını okuyup temizlenmiş yapıyı döndürür."""
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)

    title = (data.get("title") or "").strip()
    slug = (data.get("slug") or "").strip()
    raw_list = data.get("hadisler") or []

    cleaned = []
    for entry in raw_list:
        if is_noise_entry(entry):
            continue
        c = clean_entry(entry)
        if not c["hadis"] and not c["fasil"] and not c["konu"] and not c["ravi"] and not c["kaynak"]:
            continue
        cleaned.append(c)

    return {
        "title": title,
        "slug": slug,
        "hadis_count": len(cleaned),
        "hadisler": cleaned,
    }


def run():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    os.makedirs(STRUCTURED_HADIS_JSON_DIR, exist_ok=True)

    pattern = os.path.join(HADIS_JSON_DIR, "*.json")
    files = sorted(glob.glob(pattern))
    if not files:
        logger.warning("Hadis_Jsons altında JSON dosyası bulunamadı: %s", pattern)
        return

    logger.info("Toplam %d hadis JSON dosyası işlenecek.", len(files))

    for path in files:
        name = os.path.basename(path)
        try:
            out = process_file(path)
            out_path = os.path.join(STRUCTURED_HADIS_JSON_DIR, name)
            with open(out_path, "w", encoding="utf-8") as f:
                json.dump(out, f, ensure_ascii=False, indent=4)
            logger.info("Yazıldı: %s (hadis_count=%d)", name, out["hadis_count"])
        except Exception as e:
            logger.exception("İşlenirken hata: %s - %s", path, e)

    logger.info("Tamamlandı. Çıktı: %s", STRUCTURED_HADIS_JSON_DIR)


if __name__ == "__main__":
    run()
