from sqlalchemy.orm import Session
from sqlalchemy import desc
from datetime import datetime, timezone
from typing import Optional, Dict, Any, List
from app.models.manager import Manager
from app.services.common.password import PasswordService
from app.core.ServiceResponse import ServiceResponse
from app.dto.services.requests.admin.manager.ListManagersDTO import ListManagersDTO
from app.dto.services.requests.admin.manager.CreateManagerDTO import CreateManagerDTO
from app.dto.services.requests.admin.manager.UpdateManagerDTO import UpdateManagerDTO
from app.dto.services.requests.admin.manager.GetManagerDTO import GetManagerDTO
from app.dto.services.requests.admin.manager.DeleteManagerDTO import DeleteManagerDTO


def _managerToDict(m: Manager) -> Dict[str, Any]:
    return {
        "id": m.id,
        "company_id": m.company_id,
        "email": m.email,
        "username": m.username,
        "first_name": m.first_name,
        "last_name": m.last_name,
        "is_active": m.is_active,
        "last_login_at": m.last_login_at.isoformat() if m.last_login_at else None,
        "created_at": m.created_at.isoformat() if m.created_at else None,
        "updated_at": m.updated_at.isoformat() if m.updated_at else None,
    }


class AdminManagerService:
    def __init__(self, db: Session):
        self.db = db
        self.passwordService = PasswordService()

    def listManagers(self, dto: ListManagersDTO) -> ServiceResponse:
        try:
            offset = dto.getPageIndex() * dto.getPageSize()
            query = self.db.query(Manager).filter(
                Manager.company_id == dto.getCompanyId(),
                Manager.deleted_at.is_(None),
            )
            total = query.count()
            managers = (
                query.order_by(desc(Manager.updated_at))
                .offset(offset)
                .limit(dto.getPageSize())
                .all()
            )
            items = [_managerToDict(m) for m in managers]
            totalPages = (total + dto.getPageSize() - 1) // dto.getPageSize() if total > 0 else 0
            data = {
                "items": items,
                "meta": {
                    "pageIndex": dto.getPageIndex(),
                    "pageSize": dto.getPageSize(),
                    "total": total,
                    "totalPages": totalPages,
                },
            }
            return ServiceResponse(
                isSuccess=True,
                message="Managers retrieved successfully",
                statusCode=200,
                data=data,
            )
        except Exception as e:
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to list managers: {str(e)}",
                statusCode=500,
                data=None,
            )

    def getManager(self, dto: GetManagerDTO) -> ServiceResponse:
        manager = (
            self.db.query(Manager)
            .filter(
                Manager.id == dto.getManagerId(),
                Manager.company_id == dto.getCompanyId(),
                Manager.deleted_at.is_(None),
            )
            .first()
        )
        if not manager:
            return ServiceResponse(
                isSuccess=False,
                message="Manager not found",
                statusCode=404,
                data=None,
            )
        return ServiceResponse(
            isSuccess=True,
            message="Manager retrieved successfully",
            statusCode=200,
            data=_managerToDict(manager),
        )

    def createManager(self, dto: CreateManagerDTO) -> ServiceResponse:
        if self.db.query(Manager).filter(Manager.email == dto.getEmail(), Manager.deleted_at.is_(None)).first():
            return ServiceResponse(
                isSuccess=False,
                message="A manager with this email already exists",
                statusCode=422,
                data=None,
            )
        if self.db.query(Manager).filter(Manager.username == dto.getUsername(), Manager.deleted_at.is_(None)).first():
            return ServiceResponse(
                isSuccess=False,
                message="A manager with this username already exists",
                statusCode=422,
                data=None,
            )
        passwordHash = self.passwordService.hash(dto.getPassword())
        manager = Manager(
            company_id=dto.getCompanyId(),
            email=dto.getEmail(),
            username=dto.getUsername(),
            password_hash=passwordHash,
            first_name=dto.getFirstName(),
            last_name=dto.getLastName(),
            is_active=dto.getIsActive(),
        )
        self.db.add(manager)
        self.db.commit()
        self.db.refresh(manager)
        return ServiceResponse(
            isSuccess=True,
            message="Manager created successfully",
            statusCode=201,
            data=_managerToDict(manager),
        )

    def updateManager(self, dto: UpdateManagerDTO) -> ServiceResponse:
        manager = (
            self.db.query(Manager)
            .filter(
                Manager.id == dto.getManagerId(),
                Manager.company_id == dto.getCompanyId(),
                Manager.deleted_at.is_(None),
            )
            .first()
        )
        if not manager:
            return ServiceResponse(
                isSuccess=False,
                message="Manager not found",
                statusCode=404,
                data=None,
            )
        if dto.getFirstName() is not None:
            manager.first_name = dto.getFirstName()
        if dto.getLastName() is not None:
            manager.last_name = dto.getLastName()
        if dto.getIsActive() is not None:
            manager.is_active = dto.getIsActive()
        if dto.getPassword() is not None and dto.getPassword():
            manager.password_hash = self.passwordService.hash(dto.getPassword())
        manager.updated_at = datetime.now(timezone.utc)
        self.db.commit()
        self.db.refresh(manager)
        return ServiceResponse(
            isSuccess=True,
            message="Manager updated successfully",
            statusCode=200,
            data=_managerToDict(manager),
        )

    def deleteManager(self, dto: DeleteManagerDTO) -> ServiceResponse:
        manager = (
            self.db.query(Manager)
            .filter(
                Manager.id == dto.getManagerId(),
                Manager.company_id == dto.getCompanyId(),
                Manager.deleted_at.is_(None),
            )
            .first()
        )
        if not manager:
            return ServiceResponse(
                isSuccess=False,
                message="Manager not found",
                statusCode=404,
                data=None,
            )
        manager.deleted_at = datetime.now(timezone.utc)
        self.db.commit()
        return ServiceResponse(
            isSuccess=True,
            message="Manager deleted successfully",
            statusCode=200,
            data=None,
        )
