from pydantic import BaseModel, Field


class StreamMessageRequest(BaseModel):
    sessionId: int = Field(0, alias="session_id")
    content: str = Field(..., min_length=1, max_length=10000, description="Message content")
