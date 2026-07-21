from fastapi import Depends, HTTPException, status, Request
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.core.controller import Controller
from app.utils.http_response import httpResponse
from app.services.sqlalchemy.AdminService import AdminService
from app.services.common.jwt import JWTService
from app.http.requests.api.admin.authController.LoginRequest import LoginRequest
from app.http.requests.api.admin.authController.RefreshTokenRequest import RefreshTokenRequest
from app.dto.services.requests.admin.auth.LoginDTO import LoginDTO
from app.dto.services.requests.admin.auth.RefreshTokenDTO import RefreshTokenDTO
from app.dto.services.requests.admin.auth.GetCurrentUserDTO import GetCurrentUserDTO
from app.dto.services.requests.admin.auth.LogoutDTO import LogoutDTO
from app.http.requests.api.admin.authController.GetCurrentUserRequest import GetCurrentUserRequest
from app.http.requests.api.admin.authController.LogoutRequest import LogoutRequest


class AuthController(Controller):
    def __init__(self, adminService: AdminService):
        self.adminService = adminService

    def login(self, request: LoginRequest, httpRequest: Request, db: Session = Depends(get_db)):
        response = self.adminService.login(LoginDTO(request), httpRequest)
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def refreshToken(self, request: RefreshTokenRequest, db: Session = Depends(get_db)):
        response = self.adminService.refreshToken(RefreshTokenDTO(request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def me(self, request: GetCurrentUserRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.adminService.getById(GetCurrentUserDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def logout(self, request: LogoutRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.adminService.logout(LogoutDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )
