from pydantic import BaseModel, Field, ConfigDict
from typing import Optional


class ListAgentsRequest(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    page: Optional[int] = Field(1, ge=1, alias="page")
    perPage: Optional[int] = Field(20, ge=1, le=100, alias="per_page")
