from fastapi.middleware.cors import CORSMiddleware
from fastapi import FastAPI
from app.core.config import settings
import json


def _parse_cors_list(value: str) -> list:
    if not value:
        return ["*"]
    if value.startswith("["):
        try:
            return json.loads(value)
        except:
            pass
    return [item.strip() for item in value.split(",") if item.strip()] or ["*"]


def setup_cors(app: FastAPI):
    origins = _parse_cors_list(settings.CORS_ALLOWED_ORIGINS_STR)
    methods = _parse_cors_list(settings.CORS_ALLOWED_METHODS_STR)
    headers = _parse_cors_list(settings.CORS_ALLOWED_HEADERS_STR)
    
    app.add_middleware(
        CORSMiddleware,
        allow_origins=origins,
        allow_credentials=True,
        allow_methods=methods,
        allow_headers=headers,
    )
