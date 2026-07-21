from typing import Dict
from app.http.requests.api.employee.chatController.GetMessagesRequest import GetMessagesRequest


class GetMessagesDTO:
    def __init__(self, currentUser: Dict, request: GetMessagesRequest):
        self._companyId = currentUser.get("company_id")
        self._sessionId = request.sessionId
        # Convert page/per_page to pageIndex/pageSize (pageIndex is 0-based)
        page = request.page or 1
        self._pageIndex = page - 1 if page > 0 else 0
        self._pageSize = request.perPage or 50

    def getCompanyId(self) -> int:
        return self._companyId

    def getSessionId(self) -> int:
        return self._sessionId

    def getPageIndex(self) -> int:
        return self._pageIndex

    def getPageSize(self) -> int:
        return self._pageSize
