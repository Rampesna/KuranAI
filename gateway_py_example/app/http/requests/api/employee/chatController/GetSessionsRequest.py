from pydantic import BaseModel, Field
from typing import Optional


class GetSessionsRequest(BaseModel):
    page: Optional[int] = Field(1, ge=1, alias="page")
    perPage: Optional[int] = Field(20, ge=1, le=100, alias="per_page")
