from sqlalchemy import Column, String, Boolean, Text, BigInteger, DateTime, Table
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.sql import func
from app.core.database import Base


class Company(Base):
    __table__ = Table(
        'companies',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Required/unique fields
        Column('name', String(255), nullable=False),
        Column('slug', String(255), unique=True, nullable=False),
        # Optional fields
        Column('description', Text, nullable=True),
        Column('logo_url', String(500), nullable=True),
        Column('website', String(255), nullable=True),
        Column('email', String(255), nullable=True),
        Column('phone', String(50), nullable=True),
        Column('address', Text, nullable=True),
        # Boolean flags
        Column('is_active', Boolean, default=True, nullable=False),
        # Database connection fields
        Column('database_type', String(50), nullable=True),  # Oracle, MySQL, PostgreSQL, MSSQL, MongoDB, etc.
        Column('db_ip', String(255), nullable=True),
        Column('db_port', String(10), nullable=True),
        Column('db_username', String(255), nullable=True),
        Column('db_password', String(500), nullable=True),  # Encrypted password
        Column('db_service_name', String(255), nullable=True),  # Oracle service name
        Column('db_sid', String(255), nullable=True),  # Oracle SID (alternative to service_name)
        Column('db_schema', String(255), nullable=True),  # Database schema prefix (e.g., "OPERA" for OPERA.REP_MANAGER)
        # JSON fields
        Column('settings', JSONB, nullable=True),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<Company(id={self.id}, name='{self.name}')>"
