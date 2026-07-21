from typing import Dict, Optional
from app.http.requests.api.manager.chatController.CreateSessionRequest import CreateSessionRequest


class CreateSessionDTO:
    def __init__(self, currentUser: Dict, request: CreateSessionRequest):
        self._companyId = currentUser.get("company_id")
        self._userType = currentUser.get("user_type")
        self._userId = currentUser.get("id")
        self._title = request.title

    def getCompanyId(self) -> int:
        return self._companyId

    def getUserType(self) -> str:
        return self._userType

    def getUserId(self) -> int:
        return self._userId

    def getTitle(self) -> Optional[str]:
        return self._title
