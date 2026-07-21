from typing import Optional
from app.http.requests.api.manager.authController.LoginRequest import LoginRequest


class LoginDTO:
    def __init__(self, request: LoginRequest):
        self._credential = request.credential
        self._password = request.password

    def getCredential(self) -> str:
        return self._credential

    def getPassword(self) -> str:
        return self._password
