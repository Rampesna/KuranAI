from sqlalchemy.orm import Session
from typing import Dict, Any, List
from app.models.company import Company
from app.core.ServiceResponse import ServiceResponse


def _company_to_safe_dict(c: Company) -> Dict[str, Any]:
    """Return company fields for API; exclude sensitive DB credentials from response."""
    d = {
        "id": c.id,
        "name": c.name,
        "slug": c.slug,
        "description": c.description,
        "logo_url": c.logo_url,
        "website": c.website,
        "email": c.email,
        "phone": c.phone,
        "address": c.address,
        "is_active": c.is_active,
        "database_type": c.database_type,
        "created_at": c.created_at.isoformat() if c.created_at else None,
        "updated_at": c.updated_at.isoformat() if c.updated_at else None,
    }
    return d


class AdminCompanyService:
    def __init__(self, db: Session):
        self.db = db

    def _can_access_company(self, current_user: dict, company_id: int) -> bool:
        if current_user.get("is_super_admin"):
            return True
        return current_user.get("company_id") == company_id

    def list_companies(self, current_user: dict) -> ServiceResponse:
        try:
            query = self.db.query(Company).filter(Company.deleted_at.is_(None))
            if not current_user.get("is_super_admin"):
                query = query.filter(Company.id == current_user.get("company_id"))
            companies = query.order_by(Company.name).all()
            items = [_company_to_safe_dict(c) for c in companies]
            return ServiceResponse(
                isSuccess=True,
                message="Companies retrieved successfully",
                statusCode=200,
                data={"items": items},
            )
        except Exception as e:
            return ServiceResponse(
                isSuccess=False,
                message=f"Failed to list companies: {str(e)}",
                statusCode=500,
                data=None,
            )

    def get_company(self, company_id: int, current_user: dict) -> ServiceResponse:
        if not self._can_access_company(current_user, company_id):
            return ServiceResponse(
                isSuccess=False,
                message="Access denied to this company",
                statusCode=403,
                data=None,
            )
        company = (
            self.db.query(Company)
            .filter(Company.id == company_id, Company.deleted_at.is_(None))
            .first()
        )
        if not company:
            return ServiceResponse(
                isSuccess=False,
                message="Company not found",
                statusCode=404,
                data=None,
            )
        return ServiceResponse(
            isSuccess=True,
            message="Company retrieved successfully",
            statusCode=200,
            data=_company_to_safe_dict(company),
        )
