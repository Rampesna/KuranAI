from pydantic import BaseModel, Field, ConfigDict
from typing import Optional


class ListManagersRequest(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    pageIndex: Optional[int] = Field(0, ge=0, alias="pageIndex")
    pageSize: Optional[int] = Field(20, ge=1, le=100, alias="pageSize")
