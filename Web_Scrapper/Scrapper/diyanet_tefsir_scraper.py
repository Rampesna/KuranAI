# -*- coding: utf-8 -*-
"""
Diyanet Tefsir scraper: kuran.diyanet.gov.tr/Tefsir/ kaynağından
- Sure sayfalarından: Hakkında, Nüzul, Konu, Fazilet -> sure_tefsirleri_jsons
- Her ayet sayfasından: meal, tefsir -> ayet_tefsirleri_jsons (her ayet ayrı JSON)
Selenium kullanır. Not: Site bazen CAPTCHA gösterebilir; headless=False ve
başta bekleme süresi ile tarayıcıda çözülebilir.
"""

import json
import logging
import os
import re
import time
from urllib.parse import urljoin

from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager

from config import (
    DIYANET_BASE_URL,
    AYET_TEFSIRLERI_JSON_DIR,
    REQUEST_HEADERS,
    SURE_TEFSIRLERI_JSON_DIR,
)
from diyanet_surah_slugs import DIYANET_SURAH_SLUGS

logger = logging.getLogger(__name__)

# Sure sayfası URL: /tefsir/sure/1-fatiha-suresi
SURAH_PAGE_PATH = "/tefsir/sure/"


def setup_logging():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )


def get_chrome_driver():
    options = Options()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-gpu")
    options.add_argument("--window-size=1920,1080")
    if REQUEST_HEADERS.get("User-Agent"):
        options.add_argument('--user-agent=' + REQUEST_HEADERS["User-Agent"])
    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=options)


def get_section_text(soup, heading_names):
    """İlk eşleşen başlıktan sonraki metni, sonraki aynı seviye başlığa kadar alır."""
    for tag in soup.find_all(["h2", "h3", "h4", "h5"]):
        t = (tag.get_text() or "").strip()
        if not t:
            continue
        for name in heading_names:
            if name.lower() in t.lower():
                parts = []
                for sib in tag.find_next_siblings():
                    if sib.name and sib.name.lower() in ("h2", "h3", "h4", "h5"):
                        break
                    parts.append(sib.get_text(separator="\n", strip=True))
                return "\n\n".join(p for p in parts if p).strip()
    return ""


def extract_surah_sections(soup):
    """Sure sayfasından Hakkında, Nüzul, Konu, Fazilet metinlerini çıkarır."""
    hakkinda = get_section_text(soup, ["Hakkında"])
    nuzul = get_section_text(soup, ["Nüzul", "Nüzûl"])
    konu = get_section_text(soup, ["Konusu", "Konu"])
    fazilet = get_section_text(soup, ["Fazileti", "Fazilet"])
    return {"hakkinda": hakkinda, "nuzul": nuzul, "konu": konu, "fazilet": fazilet}


def extract_surah_title_and_meta(soup):
    """h1 ve sayfa meta bilgilerini çıkarır (sure adı, ayet sayısı vb.)."""
    title = ""
    h1 = soup.find("h1")
    if h1:
        title = (h1.get_text() or "").strip()
    return title


def extract_ayet_links(soup, base_url):
    """Sure sayfasındaki ayet tefsir linklerini döndürür: [(sure_no, ayet_no, url), ...]"""
    links = []
    for a in soup.find_all("a", href=True):
        href = a.get("href", "")
        if "ayet-tefsiri" not in href:
            continue
        full = urljoin(base_url, href)
        # URL örn: /tefsir/Fâtiha-suresi/1/1-ayet-tefsiri veya .../1/2-ayet-tefsiri
        m = re.search(r"/(\d+)/(\d+)-ayet-tefsiri", full)
        if m:
            sura_no = int(m.group(1))
            ayet_no = int(m.group(2))
            links.append((sura_no, ayet_no, full))
    # Aynı (sure, ayet) tekrarlarını kaldır, sırala
    seen = set()
    unique = []
    for s, a, u in sorted(links, key=lambda x: (x[0], x[1])):
        if (s, a) not in seen:
            seen.add((s, a))
            unique.append((s, a, u))
    return unique


def extract_meal_and_tefsir(soup):
    """Ayet tefsir sayfasından meal ve tefsir metnini çıkarır."""
    meal = ""
    tefsir = ""
    for tag in soup.find_all(["h2", "h3"]):
        t = (tag.get_text() or "").strip()
        if "Meal" in t or "meal" in t:
            for sib in tag.find_next_siblings():
                if sib.name and sib.name.lower() in ("h2", "h3"):
                    break
                meal = sib.get_text(separator="\n", strip=True)
                if meal:
                    break
        if "Tefsir" in t or "tefsir" in t:
            for sib in tag.find_next_siblings():
                if sib.name and sib.name.lower() in ("h2", "h3"):
                    break
                tefsir = sib.get_text(separator="\n", strip=True)
                if tefsir:
                    break
    # Kaynak satırını tefsirden kes (varsa)
    if "Kaynak:" in tefsir or "_Kaynak:_" in tefsir:
        tefsir = re.split(r"\n*_?Kaynak:_?", tefsir, 1, re.IGNORECASE)[0].strip()
    return meal.strip(), tefsir.strip()


def scrape_surah_page(driver, surah_number, slug):
    """Tek bir sure sayfasını açar, bölümleri ve ayet linklerini döndürür."""
    url = urljoin(DIYANET_BASE_URL, SURAH_PAGE_PATH + slug)
    driver.get(url)
    time.sleep(2)
    soup = BeautifulSoup(driver.page_source, "html.parser")
    title = extract_surah_title_and_meta(soup)
    sections = extract_surah_sections(soup)
    ayet_links = extract_ayet_links(soup, DIYANET_BASE_URL)
    return {
        "title": title,
        "surah_number": surah_number,
        "slug": slug,
        "hakkinda": sections["hakkinda"],
        "nuzul": sections["nuzul"],
        "konu": sections["konu"],
        "fazilet": sections["fazilet"],
        "ayet_links": ayet_links,
    }


def save_surah_json(data, out_dir):
    """Sure verisini JSON olarak kaydeder (ayet_links hariç, sadece metin alanları)."""
    payload = {
        "sura_number": data["surah_number"],
        "title": data["title"],
        "slug": data["slug"],
        "hakkinda": data["hakkinda"],
        "nuzul": data["nuzul"],
        "konu": data["konu"],
        "fazilet": data["fazilet"],
    }
    name_part = data["slug"].split("-", 1)[1].replace("-suresi", "") if "-" in data["slug"] else data["slug"]
    filename = f"{data['surah_number']:03d}_{name_part}.json"
    path = os.path.join(out_dir, filename)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(payload, f, ensure_ascii=False, indent=4)
    return path


def scrape_ayet_page(driver, sura_no, ayet_no, url):
    """Tek bir ayet tefsir sayfasından meal ve tefsir döndürür."""
    driver.get(url)
    time.sleep(1.5)
    soup = BeautifulSoup(driver.page_source, "html.parser")
    meal, tefsir = extract_meal_and_tefsir(soup)
    return {
        "sure": sura_no,
        "ayet": ayet_no,
        "meal": meal,
        "tefsir": tefsir,
    }


def save_ayet_json(data, out_dir):
    """Ayet tefsirini JSON olarak kaydeder."""
    filename = f"{data['sure']:03d}_{data['ayet']:03d}.json"
    path = os.path.join(out_dir, filename)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=4)
    return path


def run():
    setup_logging()
    os.makedirs(SURE_TEFSIRLERI_JSON_DIR, exist_ok=True)
    os.makedirs(AYET_TEFSIRLERI_JSON_DIR, exist_ok=True)

    driver = get_chrome_driver()
    try:
        for idx, slug in enumerate(DIYANET_SURAH_SLUGS, start=1):
            surah_number = idx
            logger.info("[%d/114] Sure: %s", surah_number, slug)

            try:
                surah_data = scrape_surah_page(driver, surah_number, slug)
            except Exception as e:
                logger.exception("Sure sayfası alınamadı: %s - %s", slug, e)
                continue

            save_surah_json(surah_data, SURE_TEFSIRLERI_JSON_DIR)
            logger.info("  -> Sure JSON yazıldı.")

            ayet_links = surah_data.get("ayet_links") or []
            for sura_no, ayet_no, ayet_url in ayet_links:
                try:
                    ayet_data = scrape_ayet_page(driver, sura_no, ayet_no, ayet_url)
                    save_ayet_json(ayet_data, AYET_TEFSIRLERI_JSON_DIR)
                    logger.info("  -> Ayet %d/%d kaydedildi.", sura_no, ayet_no)
                except Exception as e:
                    logger.warning("Ayet %d/%d atlandı: %s", sura_no, ayet_no, e)
                time.sleep(0.8)

            time.sleep(1)
    finally:
        driver.quit()

    logger.info("Diyanet tefsir scrape tamamlandı.")
    logger.info("Sure tefsirleri: %s", SURE_TEFSIRLERI_JSON_DIR)
    logger.info("Ayet tefsirleri: %s", AYET_TEFSIRLERI_JSON_DIR)


if __name__ == "__main__":
    run()
