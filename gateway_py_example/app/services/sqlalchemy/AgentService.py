from sqlalchemy.orm import Session
from sqlalchemy import and_, func
from datetime import datetime, timedelta
from typing import List, Optional, Dict, Any
from app.models.agent import Agent, AgentReport
from app.core.ServiceResponse import ServiceResponse
from app.dto.services.requests.manager.agent.ListAgentsDTO import ListAgentsDTO
from app.dto.services.requests.manager.agent.GetAgentDTO import GetAgentDTO
from app.dto.services.requests.manager.agent.CreateAgentDTO import CreateAgentDTO
from app.dto.services.requests.manager.agent.UpdateAgentDTO import UpdateAgentDTO
from app.dto.services.requests.manager.agent.DeleteAgentDTO import DeleteAgentDTO


class AgentService:
    def __init__(self, db: Session):
        self.db = db

    def listAgents(self, dto: ListAgentsDTO) -> ServiceResponse:
        """List agents with pagination"""
        try:
            offset = (dto.getPage() - 1) * dto.getPerPage()
            
            # Query agents
            query = self.db.query(Agent).filter(
                and_(
                    Agent.company_id == dto.getCompanyId(),
                    Agent.deleted_at.is_(None)
                )
            )
            
            total = query.count()
            agents = query.order_by(Agent.created_at.desc()).offset(offset).limit(dto.getPerPage()).all()
            
            # Load reports for each agent
            for agent in agents:
                agent.reports = self.db.query(AgentReport).filter(
                    and_(
                        AgentReport.agent_id == agent.id,
                        AgentReport.deleted_at.is_(None)
                    )
                ).all()
            
            items = [agent.to_dict() for agent in agents]
            
            totalPages = (total + dto.getPerPage() - 1) // dto.getPerPage() if total > 0 else 0
            
            data = {
                "items": items,
                "meta": {
                    "current_page": dto.getPage(),
                    "per_page": dto.getPerPage(),
                    "total": total,
                    "total_pages": totalPages
                }
            }
            
            return ServiceResponse(
                isSuccess=True,
                message="Agents retrieved successfully",
                statusCode=200,
                data=data
            )
        except Exception as e:
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to retrieve agents: {str(e)}",
                statusCode=500,
                data=None
            )

    def getAgent(self, dto: GetAgentDTO) -> ServiceResponse:
        """Get a single agent by ID"""
        try:
            agent = self.db.query(Agent).filter(
                and_(
                    Agent.id == dto.getAgentId(),
                    Agent.company_id == dto.getCompanyId(),
                    Agent.deleted_at.is_(None)
                )
            ).first()
            
            if not agent:
                return ServiceResponse(
                    isSuccess=False,
                    message="Agent not found",
                    statusCode=404,
                    data=None
                )
            
            # Load reports
            agent.reports = self.db.query(AgentReport).filter(
                and_(
                    AgentReport.agent_id == agent.id,
                    AgentReport.deleted_at.is_(None)
                )
            ).all()
            
            return ServiceResponse(
                isSuccess=True,
                message="Agent retrieved successfully",
                statusCode=200,
                data=agent.to_dict()
            )
        except Exception as e:
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to retrieve agent: {str(e)}",
                statusCode=500,
                data=None
            )

    def createAgent(self, dto: CreateAgentDTO) -> ServiceResponse:
        """Create a new agent"""
        try:
            # Build schedule string
            schedule = self._buildScheduleString(
                dto.getScheduleType(),
                dto.getScheduleTime(),
                dto.getScheduleDay(),
                dto.getScheduleDate()
            )
            
            # Calculate next run time
            nextRunAt = self._calculateNextRun(
                dto.getScheduleType(),
                dto.getScheduleTime(),
                dto.getScheduleDay(),
                dto.getScheduleDate()
            )
            
            # Create agent
            agent = Agent(
                company_id=dto.getCompanyId(),
                name=dto.getName(),
                description=dto.getDescription(),
                schedule_type=dto.getScheduleType(),
                schedule_time=dto.getScheduleTime(),
                schedule_day=dto.getScheduleDay(),
                schedule_date=dto.getScheduleDate(),
                schedule=schedule,
                status=dto.getStatus(),
                next_run_at=nextRunAt
            )
            
            self.db.add(agent)
            self.db.flush()  # Get the agent ID
            
            # Create reports
            reports = dto.getReports()
            for reportName in reports:
                agentReport = AgentReport(
                    agent_id=agent.id,
                    report_name=reportName
                )
                self.db.add(agentReport)
            
            self.db.commit()
            
            # Reload with reports
            agent.reports = self.db.query(AgentReport).filter(
                and_(
                    AgentReport.agent_id == agent.id,
                    AgentReport.deleted_at.is_(None)
                )
            ).all()
            
            return ServiceResponse(
                isSuccess=True,
                message="Agent created successfully",
                statusCode=201,
                data=agent.to_dict()
            )
        except Exception as e:
            self.db.rollback()
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to create agent: {str(e)}",
                statusCode=500,
                data=None
            )

    def updateAgent(self, dto: UpdateAgentDTO) -> ServiceResponse:
        """Update an existing agent"""
        try:
            agent = self.db.query(Agent).filter(
                and_(
                    Agent.id == dto.getAgentId(),
                    Agent.company_id == dto.getCompanyId(),
                    Agent.deleted_at.is_(None)
                )
            ).first()
            
            if not agent:
                return ServiceResponse(
                    isSuccess=False,
                    message="Agent not found",
                    statusCode=404,
                    data=None
                )
            
            data = dto.toDict()
            
            # Update fields
            if 'name' in data:
                agent.name = data['name']
            if 'description' in data:
                agent.description = data['description']
            if 'schedule_type' in data:
                agent.schedule_type = data['schedule_type']
            if 'schedule_time' in data:
                agent.schedule_time = data['schedule_time']
            if 'schedule_day' in data:
                agent.schedule_day = data['schedule_day']
            if 'schedule_date' in data:
                agent.schedule_date = data['schedule_date']
            if 'status' in data:
                agent.status = data['status']
            
            # Rebuild schedule if schedule-related fields changed
            if any(key in data for key in ['schedule_type', 'schedule_time', 'schedule_day', 'schedule_date']):
                agent.schedule = self._buildScheduleString(
                    agent.schedule_type,
                    agent.schedule_time,
                    agent.schedule_day,
                    agent.schedule_date
                )
                agent.next_run_at = self._calculateNextRun(
                    agent.schedule_type,
                    agent.schedule_time,
                    agent.schedule_day,
                    agent.schedule_date
                )
            
            # Update reports if provided
            if 'reports' in data:
                # Soft delete existing reports
                self.db.query(AgentReport).filter(
                    and_(
                        AgentReport.agent_id == agent.id,
                        AgentReport.deleted_at.is_(None)
                    )
                ).update({'deleted_at': datetime.utcnow()})
                
                # Create new reports
                for reportName in data['reports']:
                    agentReport = AgentReport(
                        agent_id=agent.id,
                        report_name=reportName
                    )
                    self.db.add(agentReport)
            
            self.db.commit()
            
            # Reload with reports
            agent.reports = self.db.query(AgentReport).filter(
                and_(
                    AgentReport.agent_id == agent.id,
                    AgentReport.deleted_at.is_(None)
                )
            ).all()
            
            return ServiceResponse(
                isSuccess=True,
                message="Agent updated successfully",
                statusCode=200,
                data=agent.to_dict()
            )
        except Exception as e:
            self.db.rollback()
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to update agent: {str(e)}",
                statusCode=500,
                data=None
            )

    def deleteAgent(self, dto: DeleteAgentDTO) -> ServiceResponse:
        """Soft delete an agent"""
        try:
            agent = self.db.query(Agent).filter(
                and_(
                    Agent.id == dto.getAgentId(),
                    Agent.company_id == dto.getCompanyId(),
                    Agent.deleted_at.is_(None)
                )
            ).first()
            
            if not agent:
                return ServiceResponse(
                    isSuccess=False,
                    message="Agent not found",
                    statusCode=404,
                    data=None
                )
            
            # Soft delete
            agent.deleted_at = datetime.utcnow()
            self.db.commit()
            
            return ServiceResponse(
                isSuccess=True,
                message="Agent deleted successfully",
                statusCode=200,
                data=None
            )
        except Exception as e:
            self.db.rollback()
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to delete agent: {str(e)}",
                statusCode=500,
                data=None
            )

    def _buildScheduleString(self, scheduleType: str, scheduleTime: str, 
                               scheduleDay: Optional[str] = None, 
                               scheduleDate: Optional[int] = None) -> str:
        """Build human-readable schedule string"""
        if scheduleType == 'daily':
            return f"Her gün {scheduleTime}"
        elif scheduleType == 'weekly':
            day = scheduleDay or 'Pazartesi'
            return f"Her {day} {scheduleTime}"
        elif scheduleType == 'monthly':
            date = scheduleDate or 1
            return f"Her ayın {date} {scheduleTime}"
        else:
            return scheduleTime

    def _calculateNextRun(self, scheduleType: str, scheduleTime: str,
                            scheduleDay: Optional[str] = None,
                            scheduleDate: Optional[int] = None) -> Optional[datetime]:
        """Calculate next run time"""
        try:
            now = datetime.utcnow()
            timeParts = scheduleTime.split(':')
            hour = int(timeParts[0])
            minute = int(timeParts[1]) if len(timeParts) > 1 else 0
            
            # For simplicity, set next run to tomorrow at the specified time
            # In production, this should be more sophisticated
            nextRun = now + timedelta(days=1)
            nextRun = nextRun.replace(hour=hour, minute=minute, second=0, microsecond=0)
            
            return nextRun
        except:
            return None
