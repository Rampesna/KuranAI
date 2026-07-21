from fastapi import Depends
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.core.controller import Controller
from app.services.sqlalchemy.ChatService import ChatService
from app.services.intelor_services.RAGService import RAGService
from app.services.openai.ChatService import ChatService as OpenAIChatService
from app.http.middlewares.Auth import getCurrentEmployee
from app.utils.http_response import httpResponse
from app.http.requests.api.employee.chatController.GetSessionsRequest import GetSessionsRequest
from app.http.requests.api.manager.chatController.CreateSessionRequest import CreateSessionRequest
from app.http.requests.api.employee.chatController.GetMessagesRequest import GetMessagesRequest
from app.http.requests.api.employee.chatController.SendMessageRequest import SendMessageRequest
from app.http.requests.api.employee.chatController.DeleteSessionRequest import DeleteSessionRequest
from app.dto.services.requests.employee.chat.GetSessionsDTO import GetSessionsDTO
from app.dto.services.requests.employee.chat.CreateSessionDTO import CreateSessionDTO
from app.dto.services.requests.employee.chat.GetMessagesDTO import GetMessagesDTO
from app.dto.services.requests.employee.chat.SendMessageDTO import SendMessageDTO
from app.dto.services.requests.employee.chat.DeleteSessionDTO import DeleteSessionDTO


class ChatController(Controller):
    def __init__(self, chatService: ChatService, ragService: RAGService = None):
        self.chatService = chatService
        self.ragService = ragService

    def getSessions(
        self,
        request: GetSessionsRequest,
        currentUser: dict,
        db: Session = Depends(get_db)
    ):
        return self.chatService.getSessions(GetSessionsDTO(currentUser, request))

    def getSession(self, sessionId: int, currentUser: dict):
        companyId = currentUser.get("company_id")
        userType = currentUser.get("user_type", "employee")
        userId = currentUser.get("user_id") or currentUser.get("id")
        item = self.chatService.getSession(sessionId, companyId, userType, userId)
        if not item:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None,
            )
        return httpResponse(
            isSuccess=True,
            message="Session retrieved successfully",
            statusCode=200,
            data=item,
        )

    def createSession(
        self,
        request: CreateSessionRequest,
        currentUser: dict,
        db: Session = Depends(get_db)
    ):
        return self.chatService.createSession(CreateSessionDTO(currentUser, request))

    def getMessages(
        self,
        request: GetMessagesRequest,
        currentUser: dict,
        db: Session = Depends(get_db)
    ):
        # Verify ownership
        session = self.chatService.getSessionById(request.sessionId)
        if not session:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None
            )

        companyId = currentUser.get("company_id")
        if session.company_id != companyId:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None
            )

        return self.chatService.getMessages(GetMessagesDTO(currentUser, request))

    def sendMessage(
        self,
        request: SendMessageRequest,
        currentUser: dict,
        db: Session = Depends(get_db)
    ):
        # Verify ownership
        session = self.chatService.getSessionById(request.sessionId)
        if not session:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None
            )

        companyId = currentUser.get("company_id")
        if session.company_id != companyId:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None
            )

        return self.chatService.sendMessage(SendMessageDTO(currentUser, request))

    def deleteSession(
        self,
        request: DeleteSessionRequest,
        currentUser: dict,
        db: Session = Depends(get_db)
    ):
        # Verify ownership
        session = self.chatService.getSessionById(request.sessionId)
        if not session:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None
            )

        companyId = currentUser.get("company_id")
        if session.company_id != companyId:
            return httpResponse(
                isSuccess=False,
                message="Session not found or access denied",
                statusCode=404,
                data=None
            )

        return self.chatService.deleteSession(DeleteSessionDTO(currentUser, request))
