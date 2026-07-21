from typing import Dict
from app.http.requests.api.admin.authController.GetCurrentUserRequest import GetCurrentUserRequest


class GetCurrentUserDTO:
    def __init__(self, currentUser: Dict, request: GetCurrentUserRequest):
        self._userId = currentUser.get("user_id") or currentUser.get("id")

    def getUserId(self) -> int:
        return self._userId
