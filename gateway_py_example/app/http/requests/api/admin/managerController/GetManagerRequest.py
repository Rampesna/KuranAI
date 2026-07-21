from pydantic import BaseModel


class GetManagerRequest(BaseModel):
    managerId: int = 0
