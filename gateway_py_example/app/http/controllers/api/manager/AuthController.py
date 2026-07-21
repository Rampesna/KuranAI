from fastapi import Depends, HTTPException, status, Request
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.core.controller import Controller
from app.utils.http_response import httpResponse
from app.services.sqlalchemy.ManagerService import ManagerService
from app.services.common.TokenService import TokenService
from app.http.requests.api.manager.authController.LoginRequest import LoginRequest
from app.http.requests.api.manager.authController.RefreshTokenRequest import RefreshTokenRequest
from app.dto.services.requests.manager.auth.LoginDTO import LoginDTO
from app.dto.services.requests.manager.auth.RefreshTokenDTO import RefreshTokenDTO
from app.dto.services.requests.manager.auth.GetCurrentUserDTO import GetCurrentUserDTO
from app.dto.services.requests.manager.auth.LogoutDTO import LogoutDTO
from app.http.requests.api.manager.authController.GetCurrentUserRequest import GetCurrentUserRequest
from app.http.requests.api.manager.authController.LogoutRequest import LogoutRequest


class AuthController(Controller):
    def __init__(self, managerService: ManagerService):
        self.managerService = managerService

    def login(self, request: LoginRequest, httpRequest: Request, db: Session = Depends(get_db)):
        response = self.managerService.login(LoginDTO(request), httpRequest)
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def refreshToken(self, request: RefreshTokenRequest, db: Session = Depends(get_db)):
        response = self.managerService.refreshToken(RefreshTokenDTO(request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def me(self, request: GetCurrentUserRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.managerService.getById(GetCurrentUserDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def logout(self, request: LogoutRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.managerService.logout(LogoutDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )
