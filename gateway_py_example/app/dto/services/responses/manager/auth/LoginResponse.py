from typing import Optional, Dict


class LoginSuccessResponse:
    def __init__(self, access_token: str, refresh_token: str, user: Dict):
        self._access_token = access_token
        self._refresh_token = refresh_token
        self._user = user

    def to_dict(self) -> Dict:
        return {
            "access_token": self._access_token,
            "refresh_token": self._refresh_token,
            "user": self._user
        }


class LoginInvalidCredentialsResponse:
    def to_dict(self) -> Dict:
        return None
