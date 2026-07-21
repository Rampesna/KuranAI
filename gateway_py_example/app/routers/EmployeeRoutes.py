from fastapi import APIRouter, Depends, Query, Request
from fastapi import Request
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.services.sqlalchemy.EmployeeService import EmployeeService
from app.services.sqlalchemy.ChatService import ChatService as ChatServiceService
from app.services.common.jwt import JWTService
from app.services.openai.ChatService import ChatService
from app.services.intelor_services.RAGService import RAGService
from app.http.controllers.api.employee.AuthController import AuthController
from app.http.controllers.api.employee.ChatController import ChatController
from app.http.middlewares.Auth import getCurrentEmployee
from app.http.requests.api.employee.authController.LoginRequest import LoginRequest
from app.http.requests.api.employee.authController.RefreshTokenRequest import RefreshTokenRequest
from app.http.requests.api.employee.authController.GetCurrentUserRequest import GetCurrentUserRequest
from app.http.requests.api.employee.authController.LogoutRequest import LogoutRequest
from app.http.requests.api.manager.chatController.CreateSessionRequest import CreateSessionRequest
from app.http.requests.api.employee.chatController.GetSessionsRequest import GetSessionsRequest
from app.http.requests.api.employee.chatController.GetMessagesRequest import GetMessagesRequest
from app.http.requests.api.employee.chatController.SendMessageRequest import SendMessageRequest
from app.http.requests.api.employee.chatController.DeleteSessionRequest import DeleteSessionRequest

router = APIRouter(prefix="/api/employee", tags=["Employee"])


def getEmployeeService(db: Session = Depends(get_db)) -> EmployeeService:
    jwtService = JWTService()
    return EmployeeService(db, jwtService)


def getAuthController(db: Session = Depends(get_db)) -> AuthController:
    employeeService = getEmployeeService(db)
    return AuthController(employeeService)


def getChatController(db: Session = Depends(get_db)) -> ChatController:
    chatServiceService = ChatServiceService(db)
    # Initialize RAG service for chat
    openaiChatService = ChatService()
    ragService = RAGService(openaiChatService)
    chatServiceService.ragService = ragService
    return ChatController(chatServiceService, ragService)


@router.post("/auth/login", summary="Employee Login", description="Login endpoint for employee users")
def login(request: LoginRequest, httpRequest: Request, db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.login(request, httpRequest, db)


@router.post("/auth/refresh", summary="Refresh Token", description="Refresh access token using refresh token")
def refreshToken(request: RefreshTokenRequest, db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.refreshToken(request, db)


@router.get("/auth/me", summary="Get Current User", description="Get current authenticated employee user")
def me(request: GetCurrentUserRequest = Depends(), currentUser: dict = Depends(getCurrentEmployee), db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.me(request, currentUser, db)


@router.post("/auth/logout", summary="Logout", description="Logout current employee user")
def logout(request: LogoutRequest, currentUser: dict = Depends(getCurrentEmployee), db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.logout(request, currentUser, db)


# Chat routes
@router.get("/chat/sessions", summary="Get Chat Sessions", description="Get paginated chat sessions")
def getSessions(
    page: int = Query(1, ge=1, alias="page"),
    perPage: int = Query(20, ge=1, le=100, alias="per_page"),
    currentUser: dict = Depends(getCurrentEmployee),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    from app.http.requests.api.employee.chatController.GetSessionsRequest import GetSessionsRequest
    request = GetSessionsRequest(page=page, perPage=perPage)
    return controller.getSessions(request, currentUser, db)


@router.get("/chat/sessions/{session_id}", summary="Get Chat Session", description="Get a single chat session by id")
def getSession(
    session_id: int,
    currentUser: dict = Depends(getCurrentEmployee),
    db: Session = Depends(get_db),
):
    controller = getChatController(db)
    return controller.getSession(session_id, currentUser)


@router.post("/chat/sessions", summary="Create Chat Session", description="Create a new chat session")
def createSession(
    request: CreateSessionRequest,
    currentUser: dict = Depends(getCurrentEmployee),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    return controller.createSession(request, currentUser, db)


@router.get("/chat/sessions/{session_id}/messages", summary="Get Chat Messages", description="Get paginated messages for a session")
def getMessages(
    session_id: int,
    page: int = Query(1, ge=1, alias="page"),
    perPage: int = Query(50, ge=1, le=100, alias="per_page"),
    currentUser: dict = Depends(getCurrentEmployee),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = GetMessagesRequest(**{"session_id": session_id, "page": page, "per_page": perPage})
    return controller.getMessages(request, currentUser, db)


@router.post("/chat/sessions/{session_id}/messages", summary="Send Message", description="Send a message in a chat session")
def sendMessage(
    session_id: int,
    request: SendMessageRequest,
    currentUser: dict = Depends(getCurrentEmployee),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request.sessionId = session_id
    return controller.sendMessage(request, currentUser, db)


@router.delete("/chat/sessions/{session_id}", summary="Delete Chat Session", description="Delete a chat session")
def deleteSession(
    session_id: int,
    currentUser: dict = Depends(getCurrentEmployee),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = DeleteSessionRequest(**{"session_id": session_id})
    return controller.deleteSession(request, currentUser, db)
