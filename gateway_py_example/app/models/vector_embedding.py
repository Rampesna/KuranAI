from sqlalchemy import Column, String, Text, BigInteger, DateTime, Table
from sqlalchemy.sql import func
from pgvector.sqlalchemy import Vector
from app.core.database import Base


class VectorEmbedding(Base):
    __table__ = Table(
        'vector_embeddings',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Required fields
        Column('view_name', String(255), nullable=False, index=True),
        Column('column_name', String(255), nullable=True),
        Column('data_type', String(100), nullable=True),
        Column('description', Text, nullable=True),
        Column('semantic_meaning', Text, nullable=True),
        Column('embedding_text', Text, nullable=True),
        # Vector embedding (3072 dimensions for text-embedding-3-large)
        Column('embedding', Vector(3072), nullable=False),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
    )

    def __repr__(self):
        return f"<VectorEmbedding(id={self.id}, view_name='{self.view_name}')>"
