from sqlalchemy import Column, String, Boolean, BigInteger, DateTime, ForeignKey, Table
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from app.core.database import Base


class Manager(Base):
    __table__ = Table(
        'managers',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Foreign keys
        Column('company_id', BigInteger, ForeignKey('companies.id'), nullable=False, index=True),
        # Required/unique fields
        Column('email', String(255), unique=True, nullable=False, index=True),
        Column('username', String(100), unique=True, nullable=False, index=True),
        Column('password_hash', String(255), nullable=False),
        # Optional user info
        Column('first_name', String(100), nullable=True),
        Column('last_name', String(100), nullable=True),
        # Boolean flags
        Column('is_active', Boolean, default=True, nullable=False),
        # Custom timestamps
        Column('last_login_at', DateTime(timezone=True), nullable=True),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Relationships
    company = relationship("Company", backref="managers")

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<Manager(id={self.id}, email='{self.email}')>"
