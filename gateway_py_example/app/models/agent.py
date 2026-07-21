from sqlalchemy import Column, String, Boolean, BigInteger, DateTime, ForeignKey, Integer, Table, Text
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from app.core.database import Base


class Agent(Base):
    __table__ = Table(
        'agents',
        Base.metadata,
        Column('id', BigInteger, primary_key=True, index=True),
        Column('company_id', BigInteger, ForeignKey('companies.id'), nullable=False, index=True),
        Column('name', String(255), nullable=False),
        Column('description', Text, nullable=False),
        Column('schedule_type', String(20), nullable=False),
        Column('schedule_time', String(10), nullable=False),
        Column('schedule_day', String(20), nullable=True),
        Column('schedule_date', Integer, nullable=True),
        Column('schedule', String(255), nullable=False),
        Column('status', String(20), nullable=False, default='active'),
        Column('last_run_at', DateTime(timezone=True), nullable=True),
        Column('next_run_at', DateTime(timezone=True), nullable=True),
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Relationships
    company = relationship("Company", backref="agents")
    reports = relationship("AgentReport", back_populates="agent", cascade="all, delete-orphan")

    def to_dict(self):
        result = {c.name: getattr(self, c.name) for c in self.__table__.columns}
        if self.reports:
            result['reports'] = [report.report_name for report in self.reports]
        return result

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<Agent(id={self.id}, name='{self.name}')>"


class AgentReport(Base):
    __table__ = Table(
        'agent_reports',
        Base.metadata,
        Column('id', BigInteger, primary_key=True, index=True),
        Column('agent_id', BigInteger, ForeignKey('agents.id'), nullable=False, index=True),
        Column('report_name', String(255), nullable=False),
        Column('created_at', DateTime(timezone=True), server_default=func.now(), nullable=False),
        Column('updated_at', DateTime(timezone=True), onupdate=func.now(), server_default=func.now(), nullable=False),
        Column('deleted_at', DateTime(timezone=True), nullable=True),
    )

    # Relationships
    agent = relationship("Agent", back_populates="reports")

    def to_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def __repr__(self):
        return f"<AgentReport(id={self.id}, report_name='{self.report_name}')>"
