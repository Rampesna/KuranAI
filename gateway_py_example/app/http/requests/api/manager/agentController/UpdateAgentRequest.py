from pydantic import BaseModel, Field, validator
from typing import List, Optional


class UpdateAgentRequest(BaseModel):
    agentId: Optional[int] = Field(None, alias="agent_id")
    name: Optional[str] = Field(None, min_length=1, max_length=255)
    description: Optional[str] = Field(None, min_length=1)
    schedule_type: Optional[str] = Field(None, pattern="^(daily|weekly|monthly|custom)$")
    schedule_time: Optional[str] = Field(None, pattern="^\\d{2}:\\d{2}$")
    schedule_day: Optional[str] = None
    schedule_date: Optional[int] = Field(None, ge=1, le=31)
    reports: Optional[List[str]] = Field(None, min_length=1)
    status: Optional[str] = Field(None, pattern="^(active|paused)$")

    @validator('schedule_day')
    def validate_schedule_day(cls, v, values):
        if values.get('schedule_type') == 'weekly' and not v:
            raise ValueError('schedule_day is required for weekly schedule')
        return v

    @validator('schedule_date')
    def validate_schedule_date(cls, v, values):
        if values.get('schedule_type') == 'monthly' and not v:
            raise ValueError('schedule_date is required for monthly schedule')
        return v

    class Config:
        json_schema_extra = {
            "example": {
                "name": "Updated Agent Name",
                "status": "paused"
            }
        }
