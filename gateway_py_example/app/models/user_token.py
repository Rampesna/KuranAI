from sqlalchemy import Column, String, Text, BigInteger, Boolean, DateTime, Index, Table
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.sql import func
from app.core.database import Base


class UserToken(Base):
    __table__ = Table(
        'personal_access_tokens',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Required fields (user identification)
        Column('user_id', BigInteger, nullable=False, index=True),
        Column('user_type', String(50), nullable=False, index=True),  # admin, manager, employee
        # Required fields (tokens)
        Column('access_token', Text, nullable=False, index=True),
        Column('access_token_expires_at', DateTime(timezone=True), nullable=False, index=True),
        # Optional fields (refresh token)
        Column('refresh_token', Text, nullable=True, index=True),
        Column('refresh_token_expires_at', DateTime(timezone=True), nullable=True, index=True),
        # Boolean flags
        Column('is_active', Boolean, default=True, nullable=False, index=True),
        # JSON fields (device information - Netflix/Google style device management)
        Column('device_info', JSONB, nullable=True),
        # Timestamps (before final timestamps)
        Column('last_used_at', DateTime(timezone=True), nullable=True, index=True),
        Column('logged_out_at', DateTime(timezone=True), nullable=True),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Indexes for faster queries
    __table_args__ = (
        Index("idx_personal_access_tokens_user_active", "user_id", "user_type", "is_active"),
        Index("idx_personal_access_tokens_access_token", "access_token"),
        Index("idx_personal_access_tokens_refresh_token", "refresh_token"),
    )

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<UserToken(id={self.id}, user_id={self.user_id}, user_type='{self.user_type}', is_active={self.is_active})>"
