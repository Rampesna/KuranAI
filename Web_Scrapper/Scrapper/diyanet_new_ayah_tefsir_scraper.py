# -*- coding: utf-8 -*-
"""
Diyanet Ayet Tefsir Scraper (yeni) - kuran.diyanet.gov.tr/Tefsir
- Sure sayfalarından ayet tefsir linklerini çeker (tek ayet: 1, 2... veya aralık: 14-16)
- new_ayah_tefsir_jsons altına {surah:03d}_{ayah}.json olarak kaydeder
- ayah alanı text: "1", "2", "14-16" gibi
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

from config import DIYANET_BASE_URL, REQUEST_HEADERS
from diyanet_surah_slugs import DIYANET_SURAH_SLUGS

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
WEB_SCRAPPER_ROOT = os.path.dirname(SCRIPT_DIR)
NEW_AYAH_TEFSIR_JSON_DIR = os.path.join(WEB_SCRAPPER_ROOT, "new_ayah_tefsir_jsons")

SURAH_PAGE_PATH = "/tefsir/sure/"

logger = logging.getLogger(__name__)


def setup_logging():
    logging.basicConfig(
        level=logging.INFO,
        format="[SCRAPER] %(message)s",
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


def extract_ayet_links_with_ranges(soup, base_url):
    """
    Sure sayfasındaki ayet tefsir linklerini döndürür.
    (ayah_str, url) - sure numarası URL'de yanlış olabilir, döngüdeki surah_number kullanılacak.
    Tek: "1", "2"   Aralık: "14-16"
    URL örn: .../1-ayet-tefsiri veya .../14-16-ayet-tefsiri
    """
    links = []
    for a in soup.find_all("a", href=True):
        href = a.get("href", "")
        if "ayet-tefsiri" not in href:
            continue
        full = urljoin(base_url, href)
        # Tek: /(\d+)-ayet-tefsiri   Aralık: /(\d+-\d+)-ayet-tefsiri
        m = re.search(r"/(\d+(?:-\d+)?)-ayet-tefsiri", full)
        if m:
            ayah_str = m.group(1)
            links.append((ayah_str, full))

    def sort_key(x):
        a, _ = x
        first_num = int(a.split("-")[0]) if a else 0
        return first_num

    seen = set()
    unique = []
    for item in sorted(links, key=sort_key):
        if item[0] not in seen:
            seen.add(item[0])
            unique.append(item)
    return unique


def extract_meal_and_tefsir(soup):
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
    if "Kaynak:" in tefsir or "_Kaynak:_" in tefsir:
        tefsir = re.split(r"\n*_?Kaynak:_?", tefsir, 1, re.IGNORECASE)[0].strip()
    return meal.strip(), tefsir.strip()


def scrape_surah_ayet_links(driver, slug):
    url = urljoin(DIYANET_BASE_URL, SURAH_PAGE_PATH + slug)
    driver.get(url)
    time.sleep(2)
    soup = BeautifulSoup(driver.page_source, "html.parser")
    return extract_ayet_links_with_ranges(soup, DIYANET_BASE_URL)


def scrape_ayet_page(driver, surah_number, ayah_str, url):
    driver.get(url)
    time.sleep(1.5)
    soup = BeautifulSoup(driver.page_source, "html.parser")
    meal, tefsir = extract_meal_and_tefsir(soup)
    return {
        "surah": surah_number,
        "ayah": ayah_str,
        "meal": meal,
        "tefsir": tefsir,
    }


def ayah_to_filename_part(ayah_str):
    """ayah "1" -> "001", "14-16" -> "014-016" (dosya adı için)"""
    if "-" in ayah_str:
        parts = ayah_str.split("-", 1)
        return f"{int(parts[0]):03d}-{int(parts[1]):03d}"
    return f"{int(ayah_str):03d}"


def get_ayet_filepath(out_dir, surah_number, ayah_str):
    part = ayah_to_filename_part(ayah_str)
    return os.path.join(out_dir, f"{surah_number:03d}_{part}.json")


def save_ayet_json(data, out_dir):
    path = get_ayet_filepath(out_dir, data["surah"], data["ayah"])
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=4)
    return path


def run():
    setup_logging()
    os.makedirs(NEW_AYAH_TEFSIR_JSON_DIR, exist_ok=True)
    logger.info("Output: %s", NEW_AYAH_TEFSIR_JSON_DIR)

    driver = get_chrome_driver()
    try:
        total = 0
        for idx, slug in enumerate(DIYANET_SURAH_SLUGS, start=1):
            surah_number = idx
            logger.info("[%d/114] Sure: %s", surah_number, slug)

            try:
                ayet_links = scrape_surah_ayet_links(driver, slug)
            except Exception as e:
                logger.exception("Sure linkleri alınamadı: %s - %s", slug, e)
                continue

            for ayah_str, ayet_url in ayet_links:
                out_path = get_ayet_filepath(NEW_AYAH_TEFSIR_JSON_DIR, surah_number, ayah_str)
                if os.path.exists(out_path):
                    logger.info("  -> sure %d ayet %s zaten var, atlanıyor", surah_number, ayah_str)
                    total += 1
                    continue
                try:
                    ayet_data = scrape_ayet_page(driver, surah_number, ayah_str, ayet_url)
                    save_ayet_json(ayet_data, NEW_AYAH_TEFSIR_JSON_DIR)
                    total += 1
                    logger.info("  -> sure %d ayet %s kaydedildi", surah_number, ayah_str)
                except Exception as e:
                    logger.warning("Sure %d ayet %s atlandı: %s", surah_number, ayah_str, e)
                time.sleep(0.8)

            time.sleep(1)

        logger.info("Tamamlandı. Toplam %d ayet tefsiri.", total)
    finally:
        driver.quit()


if __name__ == "__main__":
    run()
