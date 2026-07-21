from typing import Dict, Optional
from app.http.requests.api.manager.agentController.ListAgentsRequest import ListAgentsRequest


class ListAgentsDTO:
    def __init__(self, currentUser: Dict, request: ListAgentsRequest):
        self._companyId = currentUser.get("company_id")
        self._page = request.page or 1
        self._perPage = request.perPage or 20

    def getCompanyId(self) -> int:
        return self._companyId

    def getPage(self) -> int:
        return self._page

    def getPerPage(self) -> int:
        return self._perPage
