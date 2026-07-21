from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.services.sqlalchemy.AdminService import AdminService
from app.services.sqlalchemy.ManagerService import ManagerService
from app.services.sqlalchemy.EmployeeService import EmployeeService
from app.services.common.jwt import JWTService
from app.services.common.TokenService import TokenService
from app.models.admin import Admin
from app.models.manager import Manager
from app.models.employee import Employee

security = HTTPBearer()


def getCurrentAdmin(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: Session = Depends(get_db)
) -> dict:
    token = credentials.credentials
    
    # Validate token using TokenService (checks Redis first, then DB)
    # This validates JWT signature, expiry, logout status, and updates last_used_at
    tokenService = TokenService(db)
    payload = tokenService.validateTokenAndUpdateUsage(token)
    
    if not payload or payload.get("token_type") != "access" or payload.get("user_type") != "admin":
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token is invalid, expired, or logged out"
        )
    
    # Get company_id and is_super_admin from database
    admin = db.query(Admin).filter(Admin.id == payload["user_id"]).first()
    company_id = admin.company_id if admin else None
    is_super_admin = getattr(admin, "is_super_admin", False) if admin else False

    return {
        "id": payload["user_id"],  # Changed from "user_id" to "id" for consistency
        "user_id": payload["user_id"],  # Keep both for backward compatibility
        "user_type": payload["user_type"],
        "email": payload["email"],
        "company_id": company_id,
        "is_super_admin": is_super_admin,
        "access_token": token  # Include token for logout
    }


def getCurrentManager(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: Session = Depends(get_db)
) -> dict:
    token = credentials.credentials
    
    # Validate token using TokenService (checks Redis first, then DB)
    # This validates JWT signature, expiry, logout status, and updates last_used_at
    tokenService = TokenService(db)
    payload = tokenService.validateTokenAndUpdateUsage(token)
    
    if not payload or payload.get("token_type") != "access" or payload.get("user_type") != "manager":
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token is invalid, expired, or logged out"
        )
    
    # Get company_id from database
    manager = db.query(Manager).filter(Manager.id == payload["user_id"]).first()
    company_id = manager.company_id if manager else None
    
    return {
        "id": payload["user_id"],  # Changed from "user_id" to "id" for consistency
        "user_id": payload["user_id"],  # Keep both for backward compatibility
        "user_type": payload["user_type"],
        "email": payload["email"],
        "company_id": company_id,
        "access_token": token  # Include token for logout
    }


def getCurrentEmployee(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: Session = Depends(get_db)
) -> dict:
    token = credentials.credentials
    
    # Validate token using TokenService (checks Redis first, then DB)
    # This validates JWT signature, expiry, logout status, and updates last_used_at
    tokenService = TokenService(db)
    payload = tokenService.validateTokenAndUpdateUsage(token)
    
    if not payload or payload.get("token_type") != "access" or payload.get("user_type") != "employee":
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token is invalid, expired, or logged out"
        )
    
    # Get company_id from database
    employee = db.query(Employee).filter(Employee.id == payload["user_id"]).first()
    company_id = employee.company_id if employee else None
    
    return {
        "id": payload["user_id"],  # Changed from "user_id" to "id" for consistency
        "user_id": payload["user_id"],  # Keep both for backward compatibility
        "user_type": payload["user_type"],
        "email": payload["email"],
        "company_id": company_id,
        "access_token": token  # Include token for logout
    }
