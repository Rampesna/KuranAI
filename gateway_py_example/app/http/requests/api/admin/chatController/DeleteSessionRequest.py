from pydantic import BaseModel, Field


class DeleteSessionRequest(BaseModel):
    sessionId: int = Field(..., alias="session_id")
