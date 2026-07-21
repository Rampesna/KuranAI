from typing import Dict
from app.http.requests.api.admin.managerController.ListManagersRequest import ListManagersRequest


class ListManagersDTO:
    def __init__(self, currentUser: Dict, request: ListManagersRequest):
        self._companyId = currentUser.get("company_id")
        self._pageIndex = request.pageIndex if request.pageIndex is not None else 0
        self._pageSize = request.pageSize if request.pageSize is not None else 20

    def getCompanyId(self) -> int:
        return self._companyId

    def getPageIndex(self) -> int:
        return self._pageIndex

    def getPageSize(self) -> int:
        return self._pageSize
