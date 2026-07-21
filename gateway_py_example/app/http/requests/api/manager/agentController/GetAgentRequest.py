from pydantic import BaseModel, Field


class GetAgentRequest(BaseModel):
    agentId: int = Field(..., alias="agent_id")
