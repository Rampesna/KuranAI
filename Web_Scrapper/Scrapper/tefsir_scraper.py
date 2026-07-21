# -*- coding: utf-8 -*-
"""
Tefsir scraper: kuran.com/tefsir listesini tarar, her maddenin iç sayfasından
tefsir metnini alır ve JSON olarak Tefsir_Jsons altına kaydeder.
Fatiha listede 2. sırada gösterildiği için gerçek sure numarası 1 olacak şekilde kaydedilir.
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

from config import (
    BASE_URL,
    REQUEST_HEADERS,
    REQUEST_TIMEOUT,
    TEFSIR_JSON_DIR,
    TEFSIR_LIST_URL,
)

logger = logging.getLogger(__name__)


def setup_logging():
    """Configure logging for tefsir scraper."""
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )


def fetch_page(url: str) -> Optional[str]:
    """Fetch HTML content of a URL."""
    try:
        r = requests.get(url, headers=REQUEST_HEADERS, timeout=REQUEST_TIMEOUT)
        r.raise_for_status()
        r.encoding = r.apparent_encoding or "utf-8"
        return r.text
    except requests.RequestException as e:
        logger.error("İstek hatası %s: %s", url, e)
        return None


def get_tefsir_list(html: str) -> list[dict]:
    """
    Parse tefsir list page and return ordered list of items.
    Each item: { "list_index": int, "title": str, "url": str, "slug": str }
    List index 1 = Sûrelerin İniş Sebepleri (intro), 2 = Fatiha (surah 1), 3 = Bakara (surah 2), ...
    """
    soup = BeautifulSoup(html, "html.parser")
    items = []
    seen_urls = set()

    # Find all links that point to /tefsir/turkce/...
    for a in soup.find_all("a", href=True):
        href = a.get("href", "")
        full_url = urljoin(BASE_URL, href)
        if "/tefsir/turkce/" not in full_url:
            continue
        if full_url in seen_urls:
            continue
        text = (a.get_text() or "").strip()
        # Pattern: "1. Sûrelerin İniş Sebepleri Hakkında" or "2. FATİHA"
        match = re.match(r"^(\d+)\.\s*(.+)$", text)
        if not match:
            continue
        list_index = int(match.group(1))
        title = match.group(2).strip()
        parsed = urlparse(full_url)
        slug = parsed.path.rstrip("/").split("/")[-1] or "unknown"
        seen_urls.add(full_url)
        items.append({
            "list_index": list_index,
            "title": title,
            "url": full_url,
            "slug": slug,
        })

    # Sort by list_index and deduplicate by list_index (keep first occurrence)
    items.sort(key=lambda x: x["list_index"])
    by_index = {}
    for it in items:
        if it["list_index"] not in by_index:
            by_index[it["list_index"]] = it
    ordered = [by_index[i] for i in sorted(by_index.keys())]
    return ordered


def extract_tefsir_content(html: str, title: str) -> str:
    """
    Extract main tefsir text from a tefsir detail page.
    Tries common content selectors, then falls back to body text.
    """
    soup = BeautifulSoup(html, "html.parser")

    # Remove nav, footer, script, style
    for tag in soup.find_all(["script", "style", "nav", "footer", "header"]):
        tag.decompose()

    # Try to find main content container
    main = (
        soup.find("main")
        or soup.find(attrs={"role": "main"})
        or soup.find("article")
        or soup.find(class_=re.compile(r"content|tefsir|entry|post", re.I))
    )
    if main:
        content_root = main
    else:
        content_root = soup.find("body") or soup

    if not content_root:
        return ""

    # Find h1 (surah/concept title) and get text after it
    h1 = content_root.find("h1")
    if h1:
        # Collect all text from next siblings until we hit another major section or end
        parts = []
        for sibling in h1.find_next_siblings():
            if sibling.name in ("h1", "h2", "nav", "footer", "div") and sibling.get("class"):
                cls = " ".join(sibling.get("class", []))
                if "nav" in cls or "footer" in cls or "sidebar" in cls.lower():
                    break
            text = sibling.get_text(separator="\n", strip=True)
            if text:
                parts.append(text)
        return "\n\n".join(parts) if parts else content_root.get_text(separator="\n", strip=True)

    return content_root.get_text(separator="\n", strip=True)


def run_tefsir_scraper():
    """Main entry: fetch tefsir list, then each page, save JSONs."""
    setup_logging()
    os.makedirs(TEFSIR_JSON_DIR, exist_ok=True)
    logger.info("Tefsir listesi alınıyor: %s", TEFSIR_LIST_URL)

    html = fetch_page(TEFSIR_LIST_URL)
    if not html:
        logger.error("Tefsir listesi sayfası alınamadı.")
        return

    items = get_tefsir_list(html)
    logger.info("Toplam %d tefsir maddesi bulundu.", len(items))

    for it in items:
        list_index = it["list_index"]
        # Gerçek sure numarası: listede 1 = intro (0), 2 = Fatiha (1), 3 = Bakara (2), ...
        surah_number = list_index - 1  # 1->0 (intro), 2->1 (Fatiha), 3->2 (Bakara), ...
        title = it["title"]
        url = it["url"]
        slug = it["slug"]

        logger.info("[%d/%d] Sure no: %d | %s | %s", list_index, len(items), surah_number, title, url)

        page_html = fetch_page(url)
        if not page_html:
            logger.warning("Sayfa alınamadı, atlanıyor: %s", url)
            continue

        content = extract_tefsir_content(page_html, title)
        if not content.strip():
            logger.warning("İçerik boş: %s", url)

        payload = {
            "list_index": list_index,
            "surah_number": surah_number,
            "title": title,
            "slug": slug,
            "url": url,
            "content": content,
        }

        # Dosya adı: 000_intro.json, 001_fatiha.json, 002_bakara.json, ...
        safe_slug = re.sub(r"[^\w\-]", "_", slug)[:80]
        filename = f"{surah_number:03d}_{safe_slug}.json"
        out_path = os.path.join(TEFSIR_JSON_DIR, filename)

        with open(out_path, "w", encoding="utf-8") as f:
            json.dump(payload, f, ensure_ascii=False, indent=2)

        logger.info("Kaydedildi: %s", out_path)
        time.sleep(0.5)  # Polite delay

    logger.info("Tefsir scrape tamamlandı. Çıktı: %s", TEFSIR_JSON_DIR)


if __name__ == "__main__":
    run_tefsir_scraper()
