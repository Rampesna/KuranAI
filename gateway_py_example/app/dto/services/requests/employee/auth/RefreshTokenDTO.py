from app.http.requests.api.employee.authController.RefreshTokenRequest import RefreshTokenRequest


class RefreshTokenDTO:
    def __init__(self, request: RefreshTokenRequest):
        self._refreshToken = request.refresh_token

    def getRefreshToken(self) -> str:
        return self._refreshToken
