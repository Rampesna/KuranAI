from pydantic import BaseModel, Field, ConfigDict
from typing import Optional


class UpdateManagerRequest(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    firstName: Optional[str] = Field(None, max_length=100, alias="firstName")
    lastName: Optional[str] = Field(None, max_length=100, alias="lastName")
    isActive: Optional[bool] = Field(None, alias="isActive")
    password: Optional[str] = Field(None, min_length=1, alias="password")
