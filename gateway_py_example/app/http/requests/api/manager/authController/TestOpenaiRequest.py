from pydantic import BaseModel, Field


class TestOpenAIRequest(BaseModel):
    message: str = Field(..., description="User message for OpenAI test")

    class Config:
        json_schema_extra = {
            "example": {
                "message": "Bu ay kaç çocuk kalmış?"
            }
        }
