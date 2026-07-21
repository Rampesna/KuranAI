from fastapi import APIRouter, Depends, Query, Request
from fastapi import Request
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.services.sqlalchemy.AdminService import AdminService
from app.services.sqlalchemy.ChatService import ChatService as ChatServiceService
from app.services.common.jwt import JWTService
from app.services.openai.ChatService import ChatService
from app.services.intelor_services.RAGService import RAGService
from app.http.controllers.api.admin.AuthController import AuthController
from app.http.controllers.api.admin.ChatController import ChatController
from app.http.controllers.api.admin.ManagerController import ManagerController
from app.http.controllers.api.admin.CompanyController import CompanyController
from app.http.middlewares.Auth import getCurrentAdmin
from app.services.sqlalchemy.AdminManagerService import AdminManagerService
from app.services.sqlalchemy.AdminCompanyService import AdminCompanyService
from fastapi import HTTPException, status
from app.http.requests.api.admin.managerController.ListManagersRequest import ListManagersRequest
from app.http.requests.api.admin.managerController.CreateManagerRequest import CreateManagerRequest
from app.http.requests.api.admin.managerController.UpdateManagerRequest import UpdateManagerRequest
from app.http.requests.api.admin.managerController.GetManagerRequest import GetManagerRequest
from app.http.requests.api.admin.managerController.DeleteManagerRequest import DeleteManagerRequest
from app.http.requests.api.admin.authController.LoginRequest import LoginRequest
from app.http.requests.api.admin.authController.RefreshTokenRequest import RefreshTokenRequest
from app.http.requests.api.admin.authController.GetCurrentUserRequest import GetCurrentUserRequest
from app.http.requests.api.admin.authController.LogoutRequest import LogoutRequest
from app.http.requests.api.manager.chatController.CreateSessionRequest import CreateSessionRequest
from app.http.requests.api.manager.chatController.SendMessageRequest import SendMessageRequest
from app.http.requests.api.manager.chatController.GetSessionsRequest import GetSessionsRequest
from app.http.requests.api.manager.chatController.GetMessagesRequest import GetMessagesRequest
from app.http.requests.api.manager.chatController.DeleteSessionRequest import DeleteSessionRequest

router = APIRouter(prefix="/api/admin", tags=["Admin"])


def getAdminService(db: Session = Depends(get_db)) -> AdminService:
    jwtService = JWTService()
    return AdminService(db, jwtService)


def getAuthController(db: Session = Depends(get_db)) -> AuthController:
    adminService = getAdminService(db)
    return AuthController(adminService)


def getChatController(db: Session = Depends(get_db)) -> ChatController:
    chatServiceService = ChatServiceService(db)
    # Initialize RAG service for chat
    openaiChatService = ChatService()
    ragService = RAGService(openaiChatService)
    chatServiceService.ragService = ragService
    return ChatController(chatServiceService, ragService)


def getAdminManagerService(db: Session = Depends(get_db)) -> AdminManagerService:
    return AdminManagerService(db)


def getManagerController(db: Session = Depends(get_db)) -> ManagerController:
    return ManagerController(getAdminManagerService(db))


def getAdminCompanyService(db: Session = Depends(get_db)) -> AdminCompanyService:
    return AdminCompanyService(db)


def getCompanyController(db: Session = Depends(get_db)) -> CompanyController:
    return CompanyController(getAdminCompanyService(db))


def _ensure_company_access(company_id: int, current_user: dict) -> dict:
    """If admin can access the company, return current_user with company_id set to path company_id. Else raise 403."""
    is_super = current_user.get("is_super_admin")
    if is_super or current_user.get("company_id") == company_id:
        return {**current_user, "company_id": company_id}
    raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Access denied to this company")


@router.post("/auth/login", summary="Admin Login", description="Login endpoint for admin users")
def login(request: LoginRequest, httpRequest: Request, db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.login(request, httpRequest, db)


@router.post("/auth/refresh", summary="Refresh Token", description="Refresh access token using refresh token")
def refreshToken(request: RefreshTokenRequest, db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.refreshToken(request, db)


@router.get("/auth/me", summary="Get Current User", description="Get current authenticated admin user")
def me(request: GetCurrentUserRequest = Depends(), currentUser: dict = Depends(getCurrentAdmin), db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.me(request, currentUser, db)


@router.post("/auth/logout", summary="Logout", description="Logout current admin user")
def logout(request: LogoutRequest, currentUser: dict = Depends(getCurrentAdmin), db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.logout(request, currentUser, db)


# Chat routes
@router.get("/chat/sessions", summary="Get Chat Sessions", description="Get paginated chat sessions")
def getSessions(
    page: int = Query(1, ge=1, alias="page"),
    perPage: int = Query(20, ge=1, le=100, alias="per_page"),
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    from app.http.requests.api.admin.chatController.GetSessionsRequest import GetSessionsRequest
    request = GetSessionsRequest(page=page, perPage=perPage)
    return controller.getSessions(request, currentUser, db)


@router.get("/chat/sessions/{session_id}", summary="Get Chat Session", description="Get a single chat session by id")
def getSession(
    session_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getChatController(db)
    return controller.getSession(session_id, currentUser)


@router.post("/chat/sessions", summary="Create Chat Session", description="Create a new chat session")
def createSession(
    request: CreateSessionRequest,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    return controller.createSession(request, currentUser, db)


@router.get("/chat/sessions/{session_id}/messages", summary="Get Chat Messages", description="Get paginated messages for a session")
def getMessages(
    session_id: int,
    page: int = Query(1, ge=1, alias="page"),
    perPage: int = Query(50, ge=1, le=100, alias="per_page"),
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = GetMessagesRequest(**{"session_id": session_id, "page": page, "per_page": perPage})
    return controller.getMessages(request, currentUser, db)


@router.post("/chat/sessions/{session_id}/messages", summary="Send Message", description="Send a message in a chat session")
def sendMessage(
    session_id: int,
    request: SendMessageRequest,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request.sessionId = session_id
    return controller.sendMessage(request, currentUser, db)


@router.delete("/chat/sessions/{session_id}", summary="Delete Chat Session", description="Delete a chat session")
def deleteSession(
    session_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = DeleteSessionRequest(**{"session_id": session_id})
    return controller.deleteSession(request, currentUser, db)


# Manager CRUD (admin manages managers of their company)
@router.get("/managers", summary="List Managers", description="List managers with pagination (pageIndex, pageSize)")
def listManagers(
    pageIndex: int = Query(0, ge=0, alias="pageIndex"),
    pageSize: int = Query(20, ge=1, le=100, alias="pageSize"),
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getManagerController(db)
    request = ListManagersRequest(pageIndex=pageIndex, pageSize=pageSize)
    return controller.listManagers(request, currentUser, db)


@router.get("/managers/{manager_id}", summary="Get Manager", description="Get a manager by id")
def getManager(
    manager_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getManagerController(db)
    request = GetManagerRequest(managerId=manager_id)
    return controller.getManager(request, currentUser, db)


@router.post("/managers", summary="Create Manager", description="Create a new manager")
def createManager(
    request: CreateManagerRequest,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getManagerController(db)
    return controller.createManager(request, currentUser, db)


@router.put("/managers/{manager_id}", summary="Update Manager", description="Update a manager")
def updateManager(
    manager_id: int,
    request: UpdateManagerRequest,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getManagerController(db)
    return controller.updateManager(manager_id, request, currentUser, db)


@router.delete("/managers/{manager_id}", summary="Delete Manager", description="Soft delete a manager")
def deleteManager(
    manager_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getManagerController(db)
    request = DeleteManagerRequest(managerId=manager_id)
    return controller.deleteManager(request, currentUser, db)


# Companies (list / detail; super_admin sees all, else own company only)
@router.get("/companies", summary="List Companies", description="List companies (super_admin: all, else: own company)")
def listCompanies(
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getCompanyController(db)
    return controller.list_companies(currentUser, db)


@router.get("/companies/{company_id}", summary="Get Company", description="Get company by id (access checked)")
def getCompany(
    company_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    controller = getCompanyController(db)
    return controller.get_company(company_id, currentUser, db)


# Managers under a specific company (company_id from path; access checked)
@router.get("/companies/{company_id}/managers", summary="List Managers by Company")
def listManagersByCompany(
    company_id: int,
    pageIndex: int = Query(0, ge=0, alias="pageIndex"),
    pageSize: int = Query(20, ge=1, le=100, alias="pageSize"),
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    user = _ensure_company_access(company_id, currentUser)
    controller = getManagerController(db)
    request = ListManagersRequest(pageIndex=pageIndex, pageSize=pageSize)
    return controller.listManagers(request, user, db)


@router.get("/companies/{company_id}/managers/{manager_id}", summary="Get Manager by Company")
def getManagerByCompany(
    company_id: int,
    manager_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    user = _ensure_company_access(company_id, currentUser)
    controller = getManagerController(db)
    request = GetManagerRequest(managerId=manager_id)
    return controller.getManager(request, user, db)


@router.post("/companies/{company_id}/managers", summary="Create Manager for Company")
def createManagerForCompany(
    company_id: int,
    request: CreateManagerRequest,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    user = _ensure_company_access(company_id, currentUser)
    controller = getManagerController(db)
    return controller.createManager(request, user, db)


@router.put("/companies/{company_id}/managers/{manager_id}", summary="Update Manager for Company")
def updateManagerForCompany(
    company_id: int,
    manager_id: int,
    request: UpdateManagerRequest,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    user = _ensure_company_access(company_id, currentUser)
    controller = getManagerController(db)
    return controller.updateManager(manager_id, request, user, db)


@router.delete("/companies/{company_id}/managers/{manager_id}", summary="Delete Manager for Company")
def deleteManagerForCompany(
    company_id: int,
    manager_id: int,
    currentUser: dict = Depends(getCurrentAdmin),
    db: Session = Depends(get_db),
):
    user = _ensure_company_access(company_id, currentUser)
    controller = getManagerController(db)
    request = DeleteManagerRequest(managerId=manager_id)
    return controller.deleteManager(request, user, db)
