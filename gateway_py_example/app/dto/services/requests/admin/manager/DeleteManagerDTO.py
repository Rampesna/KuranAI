from typing import Dict
from app.http.requests.api.admin.managerController.DeleteManagerRequest import DeleteManagerRequest


class DeleteManagerDTO:
    def __init__(self, currentUser: Dict, request: DeleteManagerRequest):
        self._companyId = currentUser.get("company_id")
        self._managerId = request.managerId

    def getCompanyId(self) -> int:
        return self._companyId

    def getManagerId(self) -> int:
        return self._managerId
