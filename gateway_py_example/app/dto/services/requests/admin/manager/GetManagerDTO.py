from typing import Dict
from app.http.requests.api.admin.managerController.GetManagerRequest import GetManagerRequest


class GetManagerDTO:
    def __init__(self, currentUser: Dict, request: GetManagerRequest):
        self._companyId = currentUser.get("company_id")
        self._managerId = request.managerId

    def getCompanyId(self) -> int:
        return self._companyId

    def getManagerId(self) -> int:
        return self._managerId
