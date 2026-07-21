from sqlalchemy import Column, String, Boolean, Text, BigInteger, ForeignKey, Index, DateTime, Table
from sqlalchemy.sql import func
from app.core.database import Base


class Setting(Base):
    __table__ = Table(
        'settings',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Foreign keys
        Column('company_id', BigInteger, ForeignKey("companies.id", ondelete="CASCADE"), nullable=False, index=True),
        # Required fields
        Column('key', String(255), nullable=False, index=True),
        Column('value', Text, nullable=False),
        Column('category', String(100), nullable=False, default="general", index=True),
        # Optional fields
        Column('description', Text, nullable=True),
        # Boolean flags
        Column('is_active', Boolean, default=True, nullable=False, index=True),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Note: No relationship defined here because Company model already has a "settings" JSONB column
    # If needed, can be accessed via: db.query(Setting).filter(Setting.company_id == company_id)

    # Unique constraint on (company_id, key)
    __table_args__ = (
        Index("uk_settings_company_key", "company_id", "key", unique=True),
    )

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<Setting(id={self.id}, company_id={self.company_id}, key='{self.key}')>"
