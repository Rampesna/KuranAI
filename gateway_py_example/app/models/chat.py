from sqlalchemy import Column, String, Text, BigInteger, Integer, ForeignKey, DateTime, Table
from sqlalchemy.orm import relationship
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.sql import func
from app.core.database import Base


class ChatSession(Base):
    __table__ = Table(
        'chat_sessions',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Foreign keys
        Column('company_id', BigInteger, ForeignKey('companies.id'), nullable=False, index=True),
        # Required fields
        Column('user_type', String(50), nullable=False, index=True),  # admin, manager, employee
        Column('user_id', BigInteger, nullable=False, index=True),
        # Optional fields
        Column('title', String(255), nullable=True),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Relationships
    company = relationship("Company", backref="chat_sessions")
    messages = relationship("ChatMessage", back_populates="session", cascade="all, delete-orphan")

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<ChatSession(id={self.id}, user_type='{self.user_type}')>"


class ChatMessage(Base):
    __table__ = Table(
        'chat_messages',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Foreign keys
        Column('chat_session_id', BigInteger, ForeignKey('chat_sessions.id'), nullable=False, index=True),
        # Required fields
        Column('role', String(50), nullable=False),  # user, assistant, system
        Column('content', Text, nullable=False),
        # Optional fields
        Column('token_count', Integer, nullable=True),
        Column('generated_sql', Text, nullable=True),
        # JSON fields
        Column('sql_result', JSONB, nullable=True),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Relationships
    session = relationship("ChatSession", back_populates="messages")

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<ChatMessage(id={self.id}, role='{self.role}')>"
