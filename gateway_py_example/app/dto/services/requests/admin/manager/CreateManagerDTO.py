from typing import Dict, Optional
from app.http.requests.api.admin.managerController.CreateManagerRequest import CreateManagerRequest


class CreateManagerDTO:
    def __init__(self, currentUser: Dict, request: CreateManagerRequest):
        self._companyId = currentUser.get("company_id")
        self._email = request.email
        self._username = request.username
        self._password = request.password
        self._firstName = request.firstName if hasattr(request, "firstName") else None
        self._lastName = request.lastName if hasattr(request, "lastName") else None
        self._isActive = getattr(request, "isActive", True)

    def getCompanyId(self) -> int:
        return self._companyId

    def getEmail(self) -> str:
        return self._email

    def getUsername(self) -> str:
        return self._username

    def getPassword(self) -> str:
        return self._password

    def getFirstName(self) -> Optional[str]:
        return self._firstName

    def getLastName(self) -> Optional[str]:
        return self._lastName

    def getIsActive(self) -> bool:
        return self._isActive if self._isActive is not None else True
