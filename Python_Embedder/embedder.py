"""
Embedder: OpenAI embedding API ile surah/ayah exegesis JSON'larını vektörleştirip DB'ye yazar.
"""
import json
import logging
import re
import sys
from pathlib import Path

import psycopg2
import tiktoken
from dotenv import load_dotenv
from openai import OpenAI

from config import (
    NEW_AYAH_JSON_DIR,
    EMBEDDING_MODEL,
    MAX_EMBEDDING_TOKENS,
    OPENAI_API_KEY,
    SURAH_JSON_DIR,
    DB_HOST,
    DB_PORT,
    DB_USER,
    DB_PASSWORD,
    DB_NAME,
)

load_dotenv()

LOG = logging.getLogger("embedder")
LOG.setLevel(logging.INFO)
handler = logging.StreamHandler(sys.stdout)
handler.setFormatter(logging.Formatter("%(message)s"))
LOG.addHandler(handler)


def get_conn():
    return psycopg2.connect(
        host=DB_HOST,
        port=int(DB_PORT),
        user=DB_USER,
        password=DB_PASSWORD,
        dbname=DB_NAME,
    )


def truncate_to_max_tokens(text: str, max_tokens: int = MAX_EMBEDDING_TOKENS) -> str:
    """Metni max_tokens (8192 limiti için 8000) ile kırpar."""
    try:
        enc = tiktoken.encoding_for_model("text-embedding-3-large")
    except KeyError:
        enc = tiktoken.get_encoding("cl100k_base")
    tokens = enc.encode(text)
    if len(tokens) <= max_tokens:
        return text
    LOG.info("[LOG] truncating text %d -> %d tokens", len(tokens), max_tokens)
    return enc.decode(tokens[:max_tokens])


def get_embedding(client: OpenAI, text: str) -> tuple[list[float], int]:
    """OpenAI embedding API çağırır. (embedding, token_count) döner. Uzun metinler kırpılır."""
    text = truncate_to_max_tokens(text)
    resp = client.embeddings.create(model=EMBEDDING_MODEL, input=text)
    vec = resp.data[0].embedding
    tok = resp.usage.total_tokens if resp.usage else 0
    return vec, tok


def vec_to_pg(vec: list[float]) -> str:
    """Python listesini PostgreSQL vector string formatına çevirir."""
    return "[" + ",".join(str(x) for x in vec) + "]"


def build_surah_text(data: dict) -> str:
    """Surah exegesis JSON içeriğinden embedding için metin oluşturur."""
    parts = []
    for key in ["title", "hakkinda", "nuzul", "konu", "fazilet"]:
        if key in data and data[key]:
            parts.append(str(data[key]))
    return "\n\n".join(parts)


def build_ayah_text(data: dict) -> str:
    """Ayah exegesis JSON içeriğinden embedding için metin oluşturur."""
    parts = []
    if data.get("meal"):
        parts.append(data["meal"])
    if data.get("tefsir"):
        parts.append(data["tefsir"])
    return "\n\n".join(parts)


def sorted_surah_files(dir_path: str) -> list[Path]:
    """sure_tefsirleri_jsons içindeki JSON dosyalarını sura numarasına göre sıralar."""
    path = Path(dir_path)
    if not path.exists():
        return []
    files = list(path.glob("*.json"))
    pattern = re.compile(r"^(\d+)_")

    def sort_key(p: Path) -> int:
        m = pattern.match(p.name)
        return int(m.group(1)) if m else 9999

    return sorted(files, key=sort_key)


def sorted_new_ayah_files(dir_path: str) -> list[Path]:
    """new_ayah_tefsir_jsons: {surah:03d}_{ayah}.json (ayah: 001 veya 014-016)"""
    path = Path(dir_path)
    if not path.exists():
        return []
    files = list(path.glob("*.json"))
    pattern = re.compile(r"^(\d+)_([\d\-]+)\.json$")

    def sort_key(p: Path):
        m = pattern.match(p.name)
        if not m:
            return (999, 999)
        surah = int(m.group(1))
        ayah_part = m.group(2)
        first_ayah = int(ayah_part.split("-")[0]) if ayah_part else 0
        return (surah, first_ayah)

    return sorted(files, key=sort_key)


def process_surah_exegesis(client: OpenAI, conn) -> int:
    """sure_tefsirleri_jsons → embedding → surah_exegesis_embeddings."""
    files = sorted_surah_files(SURAH_JSON_DIR)
    if not files:
        LOG.warning("[SURAH] no JSON files found in %s", SURAH_JSON_DIR)
        return 0

    LOG.info("[SURAH] starting %d surah files", len(files))
    cur = conn.cursor()
    inserted = 0

    for fp in files:
        try:
            with open(fp, "r", encoding="utf-8") as f:
                data = json.load(f)

            sura_num = data.get("sura_number")
            if sura_num is None:
                LOG.warning("[SURAH] skip %s: no sura_number", fp.name)
                continue

            text = build_surah_text(data)
            if not text.strip():
                LOG.warning("[SURAH] skip %s: empty text", fp.name)
                continue

            LOG.info("[SURAH] embedding sura_number=%s (%s)", sura_num, fp.name)
            vec, token_count = get_embedding(client, text)

            cur.execute(
                """
                INSERT INTO surah_exegesis_embeddings (surah_number, data, embedding, token_count)
                VALUES (%s, %s, %s::vector, %s)
                ON CONFLICT (surah_number) DO NOTHING
                """,
                (
                    sura_num,
                    json.dumps(data, ensure_ascii=False),
                    vec_to_pg(vec),
                    token_count,
                ),
            )
            conn.commit()
            inserted += 1
            LOG.info("[SURAH] saved sura_number=%s tokens=%s", sura_num, token_count)

        except Exception as e:
            conn.rollback()
            LOG.error("[SURAH] error %s: %s", fp.name, e)
            raise

    cur.close()
    LOG.info("[SURAH] done. inserted=%d", inserted)
    return inserted


def ensure_ayah_table_schema(conn):
    """ayah_exegesis_embeddings: ayah sütununu text, surah sütununu bigint yapar (gerekirse)."""
    cur = conn.cursor()
    try:
        cur.execute("""
            SELECT data_type FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = 'ayah_exegesis_embeddings' AND column_name = 'ayah'
        """)
        row = cur.fetchone()
        if row and row[0] != "text":
            LOG.info("[AYAH] altering ayah column to text")
            cur.execute("""
                ALTER TABLE ayah_exegesis_embeddings
                ALTER COLUMN ayah TYPE text USING ayah::text
            """)
            conn.commit()
        cur.execute("""
            SELECT data_type FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = 'ayah_exegesis_embeddings' AND column_name = 'surah'
        """)
        row = cur.fetchone()
        if row and row[0] not in ("bigint", "int8"):
            LOG.info("[AYAH] altering surah column to bigint")
            cur.execute("""
                ALTER TABLE ayah_exegesis_embeddings
                ALTER COLUMN surah TYPE bigint USING surah::bigint
            """)
            conn.commit()
    except Exception as e:
        conn.rollback()
        raise
    finally:
        cur.close()


def process_ayah_exegesis(client: OpenAI, conn) -> int:
    """new_ayah_tefsir_jsons → embedding → ayah_exegesis_embeddings.
    Tabloyu truncate eder, JSON'dan surah (int) ve ayah (str: \"1\" veya \"14-16\") alır."""
    ensure_ayah_table_schema(conn)

    files = sorted_new_ayah_files(NEW_AYAH_JSON_DIR)
    if not files:
        LOG.warning("[AYAH] no JSON files found in %s", NEW_AYAH_JSON_DIR)
        return 0

    LOG.info("[AYAH] truncating ayah_exegesis_embeddings")
    cur = conn.cursor()
    cur.execute("TRUNCATE TABLE ayah_exegesis_embeddings RESTART IDENTITY")
    conn.commit()

    LOG.info("[AYAH] starting %d ayah files", len(files))
    inserted = 0

    for fp in files:
        try:
            with open(fp, "r", encoding="utf-8") as f:
                data = json.load(f)

            surah = data.get("surah")
            ayah = data.get("ayah")
            if surah is None or ayah is None:
                LOG.warning("[AYAH] skip %s: no surah/ayah in JSON", fp.name)
                continue

            ayah_str = str(ayah) if not isinstance(ayah, str) else ayah

            text = build_ayah_text(data)
            if not text.strip():
                LOG.warning("[AYAH] skip %s: empty text", fp.name)
                continue

            LOG.info("[AYAH] embedding surah=%s ayah=%s (%s)", surah, ayah_str, fp.name)
            vec, token_count = get_embedding(client, text)

            cur.execute(
                """
                INSERT INTO ayah_exegesis_embeddings (surah, ayah, data, embedding, token_count)
                VALUES (%s, %s, %s, %s::vector, %s)
                """,
                (
                    surah,
                    ayah_str,
                    json.dumps(data, ensure_ascii=False),
                    vec_to_pg(vec),
                    token_count,
                ),
            )
            conn.commit()
            inserted += 1
            LOG.info("[AYAH] saved surah=%s ayah=%s tokens=%s", surah, ayah_str, token_count)

        except Exception as e:
            conn.rollback()
            LOG.error("[AYAH] error %s: %s", fp.name, e)
            raise

    cur.close()
    LOG.info("[AYAH] done. inserted=%d", inserted)
    return inserted


def main():
    if not OPENAI_API_KEY:
        LOG.error("[LOG] OPENAI_API_KEY not set in .env")
        sys.exit(1)

    LOG.info("[LOG] starting embedder")
    client = OpenAI(api_key=OPENAI_API_KEY)

    try:
        conn = get_conn()
        LOG.info("[LOG] DB connected %s:%s/%s", DB_HOST, DB_PORT, DB_NAME)
    except Exception as e:
        LOG.error("[LOG] DB connect failed: %s", e)
        sys.exit(1)

    try:
        LOG.info("[LOG] processing ayah_exegesis only (surah already done)")
        process_ayah_exegesis(client, conn)
    finally:
        conn.close()
        LOG.info("[LOG] finished")


if __name__ == "__main__":
    main()
