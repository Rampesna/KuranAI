import os
from dotenv import load_dotenv

load_dotenv()

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = os.getenv("DB_PORT", "9301")
DB_USER = os.getenv("DB_USER", "kuranai_user")
DB_PASSWORD = os.getenv("DB_PASSWORD", "")
DB_NAME = os.getenv("DB_NAME", "kuranai_pgsql")
EMBEDDING_MODEL = "text-embedding-3-large"
MAX_EMBEDDING_TOKENS = 8000

PROJECT_ROOT = os.path.dirname(os.path.abspath(__file__))
SURAH_JSON_DIR = os.path.join(PROJECT_ROOT, "..", "Web_Scrapper", "sure_tefsirleri_jsons")
AYAH_JSON_DIR = os.path.join(PROJECT_ROOT, "..", "Web_Scrapper", "ayet_tefsirleri_jsons")
NEW_AYAH_JSON_DIR = os.path.join(PROJECT_ROOT, "..", "Web_Scrapper", "new_ayah_tefsir_jsons")
