from fastapi import Depends
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.core.controller import Controller
from app.utils.http_response import httpResponse
from app.services.sqlalchemy.AdminManagerService import AdminManagerService
from app.http.requests.api.admin.managerController.ListManagersRequest import ListManagersRequest
from app.http.requests.api.admin.managerController.CreateManagerRequest import CreateManagerRequest
from app.http.requests.api.admin.managerController.UpdateManagerRequest import UpdateManagerRequest
from app.http.requests.api.admin.managerController.GetManagerRequest import GetManagerRequest
from app.http.requests.api.admin.managerController.DeleteManagerRequest import DeleteManagerRequest
from app.dto.services.requests.admin.manager.ListManagersDTO import ListManagersDTO
from app.dto.services.requests.admin.manager.CreateManagerDTO import CreateManagerDTO
from app.dto.services.requests.admin.manager.UpdateManagerDTO import UpdateManagerDTO
from app.dto.services.requests.admin.manager.GetManagerDTO import GetManagerDTO
from app.dto.services.requests.admin.manager.DeleteManagerDTO import DeleteManagerDTO


class ManagerController(Controller):
    def __init__(self, adminManagerService: AdminManagerService):
        self.adminManagerService = adminManagerService

    def listManagers(self, request: ListManagersRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.adminManagerService.listManagers(ListManagersDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )

    def getManager(self, request: GetManagerRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.adminManagerService.getManager(GetManagerDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )

    def createManager(self, request: CreateManagerRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.adminManagerService.createManager(CreateManagerDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )

    def updateManager(
        self, managerId: int, request: UpdateManagerRequest, currentUser: dict, db: Session = Depends(get_db)
    ):
        response = self.adminManagerService.updateManager(
            UpdateManagerDTO(currentUser, request, managerId)
        )
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )

    def deleteManager(self, request: DeleteManagerRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.adminManagerService.deleteManager(DeleteManagerDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )
