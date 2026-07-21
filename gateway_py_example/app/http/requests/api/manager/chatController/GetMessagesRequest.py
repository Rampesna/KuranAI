from pydantic import BaseModel, Field
from typing import Optional


class GetMessagesRequest(BaseModel):
    sessionId: int = Field(..., alias="session_id")
    pageIndex: Optional[int] = Field(0, ge=0, alias="pageIndex")
    pageSize: Optional[int] = Field(10, ge=1, le=100, alias="pageSize")
