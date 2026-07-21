from typing import Dict
from app.http.requests.api.employee.chatController.SendMessageRequest import SendMessageRequest


class SendMessageDTO:
    def __init__(self, currentUser: Dict, request: SendMessageRequest):
        self._companyId = currentUser.get("company_id")
        self._sessionId = request.sessionId
        self._content = request.content

    def getCompanyId(self) -> int:
        return self._companyId

    def getSessionId(self) -> int:
        return self._sessionId

    def getContent(self) -> str:
        return self._content
