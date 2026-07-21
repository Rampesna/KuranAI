from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import field_validator, model_validator
from typing import List, Union, Optional, Any
import json
import os

# Fix REDIS_PASSWORD if it's empty string (causes None in Pydantic)
if os.getenv("REDIS_PASSWORD") == "":
    os.environ["REDIS_PASSWORD"] = " "  # Set to space, will be stripped


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        case_sensitive=True,
        env_parse_none_str="",  # Parse None as empty string
        extra="ignore",  # Ignore extra environment variables like CORS_ALLOWED_ORIGINS
        env_ignore_empty=False,  # Don't ignore empty strings
    )

    # App
    APP_ENV: str = "production"
    APP_DEBUG: bool = False
    GATEWAY_PORT: int = 4551
    GATEWAY_HOST: str = "0.0.0.0"
    API_BASE_URL: str = "http://localhost:4551"  # Base URL for API documentation (Scalar/OpenAPI)

    # JWT
    JWT_SECRET: str
    JWT_ACCESS_TOKEN_EXPIRY_HOURS: int = 24
    JWT_REFRESH_TOKEN_EXPIRY_HOURS: int = 168  # 7 days

    # Database
    DATABASE_URL: str = ""
    DB_MAX_CONNECTIONS: int = 20
    DB_MAX_IDLE_CONNECTIONS: int = 5
    DB_CONNECTION_MAX_LIFETIME_SECONDS: int = 3600

    # Redis
    REDIS_HOST: str = "redis"
    REDIS_PORT: int = 6379
    REDIS_PASSWORD: str = ""  # Default empty
    REDIS_DB: int = 0
    REDIS_CACHE_TTL_SECONDS: int = 3600
    
    @field_validator("REDIS_PASSWORD", mode="before")
    @classmethod
    def normalize_redis_password(cls, v):
        if v is None or v == " ":
            return ""
        return str(v).strip()

    # RabbitMQ
    RABBITMQ_HOST: str = "rabbitmq"
    RABBITMQ_PORT: int = 5672
    RABBITMQ_USER: str = "guest"
    RABBITMQ_PASSWORD: str = "guest"
    RABBITMQ_VHOST: str = "/"

    # OpenAI
    OPENAI_API_KEY: str
    OPENAI_MODEL: str = "gpt-5.2"
    OPENAI_EMBEDDING_MODEL: str = "text-embedding-3-large"
    RAG_TOP_K: int = 3
    MAX_SQL_TRIES: int = 3
    SQL_STATEMENT_TIMEOUT_SECONDS: int = 30

    # CORS - Store as string, will be parsed after initialization
    CORS_ALLOWED_ORIGINS_STR: str = "*"
    CORS_ALLOWED_METHODS_STR: str = "*"
    CORS_ALLOWED_HEADERS_STR: str = "*"


def _parse_cors(value: str) -> List[str]:
    if not value:
        return ["*"]
    if value.startswith("["):
        try:
            return json.loads(value)
        except:
            pass
    return [item.strip() for item in value.split(",") if item.strip()] or ["*"]


settings = Settings()
