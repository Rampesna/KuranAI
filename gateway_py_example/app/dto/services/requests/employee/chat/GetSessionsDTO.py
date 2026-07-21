from typing import Dict, Optional
from app.http.requests.api.employee.chatController.GetSessionsRequest import GetSessionsRequest


class GetSessionsDTO:
    def __init__(self, currentUser: Dict, request: GetSessionsRequest):
        self._companyId = currentUser.get("company_id")
        self._userType = currentUser.get("user_type")
        self._userId = currentUser.get("id")
        # Convert page/per_page to pageIndex/pageSize (pageIndex is 0-based)
        page = request.page or 1
        self._pageIndex = page - 1 if page > 0 else 0
        self._pageSize = request.perPage or 20

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
