from datetime import datetime, timedelta
from typing import Optional, Dict
from jose import JWTError, jwt
from app.core.config import settings


class JWTService:
    def __init__(self):
        self.secret = settings.JWT_SECRET
        self.access_token_expiry = timedelta(hours=settings.JWT_ACCESS_TOKEN_EXPIRY_HOURS)
        self.refresh_token_expiry = timedelta(hours=settings.JWT_REFRESH_TOKEN_EXPIRY_HOURS)

    def generateTokenPair(self, userId: int, userType: str, email: str) -> Dict[str, str]:
        accessToken = self._generateToken(
            userId=userId,
            userType=userType,
            email=email,
            expiry=self.access_token_expiry,
            tokenType="access"
        )
        refreshToken = self._generateToken(
            userId=userId,
            userType=userType,
            email=email,
            expiry=self.refresh_token_expiry,
            tokenType="refresh"
        )
        return {
            "access_token": accessToken,
            "refresh_token": refreshToken
        }

    def _generateToken(self, userId: int, userType: str, email: str, expiry: timedelta, tokenType: str) -> str:
        payload = {
            "user_id": userId,
            "user_type": userType,
            "email": email,
            "token_type": tokenType,
            "exp": datetime.utcnow() + expiry,
            "iat": datetime.utcnow()
        }
        return jwt.encode(payload, self.secret, algorithm="HS256")

    def validateToken(self, token: str) -> Optional[Dict]:
        try:
            payload = jwt.decode(token, self.secret, algorithms=["HS256"])
            return payload
        except JWTError:
            return None

    def refreshAccessToken(self, refreshToken: str) -> Optional[Dict[str, str]]:
        payload = self.validateToken(refreshToken)
        if not payload or payload.get("token_type") != "refresh":
            return None

        return {
            "access_token": self._generateToken(
                userId=payload["user_id"],
                userType=payload["user_type"],
                email=payload["email"],
                expiry=self.access_token_expiry,
                tokenType="access"
            )
        }
