from fastapi import Depends
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.core.controller import Controller
from app.utils.http_response import httpResponse
from app.services.sqlalchemy.AdminCompanyService import AdminCompanyService


class CompanyController(Controller):
    def __init__(self, admin_company_service: AdminCompanyService):
        self.admin_company_service = admin_company_service

    def list_companies(self, current_user: dict, db: Session = Depends(get_db)):
        response = self.admin_company_service.list_companies(current_user)
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )

    def get_company(self, company_id: int, current_user: dict, db: Session = Depends(get_db)):
        response = self.admin_company_service.get_company(company_id, current_user)
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData(),
        )
