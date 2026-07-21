from pydantic import BaseModel, Field
from typing import Optional


class GetMessagesRequest(BaseModel):
    sessionId: int = Field(..., alias="session_id")
    page: Optional[int] = Field(1, ge=1, alias="page")
    perPage: Optional[int] = Field(50, ge=1, le=100, alias="per_page")
