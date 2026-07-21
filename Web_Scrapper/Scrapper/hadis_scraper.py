# -*- coding: utf-8 -*-
"""
Hadis scraper: kuran.com/hadisler listesindeki başlıkları tarar, her bölümün
iç sayfasına girip hadisleri (HADİS, FASIL, KONU, RAVİ, KAYNAK) çeker.
Load more / sonsuz scroll desteklenir. Çıktı Hadis_Jsons altına yazılır.
"""

import json
import logging
import os
import re
import time
from typing import Optional
from urllib.parse import urljoin, urlparse

import requests
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager

from config import (
    BASE_URL,
    HADIS_JSON_DIR,
    HADIS_LIST_URL,
    REQUEST_HEADERS,
    REQUEST_TIMEOUT,
    SELENIUM_LOAD_MORE_MAX_SCROLLS,
    SELENIUM_SCROLL_PAUSE,
)

logger = logging.getLogger(__name__)


def setup_logging():
    """Configure logging for hadis scraper."""
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )


def get_chrome_driver():
    """Create headless Chrome WebDriver."""
    options = Options()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-gpu")
    options.add_argument("--window-size=1920,1080")
    if REQUEST_HEADERS.get("User-Agent"):
        options.add_argument(f'--user-agent={REQUEST_HEADERS["User-Agent"]}')
    service = Service(ChromeDriverManager().install())
    driver = webdriver.Chrome(service=service, options=options)
    return driver


def fetch_hadis_list_html() -> Optional[str]:
    """Fetch hadis list page with requests (no JS needed for list)."""
    try:
        r = requests.get(HADIS_LIST_URL, headers=REQUEST_HEADERS, timeout=REQUEST_TIMEOUT)
        r.raise_for_status()
        r.encoding = r.apparent_encoding or "utf-8"
        return r.text
    except Exception as e:
        logger.error("Hadis listesi alınamadı: %s", e)
        return None


def get_hadis_sections(html: str) -> list[dict]:
    """
    Parse hadis list page and return sections.
    Each: { "section_index": int, "title": str, "url": str, "slug": str }
    """
    soup = BeautifulSoup(html, "html.parser")
    items = []
    seen_urls = set()

    for a in soup.find_all("a", href=True):
        href = a.get("href", "")
        full_url = urljoin(BASE_URL, href)
        if "/hadisler/turkce/" not in full_url:
            continue
        if full_url in seen_urls:
            continue
        text = (a.get_text() or "").strip()
        match = re.match(r"^(\d+)\.\s*(.+)$", text)
        if not match:
            continue
        section_index = int(match.group(1))
        title = match.group(2).strip()
        parsed = urlparse(full_url)
        slug = parsed.path.rstrip("/").split("/")[-1] or "unknown"
        seen_urls.add(full_url)
        items.append({
            "section_index": section_index,
            "title": title,
            "url": full_url,
            "slug": slug,
        })

    items.sort(key=lambda x: x["section_index"])
    by_index = {}
    for it in items:
        if it["section_index"] not in by_index:
            by_index[it["section_index"]] = it
    return [by_index[i] for i in sorted(by_index.keys())]


def parse_hadis_entries(page_text: str) -> list[dict]:
    """
    Parse page text into list of hadis entries.
    Each entry: { "hadis": str, "fasil": str, "konu": str, "ravi": str, "kaynak": str }
    """
    entries = []
    # Split by "HADİS:" to get blocks (first block may be header, skip empty)
    blocks = re.split(r"\n\s*HADİS\s*:\s*", page_text, flags=re.IGNORECASE)
    for block in blocks:
        block = block.strip()
        if not block or len(block) < 5:
            continue

        entry = {"hadis": "", "fasil": "", "konu": "", "ravi": "", "kaynak": ""}

        # Hadis text: from start of block until first label (FASIL/KONU/RAVİ/KAYNAK)
        label_pattern = r"\n\s*(?:FASIL|KONU|RAVİ|KAYNAK)\s*:"
        first_label = re.search(label_pattern, block, re.IGNORECASE)
        if first_label:
            entry["hadis"] = block[: first_label.start()].strip()
        else:
            entry["hadis"] = block.strip()

        # Extract FASIL:, KONU:, RAVİ:, KAYNAK:
        fasil_m = re.search(r"\n\s*FASIL\s*:\s*(.+?)(?=\n\s*(?:HADİS|FASIL|KONU|RAVİ|KAYNAK)\s*:|\Z)", block, re.DOTALL | re.IGNORECASE)
        konu_m = re.search(r"\n\s*KONU\s*:\s*(.+?)(?=\n\s*(?:HADİS|FASIL|KONU|RAVİ|KAYNAK)\s*:|\Z)", block, re.DOTALL | re.IGNORECASE)
        ravi_m = re.search(r"\n\s*RAVİ\s*:\s*(.+?)(?=\n\s*(?:HADİS|FASIL|KONU|RAVİ|KAYNAK)\s*:|\Z)", block, re.DOTALL | re.IGNORECASE)
        kaynak_m = re.search(r"\n\s*KAYNAK\s*:\s*(.+?)(?=\n\s*(?:HADİS|FASIL|KONU|RAVİ|KAYNAK)\s*:|\Z)", block, re.DOTALL | re.IGNORECASE)

        if fasil_m:
            entry["fasil"] = fasil_m.group(1).strip()
        if konu_m:
            entry["konu"] = konu_m.group(1).strip()
        if ravi_m:
            entry["ravi"] = ravi_m.group(1).strip()
        if kaynak_m:
            entry["kaynak"] = kaynak_m.group(1).strip()

        # Skip header/noise: must have hadis text or at least one meta field
        if len(entry["hadis"]) > 15 or entry["fasil"] or entry["konu"] or entry["ravi"] or entry["kaynak"]:
            entries.append(entry)

    return entries


def get_full_page_content_with_load_more(driver: webdriver.Chrome, url: str) -> str:
    """
    Open URL in Selenium, scroll to trigger load more until no new content,
    then return full page text (main content area).
    """
    driver.get(url)
    time.sleep(2)

    last_height = driver.execute_script("return document.body.scrollHeight")
    scroll_count = 0

    while scroll_count < SELENIUM_LOAD_MORE_MAX_SCROLLS:
        driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
        time.sleep(SELENIUM_SCROLL_PAUSE)
        new_height = driver.execute_script("return document.body.scrollHeight")
        if new_height == last_height:
            # Try clicking "Daha fazla" / "Load more" if present
            try:
                more_btn = driver.find_elements(By.XPATH, "//*[contains(text(),'Daha fazla') or contains(text(),'Load more') or contains(text(),'Devam')]")
                if more_btn:
                    more_btn[0].click()
                    time.sleep(SELENIUM_SCROLL_PAUSE)
                    new_height = driver.execute_script("return document.body.scrollHeight")
            except Exception:
                pass
        if new_height == last_height:
            break
        last_height = new_height
        scroll_count += 1

    # Get main content text
    soup = BeautifulSoup(driver.page_source, "html.parser")
    for tag in soup.find_all(["script", "style", "nav", "footer", "header"]):
        tag.decompose()

    main = soup.find("main") or soup.find(attrs={"role": "main"}) or soup.find("article")
    if main:
        return main.get_text(separator="\n", strip=True)
    body = soup.find("body")
    return body.get_text(separator="\n", strip=True) if body else ""


def run_hadis_scraper():
    """Main entry: fetch section list, then each section with Selenium, parse hadis, save JSONs."""
    setup_logging()
    os.makedirs(HADIS_JSON_DIR, exist_ok=True)

    logger.info("Hadis bölüm listesi alınıyor: %s", HADIS_LIST_URL)
    list_html = fetch_hadis_list_html()
    if not list_html:
        return

    sections = get_hadis_sections(list_html)
    logger.info("Toplam %d hadis bölümü bulundu.", len(sections))

    driver = get_chrome_driver()
    try:
        for i, sec in enumerate(sections):
            title = sec["title"]
            url = sec["url"]
            slug = sec["slug"]
            section_index = sec["section_index"]

            logger.info("[%d/%d] %s | %s", i + 1, len(sections), title, url)

            page_text = get_full_page_content_with_load_more(driver, url)
            entries = parse_hadis_entries(page_text)

            logger.info("  -> %d hadis kaydı parse edildi.", len(entries))

            payload = {
                "section_index": section_index,
                "title": title,
                "slug": slug,
                "url": url,
                "hadis_count": len(entries),
                "hadisler": entries,
            }

            safe_slug = re.sub(r"[^\w\-]", "_", slug)[:80]
            filename = f"{section_index:03d}_{safe_slug}.json"
            out_path = os.path.join(HADIS_JSON_DIR, filename)

            with open(out_path, "w", encoding="utf-8") as f:
                json.dump(payload, f, ensure_ascii=False, indent=2)

            logger.info("  -> Kaydedildi: %s", out_path)
            time.sleep(0.5)
    finally:
        driver.quit()

    logger.info("Hadis scrape tamamlandı. Çıktı: %s", HADIS_JSON_DIR)


if __name__ == "__main__":
    run_hadis_scraper()
