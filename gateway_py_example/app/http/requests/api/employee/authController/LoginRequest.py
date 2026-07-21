from pydantic import BaseModel, Field


class LoginRequest(BaseModel):
    credential: str = Field(..., description="Email or username")
    password: str = Field(..., description="Password")

    class Config:
        json_schema_extra = {
            "example": {
                "credential": "employee@intelor.com",
                "password": "admin123"
            }
        }
