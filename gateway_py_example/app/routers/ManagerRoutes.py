from fastapi import APIRouter, Depends, Query, Request
from fastapi import Request
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.services.sqlalchemy.ManagerService import ManagerService
from app.services.sqlalchemy.ChatService import ChatService as ChatServiceService
from app.services.common.jwt import JWTService
from app.services.openai.EmbeddingService import EmbeddingService
from app.services.openai.ChatService import ChatService
from app.services.intelor_services.RAGService import RAGService
from app.http.controllers.api.manager.AuthController import AuthController
from app.http.controllers.api.manager.TestController import TestController
from app.http.controllers.api.manager.ChatController import ChatController
from app.http.controllers.api.manager.AgentController import AgentController
from app.http.middlewares.Auth import getCurrentManager
from app.http.requests.api.manager.authController.LoginRequest import LoginRequest
from app.http.requests.api.manager.authController.RefreshTokenRequest import RefreshTokenRequest
from app.http.requests.api.manager.authController.GetCurrentUserRequest import GetCurrentUserRequest
from app.http.requests.api.manager.authController.LogoutRequest import LogoutRequest
from app.http.requests.api.manager.authController.TestOpenaiRequest import TestOpenAIRequest
from app.http.requests.api.manager.chatController.CreateSessionRequest import CreateSessionRequest
from app.http.requests.api.manager.chatController.SendMessageRequest import SendMessageRequest
from app.http.requests.api.manager.chatController.GetSessionsRequest import GetSessionsRequest
from app.http.requests.api.manager.chatController.GetMessagesRequest import GetMessagesRequest
from app.http.requests.api.manager.chatController.DeleteSessionRequest import DeleteSessionRequest
from app.http.requests.api.manager.chatController.StreamMessageRequest import StreamMessageRequest
from app.http.requests.api.manager.agentController.CreateAgentRequest import CreateAgentRequest
from app.http.requests.api.manager.agentController.UpdateAgentRequest import UpdateAgentRequest
from app.http.requests.api.manager.agentController.ListAgentsRequest import ListAgentsRequest
from app.http.requests.api.manager.agentController.GetAgentRequest import GetAgentRequest
from app.http.requests.api.manager.agentController.DeleteAgentRequest import DeleteAgentRequest
from app.services.sqlalchemy.AgentService import AgentService

router = APIRouter(prefix="/api/manager", tags=["Manager"])


def getManagerService(db: Session = Depends(get_db)) -> ManagerService:
    jwtService = JWTService()
    return ManagerService(db, jwtService)


def getAuthController(db: Session = Depends(get_db)) -> AuthController:
    managerService = getManagerService(db)
    return AuthController(managerService)


def getTestController(db: Session = Depends(get_db)) -> TestController:
    embeddingService = EmbeddingService()
    chatService = ChatService()
    ragService = RAGService(chatService)
    return TestController(embeddingService, chatService, ragService, db)


def getChatController(db: Session = Depends(get_db)) -> ChatController:
    chatServiceService = ChatServiceService(db)
    # Initialize RAG service for chat
    openaiChatService = ChatService()
    ragService = RAGService(openaiChatService)
    chatServiceService.ragService = ragService
    return ChatController(chatServiceService, ragService)


def getAgentController(db: Session = Depends(get_db)) -> AgentController:
    agentService = AgentService(db)
    return AgentController(agentService)


@router.post("/auth/login", summary="Manager Login", description="Login endpoint for manager users")
def login(request: LoginRequest, httpRequest: Request, db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.login(request, httpRequest, db)


@router.post("/auth/refresh", summary="Refresh Token", description="Refresh access token using refresh token")
def refreshToken(request: RefreshTokenRequest, db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.refreshToken(request, db)


@router.get("/auth/me", summary="Get Current User", description="Get current authenticated manager user")
def me(request: GetCurrentUserRequest = Depends(), currentUser: dict = Depends(getCurrentManager), db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.me(request, currentUser, db)


@router.post("/auth/logout", summary="Logout", description="Logout current manager user")
def logout(request: LogoutRequest, currentUser: dict = Depends(getCurrentManager), db: Session = Depends(get_db)):
    controller = getAuthController(db)
    return controller.logout(request, currentUser, db)


@router.post("/test/openai", summary="Test OpenAI", description="Test OpenAI embedding and RAG services")
async def testOpenai(request: TestOpenAIRequest, currentUser: dict = Depends(getCurrentManager), db: Session = Depends(get_db)):
    controller = getTestController(db)
    return await controller.testOpenai(request, currentUser, db)


# Chat routes
@router.get("/chat/sessions", summary="Get Chat Sessions", description="Get paginated chat sessions")
def getSessions(
    pageIndex: int = Query(0, ge=0, alias="pageIndex"),
    pageSize: int = Query(50, ge=1, le=100, alias="pageSize"),
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = GetSessionsRequest(pageIndex=pageIndex, pageSize=pageSize)
    return controller.getSessions(currentUser, request, db)


@router.get("/chat/sessions/{session_id}", summary="Get Chat Session", description="Get a single chat session by id")
def getSession(
    session_id: int,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db),
):
    controller = getChatController(db)
    return controller.getSession(session_id, currentUser)


@router.post("/chat/sessions", summary="Create Chat Session", description="Create a new chat session")
def createSession(
    request: CreateSessionRequest,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    return controller.createSession(currentUser, request, db)


@router.get("/chat/sessions/{session_id}/messages", summary="Get Chat Messages", description="Get paginated messages for a session")
def getMessages(
    session_id: int,
    pageIndex: int = Query(0, ge=0, alias="pageIndex"),
    pageSize: int = Query(10, ge=1, le=100, alias="pageSize"),
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = GetMessagesRequest(**{"session_id": session_id, "pageIndex": pageIndex, "pageSize": pageSize})
    return controller.getMessages(request, currentUser, db)


@router.post("/chat/sessions/{session_id}/messages", summary="Send Message", description="Send a message in a chat session")
def sendMessage(
    session_id: int,
    request: SendMessageRequest,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request.sessionId = session_id
    return controller.sendMessage(request, currentUser, db)


@router.delete("/chat/sessions/{session_id}", summary="Delete Chat Session", description="Delete a chat session")
def deleteSession(
    session_id: int,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request = DeleteSessionRequest(**{"session_id": session_id})
    return controller.deleteSession(request, currentUser, db)


@router.post("/chat/sessions/{session_id}/stream", summary="Stream Message", description="Send a message and stream the response via SSE")
async def streamMessage(
    session_id: int,
    request: StreamMessageRequest,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getChatController(db)
    request.sessionId = session_id
    return await controller.streamMessage(request, currentUser, db)


# Agent routes
@router.get("/agents", summary="List Agents", description="Get paginated list of agents")
def listAgents(
    request: ListAgentsRequest = Depends(),
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getAgentController(db)
    return controller.listAgents(request, currentUser, db)


@router.get("/agents/{agent_id}", summary="Get Agent", description="Get a single agent by ID")
def getAgent(
    agent_id: int,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getAgentController(db)
    request = GetAgentRequest(**{"agent_id": agent_id})
    return controller.getAgent(request, currentUser, db)


@router.post("/agents", summary="Create Agent", description="Create a new agent")
def createAgent(
    request: CreateAgentRequest,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getAgentController(db)
    return controller.createAgent(request, currentUser, db)


@router.put("/agents/{agent_id}", summary="Update Agent", description="Update an existing agent")
def updateAgent(
    agent_id: int,
    request: UpdateAgentRequest,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getAgentController(db)
    # Set agentId from path parameter
    request.agentId = agent_id
    return controller.updateAgent(request, currentUser, db)


@router.delete("/agents/{agent_id}", summary="Delete Agent", description="Delete an agent")
def deleteAgent(
    agent_id: int,
    currentUser: dict = Depends(getCurrentManager),
    db: Session = Depends(get_db)
):
    controller = getAgentController(db)
    request = DeleteAgentRequest(**{"agent_id": agent_id})
    return controller.deleteAgent(request, currentUser, db)
