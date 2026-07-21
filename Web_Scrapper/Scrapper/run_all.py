# -*- coding: utf-8 -*-
"""
Tefsir ve Hadis scraper'ları sırayla çalıştırır.
Terminalden izleyebilmeniz için log'lar INFO seviyesinde stdout'a yazılır.
"""

import logging
import sys

# Logging: hem dosyaya hem konsola
def setup_root_logging():
    root = logging.getLogger()
    root.setLevel(logging.INFO)
    fmt = logging.Formatter("%(asctime)s [%(levelname)s] %(name)s: %(message)s", datefmt="%Y-%m-%d %H:%M:%S")
    if not root.handlers:
        h = logging.StreamHandler(sys.stdout)
        h.setFormatter(fmt)
        root.addHandler(h)


def main():
    setup_root_logging()
    logger = logging.getLogger("run_all")

    logger.info("========== 1/2 TEFSİR SCRAPER BAŞLIYOR ==========")
    from tefsir_scraper import run_tefsir_scraper
    run_tefsir_scraper()

    logger.info("========== 2/2 HADİS SCRAPER BAŞLIYOR ==========")
    from hadis_scraper import run_hadis_scraper
    run_hadis_scraper()

    logger.info("========== TÜM VERİ ÇEKME İŞLEMLERİ TAMAMLANDI ==========")


if __name__ == "__main__":
    main()
