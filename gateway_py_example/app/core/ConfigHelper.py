import os
from app.core.config import settings

# Build DATABASE_URL if not provided
if not settings.DATABASE_URL:
    db_host = os.getenv("DB_HOST", "postgres")
    db_port = os.getenv("DB_PORT", "5432")
    db_user = os.getenv("DB_USER", "intelor")
    db_password = os.getenv("DB_PASSWORD", "intelor_secret_password")
    db_name = os.getenv("DB_NAME", "intelor_db")
    settings.DATABASE_URL = f"postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}"
