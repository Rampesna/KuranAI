from pydantic import BaseModel, Field


class DeleteAgentRequest(BaseModel):
    agentId: int = Field(..., alias="agent_id")
