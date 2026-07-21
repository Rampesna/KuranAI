from pydantic import BaseModel, Field
from typing import Optional


class CreateSessionRequest(BaseModel):
    title: Optional[str] = Field(None, max_length=255, description="Session title")
