from pydantic import BaseModel, Field, validator
from typing import List, Optional


class CreateAgentRequest(BaseModel):
    name: str = Field(..., min_length=1, max_length=255)
    description: str = Field(..., min_length=1)
    schedule_type: str = Field(..., pattern="^(daily|weekly|monthly|custom)$")
    schedule_time: str = Field(..., pattern="^\\d{2}:\\d{2}$")
    schedule_day: Optional[str] = None
    schedule_date: Optional[int] = Field(None, ge=1, le=31)
    reports: List[str] = Field(..., min_length=1)
    status: Optional[str] = Field("active", pattern="^(active|paused)$")

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
                "name": "Günlük Gelir Raporu",
                "description": "Her gün sabah 08:00'de toplam gelir raporunu hazırlar",
                "schedule_type": "daily",
                "schedule_time": "08:00",
                "reports": ["Toplam Gelir", "Oda Bazlı Gelir"],
                "status": "active"
            }
        }
