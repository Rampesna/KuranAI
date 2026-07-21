from sqlalchemy.orm import Session
from sqlalchemy import or_
from datetime import datetime, timedelta, timezone
from fastapi import Request
from app.models.admin import Admin
from app.services.common.jwt import JWTService
from app.services.common.password import PasswordService
from app.services.common.TokenService import TokenService
from app.core.ServiceResponse import ServiceResponse
from app.core.config import settings
from app.utils.device_info import extract_device_info
from app.dto.services.requests.admin.auth.LoginDTO import LoginDTO
from app.dto.services.requests.admin.auth.RefreshTokenDTO import RefreshTokenDTO
from app.dto.services.requests.admin.auth.GetCurrentUserDTO import GetCurrentUserDTO
from app.dto.services.requests.admin.auth.LogoutDTO import LogoutDTO
from app.dto.services.responses.admin.auth.LoginResponse import LoginSuccessResponse, LoginInvalidCredentialsResponse
from typing import Optional, Dict, Any
from sqlalchemy.orm import Session


class AdminService:
    def __init__(self, db: Session, jwtService: JWTService):
        self.db = db
        self.jwtService = jwtService
        self.passwordService = PasswordService()
        self.tokenService = TokenService(db)

    def login(self, dto: LoginDTO, httpRequest: Request) -> ServiceResponse:
        deviceInfo = extract_device_info(httpRequest)
        
        admin = self.db.query(Admin).filter(
            or_(
                Admin.email == dto.getCredential(),
                Admin.username == dto.getCredential()
            ),
            Admin.deleted_at.is_(None)
        ).first()

        if not admin:
            return ServiceResponse(
                isSuccess=False,
                message="Invalid credentials",
                statusCode=401,
                data=None
            )

        if not admin.is_active:
            return ServiceResponse(
                isSuccess=False,
                message="Account is inactive",
                statusCode=403,
                data=None
            )

        if not self.passwordService.check(dto.getPassword(), admin.password_hash):
            return ServiceResponse(
                isSuccess=False,
                message="Invalid credentials",
                statusCode=401,
                data=None
            )

        # Update last login
        admin.last_login_at = datetime.now(timezone.utc)
        self.db.commit()

        # Generate tokens
        tokens = self.jwtService.generateTokenPair(
            userId=admin.id,
            userType="admin",
            email=admin.email
        )

        # Calculate expiry times (timezone-aware)
        now = datetime.now(timezone.utc)
        accessTokenExpiresAt = now + timedelta(hours=settings.JWT_ACCESS_TOKEN_EXPIRY_HOURS)
        refreshTokenExpiresAt = now + timedelta(hours=settings.JWT_REFRESH_TOKEN_EXPIRY_HOURS)

        # Save tokens to database and Redis
        self.tokenService.saveToken(
            userId=admin.id,
            userType="admin",
            accessToken=tokens["access_token"],
            refreshToken=tokens["refresh_token"],
            accessTokenExpiresAt=accessTokenExpiresAt,
            refreshTokenExpiresAt=refreshTokenExpiresAt,
            deviceInfo=deviceInfo
        )

        userData = {
            "id": admin.id,
            "email": admin.email,
            "username": admin.username,
            "first_name": admin.first_name,
            "last_name": admin.last_name,
            "user_type": "admin",
            "company_id": admin.company_id
        }

        responseData = LoginSuccessResponse(
            access_token=tokens["access_token"],
            refresh_token=tokens["refresh_token"],
            user=userData
        )

        return ServiceResponse(
            isSuccess=True,
            message="Login successful",
            statusCode=200,
            data=responseData.to_dict()
        )

    def getById(self, dto: GetCurrentUserDTO) -> ServiceResponse:
        admin = self.db.query(Admin).filter(
            Admin.id == dto.getUserId(),
            Admin.deleted_at.is_(None)
        ).first()

        if not admin:
            return ServiceResponse(
                isSuccess=False,
                message="User not found",
                statusCode=404,
                data=None
            )

        userData = {
            "id": admin.id,
            "email": admin.email,
            "username": admin.username,
            "first_name": admin.first_name,
            "last_name": admin.last_name,
            "user_type": "admin",
            "company_id": admin.company_id
        }

        return ServiceResponse(
            isSuccess=True,
            message="User retrieved successfully",
            statusCode=200,
            data=userData
        )

    def refreshToken(self, dto: RefreshTokenDTO) -> ServiceResponse:
        tokens = self.jwtService.refreshAccessToken(dto.getRefreshToken())
        if not tokens:
            return ServiceResponse(
                isSuccess=False,
                message="Invalid refresh token",
                statusCode=401,
                data=None
            )

        return ServiceResponse(
            isSuccess=True,
            message="Token refreshed successfully",
            statusCode=200,
            data=tokens
        )

    def logout(self, dto: LogoutDTO) -> ServiceResponse:
        self.tokenService.deactivateToken(dto.getAccessToken())
        return ServiceResponse(
            isSuccess=True,
            message="Logout successful",
            statusCode=200,
            data=None
        )

    def getJwtService(self) -> JWTService:
        return self.jwtService
