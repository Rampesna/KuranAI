"""
Security utilities for generating random strings for JWT secrets and passwords
"""
import secrets
import string


def generate_jwt_secret(length: int = 64) -> str:
    """
    Generate a secure random string for JWT secret (64 characters)
    Uses uppercase, lowercase, and digits
    """
    alphabet = string.ascii_letters + string.digits  # A-Z, a-z, 0-9
    return ''.join(secrets.choice(alphabet) for _ in range(length))


def generate_db_password(length: int = 16) -> str:
    """
    Generate a secure random string for database passwords (16 characters)
    Uses uppercase, lowercase, and digits
    """
    alphabet = string.ascii_letters + string.digits  # A-Z, a-z, 0-9
    return ''.join(secrets.choice(alphabet) for _ in range(length))


def generate_redis_password(length: int = 16) -> str:
    """
    Generate a secure random string for Redis password (16 characters)
    """
    return generate_db_password(length)


def generate_rabbitmq_password(length: int = 16) -> str:
    """
    Generate a secure random string for RabbitMQ password (16 characters)
    """
    return generate_db_password(length)


if __name__ == "__main__":
    print("🔐 Security String Generator")
    print("=" * 50)
    print(f"JWT Secret (64 chars): {generate_jwt_secret()}")
    print(f"DB Password (16 chars): {generate_db_password()}")
    print(f"Redis Password (16 chars): {generate_redis_password()}")
    print(f"RabbitMQ Password (16 chars): {generate_rabbitmq_password()}")
