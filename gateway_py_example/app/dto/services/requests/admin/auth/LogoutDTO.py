from typing import Dict
from app.http.requests.api.admin.authController.LogoutRequest import LogoutRequest


class LogoutDTO:
    def __init__(self, currentUser: Dict, request: LogoutRequest):
        self._accessToken = currentUser.get("access_token", "")
        self._currentUser = currentUser

    def getAccessToken(self) -> str:
        return self._accessToken

    def getCurrentUser(self) -> Dict:
        return self._currentUser
