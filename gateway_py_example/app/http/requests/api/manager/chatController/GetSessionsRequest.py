from pydantic import BaseModel, Field
from typing import Optional


class GetSessionsRequest(BaseModel):
    pageIndex: Optional[int] = Field(0, ge=0, description="Page index (0-based)")
    pageSize: Optional[int] = Field(50, ge=1, le=100, description="Page size")
