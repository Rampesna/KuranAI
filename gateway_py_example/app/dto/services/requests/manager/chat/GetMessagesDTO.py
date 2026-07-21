from typing import Dict
from app.http.requests.api.manager.chatController.GetMessagesRequest import GetMessagesRequest


class GetMessagesDTO:
    def __init__(self, currentUser: Dict, request: GetMessagesRequest):
        self._companyId = currentUser.get("company_id")
        self._sessionId = request.sessionId
        self._pageIndex = request.pageIndex or 0
        self._pageSize = request.pageSize or 10

    def getCompanyId(self) -> int:
        return self._companyId

    def getSessionId(self) -> int:
        return self._sessionId

    def getPageIndex(self) -> int:
        return self._pageIndex

    def getPageSize(self) -> int:
        return self._pageSize
