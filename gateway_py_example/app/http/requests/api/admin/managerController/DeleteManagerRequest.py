from pydantic import BaseModel


class DeleteManagerRequest(BaseModel):
    managerId: int = 0
