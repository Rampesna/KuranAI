from typing import Dict, Optional
from app.http.requests.api.manager.chatController.GetSessionsRequest import GetSessionsRequest


class GetSessionsDTO:
    def __init__(self, currentUser: Dict, request: GetSessionsRequest):
        self._companyId = currentUser.get("company_id")
        self._userType = currentUser.get("user_type")
        self._userId = currentUser.get("id")
        self._pageIndex = request.pageIndex or 0
        self._pageSize = request.pageSize or 50

    def getCompanyId(self) -> int:
        return self._companyId

    def getUserType(self) -> str:
        return self._userType

    def getUserId(self) -> int:
        return self._userId

    def getPageIndex(self) -> int:
        return self._pageIndex

    def getPageSize(self) -> int:
        return self._pageSize
