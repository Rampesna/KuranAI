from typing import Dict
from app.http.requests.api.employee.chatController.DeleteSessionRequest import DeleteSessionRequest


class DeleteSessionDTO:
    def __init__(self, currentUser: Dict, request: DeleteSessionRequest):
        self._companyId = currentUser.get("company_id")
        self._userType = currentUser.get("user_type")
        self._userId = currentUser.get("id")
        self._sessionId = request.sessionId

    def getCompanyId(self) -> int:
        return self._companyId

    def getUserType(self) -> str:
        return self._userType

    def getUserId(self) -> int:
        return self._userId

    def getSessionId(self) -> int:
        return self._sessionId
