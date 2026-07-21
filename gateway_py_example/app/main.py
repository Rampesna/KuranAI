from fastapi import FastAPI, Response, Request, HTTPException, status
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError
from contextlib import asynccontextmanager
from app.core.config import settings
from app.core.database import health_check
from app.core.migrations import run_migrations, seed_database
from app.routers import AdminRoutes, ManagerRoutes, EmployeeRoutes
from app.http.middlewares.Cors import setup_cors
from scalar_fastapi import get_scalar_api_reference
import logging
import json

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Run migrations and seed data on startup"""
    try:
        import os
        # Check if FULL_MIGRATION environment variable is set
        if os.getenv("FULL_MIGRATION", "false").lower() == "true":
            logger.info("Running FULL database migration (drop all, recreate, seed, import SQL)...")
            from app.core.migrations import run_full_migration
            run_full_migration()
        else:
            logger.info("Running database migrations...")
            run_migrations()
            logger.info("Seeding database...")
            seed_database()
        logger.info("Startup completed successfully")
    except Exception as e:
        logger.error(f"Error during startup: {e}")
        # Don't fail startup, but log the error
    yield
    # Shutdown logic (if needed)


# Determine server URL from environment or settings, default to localhost
import os
api_base_url = os.getenv("API_BASE_URL") or getattr(settings, "API_BASE_URL", None) or "http://localhost:4551"
servers = [{"url": api_base_url, "description": "API Server"}]

app = FastAPI(
    title="Intelor Gateway API",
    description="Intelor Gateway API Documentation",
    version="1.0.0",
    docs_url=None,
    redoc_url=None,
    lifespan=lifespan,
    servers=servers
)

# Setup CORS
setup_cors(app)


# Global exception handlers for consistent response format
@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException):
    """Convert HTTPException to ServiceResponse format"""
    response_data = {
        "isSuccess": False,
        "message": exc.detail if isinstance(exc.detail, str) else str(exc.detail),
        "statusCode": exc.status_code,
        "data": None
    }
    return JSONResponse(
        content=response_data,
        status_code=exc.status_code
    )


@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    """Convert validation errors to ServiceResponse format"""
    errors = exc.errors()
    error_messages = []
    for error in errors:
        field = " -> ".join(str(loc) for loc in error.get("loc", []))
        msg = error.get("msg", "Validation error")
        error_messages.append(f"{field}: {msg}")
    
    response_data = {
        "isSuccess": False,
        "message": "; ".join(error_messages) if error_messages else "Validation error",
        "statusCode": 422,
        "data": None
    }
    return JSONResponse(
        content=response_data,
        status_code=422
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    """Convert all other exceptions to ServiceResponse format"""
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    response_data = {
        "isSuccess": False,
        "message": "Internal server error",
        "statusCode": 500,
        "data": None
    }
    return JSONResponse(
        content=response_data,
        status_code=500
    )


# Include routers
app.include_router(AdminRoutes.router)
app.include_router(ManagerRoutes.router)
app.include_router(EmployeeRoutes.router)

# Scalar documentation
@app.get("/docs", include_in_schema=False)
async def scalar_html():
    return get_scalar_api_reference(
        openapi_url=app.openapi_url,
        title=app.title + " - Scalar",
    )


@app.get("/openapi.json", include_in_schema=False)
async def get_openapi():
    """Return OpenAPI schema (response bodies are empty - users should test endpoints to see actual responses)"""
    return app.openapi()


@app.get("/health")
async def health():
    db_healthy = health_check()
    return JSONResponse(
        content={
            "status": "healthy" if db_healthy else "unhealthy",
            "data": {
                "database": "healthy" if db_healthy else "unhealthy"
            }
        },
        status_code=200 if db_healthy else 503
    )


@app.get("/ready")
async def ready():
    db_healthy = health_check()
    return JSONResponse(
        content={
            "status": "ready" if db_healthy else "not ready",
            "data": {
                "database": "ready" if db_healthy else "not ready"
            }
        },
        status_code=200 if db_healthy else 503
    )
