"""
Redis client for token management and caching
"""
import redis
from app.core.config import settings
from typing import Optional
import json


class RedisClient:
    _instance: Optional['RedisClient'] = None
    _client: Optional[redis.Redis] = None

    def __init__(self):
        if RedisClient._client is None:
            RedisClient._client = redis.Redis(
                host=settings.REDIS_HOST,
                port=settings.REDIS_PORT,
                password=settings.REDIS_PASSWORD if settings.REDIS_PASSWORD else None,
                db=settings.REDIS_DB,
                decode_responses=True,  # Automatically decode bytes to strings
                socket_connect_timeout=5,
                socket_timeout=5,
                retry_on_timeout=True
            )
            # Test connection
            try:
                RedisClient._client.ping()
            except redis.ConnectionError as e:
                print(f"⚠️  Redis connection failed: {e}")
                RedisClient._client = None

    @classmethod
    def get_client(cls) -> Optional[redis.Redis]:
        if cls._client is None:
            cls()
        return cls._client

    @classmethod
    def is_available(cls) -> bool:
        client = cls.get_client()
        if client is None:
            return False
        try:
            client.ping()
            return True
        except:
            return False


# Convenience functions
def get_redis() -> Optional[redis.Redis]:
    """Get Redis client instance"""
    return RedisClient.get_client()


def redis_set(key: str, value: str, ttl: Optional[int] = None) -> bool:
    """Set a key-value pair in Redis with optional TTL"""
    client = get_redis()
    if client is None:
        return False
    try:
        if ttl:
            client.setex(key, ttl, value)
        else:
            client.set(key, value)
        return True
    except:
        return False


def redis_get(key: str) -> Optional[str]:
    """Get a value from Redis"""
    client = get_redis()
    if client is None:
        return None
    try:
        return client.get(key)
    except:
        return None


def redis_delete(key: str) -> bool:
    """Delete a key from Redis"""
    client = get_redis()
    if client is None:
        return False
    try:
        client.delete(key)
        return True
    except:
        return False


def redis_exists(key: str) -> bool:
    """Check if a key exists in Redis"""
    client = get_redis()
    if client is None:
        return False
    try:
        return client.exists(key) > 0
    except:
        return False
