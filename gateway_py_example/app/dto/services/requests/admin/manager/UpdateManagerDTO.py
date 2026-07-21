from typing import Dict, Optional
from app.http.requests.api.admin.managerController.UpdateManagerRequest import UpdateManagerRequest


class UpdateManagerDTO:
    def __init__(self, currentUser: Dict, request: UpdateManagerRequest, managerId: int):
        self._companyId = currentUser.get("company_id")
        self._managerId = managerId
        self._firstName = getattr(request, "firstName", None)
        self._lastName = getattr(request, "lastName", None)
        self._isActive = getattr(request, "isActive", None)
        self._password = getattr(request, "password", None)

    def getCompanyId(self) -> int:
        return self._companyId

    def getManagerId(self) -> int:
        return self._managerId

    def getFirstName(self) -> Optional[str]:
        return self._firstName

    def getLastName(self) -> Optional[str]:
        return self._lastName

    def getIsActive(self) -> Optional[bool]:
        return self._isActive

    def getPassword(self) -> Optional[str]:
        return self._password
