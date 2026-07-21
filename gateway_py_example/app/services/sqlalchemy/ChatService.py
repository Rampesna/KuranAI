from sqlalchemy.orm import Session
from sqlalchemy import func, desc
from datetime import datetime, timezone
from app.models.chat import ChatSession, ChatMessage
from app.types.chat_role import ChatRole
from app.utils.http_response import httpResponse
from typing import Optional, Dict, Any, List
from app.services.intelor_services.RAGService import RAGService
from app.services.openai.ChatService import ChatService as OpenAIChatService
from app.dto.services.requests.manager.chat.GetSessionsDTO import GetSessionsDTO
from app.dto.services.requests.manager.chat.CreateSessionDTO import CreateSessionDTO
from app.dto.services.requests.manager.chat.GetMessagesDTO import GetMessagesDTO
from app.dto.services.requests.manager.chat.SendMessageDTO import SendMessageDTO
from app.dto.services.requests.manager.chat.DeleteSessionDTO import DeleteSessionDTO


class ChatService:
    def __init__(self, db: Session, rag_service: Optional[RAGService] = None):
        self.db = db
        self.ragService = rag_service

    def getSessions(self, dto: GetSessionsDTO) -> Dict[str, Any]:
        # Calculate offset (Laravel style: skip = pageIndex * pageSize)
        offset = dto.getPageIndex() * dto.getPageSize()

        # Get sessions - order by updated_at DESC (newest first)
        sessions = self.db.query(ChatSession).filter(
            ChatSession.company_id == dto.getCompanyId(),
            ChatSession.user_type == dto.getUserType(),
            ChatSession.user_id == dto.getUserId(),
            ChatSession.deleted_at.is_(None)
        ).order_by(desc(ChatSession.updated_at)).offset(offset).limit(dto.getPageSize()).all()

        # Build session DTOs with message count and last message
        sessionItems = []
        for session in sessions:
            messageCount = self.db.query(ChatMessage).filter(
                ChatMessage.chat_session_id == session.id,
                ChatMessage.deleted_at.is_(None)
            ).count()

            lastMessage = self.db.query(ChatMessage).filter(
                ChatMessage.chat_session_id == session.id,
                ChatMessage.deleted_at.is_(None)
            ).order_by(desc(ChatMessage.created_at)).first()

            lastMessageContent = lastMessage.content if lastMessage else None

            sessionItems.append({
                "id": session.id,
                "title": session.title,
                "message_count": messageCount,
                "last_message": lastMessageContent,
                "created_at": session.created_at.isoformat() if session.created_at else None,
                "updated_at": session.updated_at.isoformat() if session.updated_at else None,
            })

        # Laravel style response - return data directly (not wrapped in items/meta)
        return httpResponse(
            isSuccess=True,
            message="Sessions retrieved successfully",
            statusCode=200,
            data=sessionItems
        )

    def createSession(self, dto: CreateSessionDTO) -> Dict[str, Any]:
        newSession = ChatSession(
            company_id=dto.getCompanyId(),
            user_type=dto.getUserType(),
            user_id=dto.getUserId(),
            title=dto.getTitle()
        )
        # Event listener will automatically set created_at and updated_at
        self.db.add(newSession)
        self.db.commit()
        self.db.refresh(newSession)

        try:
            from app.utils.websocket_events import notifySessionCreated
            notifySessionCreated(newSession.id)
        except Exception:
            pass

        return httpResponse(
            isSuccess=True,
            message="Session created successfully",
            statusCode=201,
            data={
                "id": newSession.id,
                "title": newSession.title,
                "created_at": newSession.created_at.isoformat() if newSession.created_at else None,
                "updated_at": newSession.updated_at.isoformat() if newSession.updated_at else None,
            }
        )

    def getMessages(self, dto: GetMessagesDTO) -> Dict[str, Any]:
        # Get total count
        totalCount = self.db.query(ChatMessage).filter(
            ChatMessage.chat_session_id == dto.getSessionId(),
            ChatMessage.deleted_at.is_(None)
        ).count()

        # For pagination: we want to show newest messages at bottom, oldest at top
        # pageIndex=0: Get the last pageSize messages (newest messages)
        # pageIndex=1: Get the previous pageSize messages (older messages)
        # We need to calculate offset from the end
        
        # Calculate how many messages to skip from the beginning
        # If totalCount = 20, pageSize = 10:
        #   pageIndex=0: skip 10, get last 10 (newest) - messages 11-20
        #   pageIndex=1: skip 0, get first 10 (oldest) - messages 1-10
        # If totalCount = 6, pageSize = 10:
        #   pageIndex=0: skip 0, get all 6 (newest) - messages 1-6
        #   pageIndex=1: skip 6, get 0 (no more) - empty array
        
        messagesToSkip = max(0, totalCount - (dto.getPageIndex() + 1) * dto.getPageSize())
        
        # If skip is >= totalCount, no messages to return
        if messagesToSkip >= totalCount:
            messages = []
        else:
            messagesToTake = min(dto.getPageSize(), totalCount - messagesToSkip)
            
            # Get messages ordered by created_at ASC (oldest first)
            # Then skip and take the appropriate range
            messages = self.db.query(ChatMessage).filter(
                ChatMessage.chat_session_id == dto.getSessionId(),
                ChatMessage.deleted_at.is_(None)
            ).order_by(ChatMessage.created_at).offset(messagesToSkip).limit(messagesToTake).all()

        messageItems = []
        for message in messages:
            messageItems.append({
                "id": message.id,
                "chat_session_id": message.chat_session_id,
                "role": message.role,
                "content": message.content,
                "token_count": message.token_count,
                "generated_sql": message.generated_sql,
                "sql_result": message.sql_result,
                "created_at": message.created_at.isoformat() if message.created_at else None,
                "updated_at": message.updated_at.isoformat() if message.updated_at else None,
            })

        # Laravel style response - return data directly (not wrapped in items/meta)
        return httpResponse(
            isSuccess=True,
            message="Messages retrieved successfully",
            statusCode=200,
            data=messageItems
        )

    def sendMessage(self, dto: SendMessageDTO) -> Dict[str, Any]:
        # Check if session exists
        session = self.db.query(ChatSession).filter(
            ChatSession.id == dto.getSessionId(),
            ChatSession.deleted_at.is_(None)
        ).first()

        if not session:
            return httpResponse(
                isSuccess=False,
                message="Chat session not found",
                statusCode=404,
                data=None
            )

        # Create user message
        userMessage = ChatMessage(
            chat_session_id=dto.getSessionId(),
            role=ChatRole.USER,
            content=dto.getContent()
        )
        self.db.add(userMessage)
        self.db.commit()

        # Get chat history (last 20 messages for context)
        previousMessages = self.db.query(ChatMessage).filter(
            ChatMessage.chat_session_id == dto.getSessionId(),
            ChatMessage.deleted_at.is_(None)
        ).order_by(ChatMessage.created_at).limit(20).all()

        # Build messages array for OpenAI
        messages = []
        for msg in previousMessages:
            messages.append({
                "role": msg.role,
                "content": msg.content
            })
        
        # Add current user message
        messages.append({
            "role": "user",
            "content": dto.getContent()
        })

        # Call OpenAI Chat API
        chatService = OpenAIChatService()
        
        try:
            assistantContent = chatService.chat(messages, temperature=0.7)
        except Exception as e:
            assistantContent = f"Bir hata oluştu: {str(e)}"

        # Create assistant message
        assistantMessage = ChatMessage(
            chat_session_id=dto.getSessionId(),
            role=ChatRole.ASSISTANT,
            content=assistantContent
        )
        self.db.add(assistantMessage)

        # Update session title if not set
        if not session.title:
            # Truncate content to 50 chars for title
            newTitle = dto.getContent()[:50] + "..." if len(dto.getContent()) > 50 else dto.getContent()
            session.title = newTitle

        # Event listener will automatically update updated_at on commit
        self.db.commit()
        self.db.refresh(assistantMessage)

        return httpResponse(
            isSuccess=True,
            message="Message processed successfully",
            statusCode=200,
            data={
                "id": assistantMessage.id,
                "chat_session_id": assistantMessage.chat_session_id,
                "role": assistantMessage.role,
                "content": assistantMessage.content,
                "token_count": assistantMessage.token_count,
                "generated_sql": assistantMessage.generated_sql,
                "sql_result": assistantMessage.sql_result,
                "created_at": assistantMessage.created_at.isoformat() if assistantMessage.created_at else None,
                "updated_at": assistantMessage.updated_at.isoformat() if assistantMessage.updated_at else None,
            }
        )

    def getChatHistory(self, sessionId: int, limit: int = 20) -> List[ChatMessage]:
        """Get chat history for context"""
        return self.db.query(ChatMessage).filter(
            ChatMessage.chat_session_id == sessionId,
            ChatMessage.deleted_at.is_(None)
        ).order_by(ChatMessage.created_at).limit(limit).all()

    def getSessionById(self, sessionId: int) -> Optional[ChatSession]:
        return self.db.query(ChatSession).filter(
            ChatSession.id == sessionId,
            ChatSession.deleted_at.is_(None)
        ).first()

    def getSession(
        self, sessionId: int, companyId: int, userType: str, userId: int
    ) -> Optional[Dict[str, Any]]:
        """Return one session as list-item DTO (same shape as getSessions items) or None."""
        session = self.db.query(ChatSession).filter(
            ChatSession.id == sessionId,
            ChatSession.company_id == companyId,
            ChatSession.user_type == userType,
            ChatSession.user_id == userId,
            ChatSession.deleted_at.is_(None),
        ).first()
        if not session:
            return None
        messageCount = self.db.query(ChatMessage).filter(
            ChatMessage.chat_session_id == session.id,
            ChatMessage.deleted_at.is_(None),
        ).count()
        lastMessage = (
            self.db.query(ChatMessage)
            .filter(
                ChatMessage.chat_session_id == session.id,
                ChatMessage.deleted_at.is_(None),
            )
            .order_by(desc(ChatMessage.created_at))
            .first()
        )
        lastMessageContent = lastMessage.content if lastMessage else None
        return {
            "id": session.id,
            "title": session.title,
            "message_count": messageCount,
            "last_message": lastMessageContent,
            "created_at": session.created_at.isoformat() if session.created_at else None,
            "updated_at": session.updated_at.isoformat() if session.updated_at else None,
        }

    def deleteSession(self, dto: DeleteSessionDTO) -> Dict[str, Any]:
        session = self.db.query(ChatSession).filter(
            ChatSession.id == dto.getSessionId(),
            ChatSession.deleted_at.is_(None)
        ).first()

        if not session:
            return httpResponse(
                isSuccess=False,
                message="Session not found",
                statusCode=404,
                data=None
            )

        sessionId = dto.getSessionId()
        # Soft delete
        session.deleted_at = datetime.now(timezone.utc)
        self.db.commit()

        try:
            from app.utils.websocket_events import notifySessionDeleted
            notifySessionDeleted(sessionId)
        except Exception:
            pass

        return httpResponse(
            isSuccess=True,
            message="Session deleted successfully",
            statusCode=200,
            data=None
        )
