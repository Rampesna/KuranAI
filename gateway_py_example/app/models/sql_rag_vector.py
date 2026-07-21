from sqlalchemy import Column, String, Text, BigInteger, ForeignKey, DateTime, Table
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from pgvector.sqlalchemy import Vector
from app.core.database import Base


class SQLGenerationRAGVector(Base):
    __table__ = Table(
        'sql_generation_rag_vectors',
        Base.metadata,
        # Primary key first
        Column('id', BigInteger, primary_key=True, index=True),
        # Foreign keys
        Column('company_id', BigInteger, ForeignKey('companies.id'), nullable=False, index=True),
        # Required fields
        Column('schema_name', String(255), nullable=False),
        Column('schema_content', Text, nullable=False),
        Column('embedding', Vector(3072), nullable=False),
        # Timestamps (always last)
        Column('created_at', DateTime(timezone=True), server_default=func.now()),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now()),
    )

    # Relationship
    company = relationship("Company", backref="sql_rag_vectors")

    def __repr__(self):
        return f"<SQLGenerationRAGVector(id={self.id}, schema_name='{self.schema_name}')>"
