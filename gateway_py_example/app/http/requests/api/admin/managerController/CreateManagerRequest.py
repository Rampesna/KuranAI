from pydantic import BaseModel, Field, ConfigDict
from typing import Optional


class CreateManagerRequest(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    email: str = Field(..., min_length=1, max_length=255)
    username: str = Field(..., min_length=1, max_length=100)
    password: str = Field(..., min_length=1)
    firstName: Optional[str] = Field(None, max_length=100, alias="firstName")
    lastName: Optional[str] = Field(None, max_length=100, alias="lastName")
    isActive: Optional[bool] = Field(True, alias="isActive")
