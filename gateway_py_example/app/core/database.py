from sqlalchemy import create_engine, text, event
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from app.core.config import settings
from typing import Generator
from datetime import datetime, timezone
import os

# Build DATABASE_URL if not provided
if not settings.DATABASE_URL:
    db_host = os.getenv("DB_HOST", "postgres")
    db_port = os.getenv("DB_PORT", "5432")
    db_user = os.getenv("DB_USER", "intelor")
    db_password = os.getenv("DB_PASSWORD", "intelor_secret_password")
    db_name = os.getenv("DB_NAME", "intelor_db")
    settings.DATABASE_URL = f"postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}"

engine = create_engine(
    settings.DATABASE_URL,
    pool_size=settings.DB_MAX_CONNECTIONS,
    max_overflow=settings.DB_MAX_CONNECTIONS - settings.DB_MAX_IDLE_CONNECTIONS,
    pool_pre_ping=True,
    pool_recycle=settings.DB_CONNECTION_MAX_LIFETIME_SECONDS,
)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Create Base for all models
Base = declarative_base()

# Event listener to automatically set created_at and updated_at
@event.listens_for(Base, 'before_insert', propagate=True)
def receive_before_insert(mapper, connection, target):
    """Set created_at and updated_at before insert"""
    if hasattr(target, 'created_at') and hasattr(target, 'updated_at'):
        now = datetime.now(timezone.utc)
        if target.created_at is None:
            target.created_at = now
        if target.updated_at is None:
            target.updated_at = now


@event.listens_for(Base, 'before_update', propagate=True)
def receive_before_update(mapper, connection, target):
    """Set updated_at before update"""
    if hasattr(target, 'updated_at'):
        target.updated_at = datetime.now(timezone.utc)


def get_db() -> Generator[Session, None, None]:
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


def health_check() -> bool:
    try:
        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))
        return True
    except:
        return False
