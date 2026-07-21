# -*- coding: utf-8 -*-
"""Scraper configuration."""

import os

BASE_URL = "https://www.kuran.com"
TEFSIR_LIST_URL = f"{BASE_URL}/tefsir"
HADIS_LIST_URL = f"{BASE_URL}/hadisler"

# Output directories (relative to Scrapper folder)
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
WEB_SCRAPPER_ROOT = os.path.dirname(SCRIPT_DIR)
TEFSIR_JSON_DIR = os.path.join(WEB_SCRAPPER_ROOT, "Tefsir_Jsons")
HADIS_JSON_DIR = os.path.join(WEB_SCRAPPER_ROOT, "Hadis_Jsons")
STRUCTURED_HADIS_JSON_DIR = os.path.join(WEB_SCRAPPER_ROOT, "structured_hadis_jsons")
SURE_TEFSIRLERI_JSON_DIR = os.path.join(WEB_SCRAPPER_ROOT, "sure_tefsirleri_jsons")
AYET_TEFSIRLERI_JSON_DIR = os.path.join(WEB_SCRAPPER_ROOT, "ayet_tefsirleri_jsons")

# Diyanet Tefsir (kuran.diyanet.gov.tr)
DIYANET_BASE_URL = "https://kuran.diyanet.gov.tr"
DIYANET_TEFSIR_INDEX = f"{DIYANET_BASE_URL}/Tefsir/"

# Request settings
REQUEST_TIMEOUT = 30
REQUEST_HEADERS = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
    "Accept-Language": "tr-TR,tr;q=0.9,en;q=0.8",
}

# Selenium (hadis load more)
SELENIUM_SCROLL_PAUSE = 1.5
SELENIUM_LOAD_MORE_MAX_SCROLLS = 100
