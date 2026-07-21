from fastapi import Depends
from fastapi.responses import StreamingResponse
from sqlalchemy.orm import Session
from app.core.database import get_db, SessionLocal
from app.core.controller import Controller
from app.services.sqlalchemy.ChatService import ChatService
from app.services.intelor_services.RAGService import RAGService
from app.services.openai.ChatService import ChatService as OpenAIChatService
from app.services.openai.EmbeddingService import EmbeddingService
from app.http.controllers.api.manager.TestController import TestController
from app.http.middlewares.Auth import getCurrentManager
from app.models.chat import ChatMessage, ChatSession
from app.types.chat_role import ChatRole
from app.utils.http_response import httpResponse
from app.http.requests.api.manager.chatController.CreateSessionRequest import CreateSessionRequest
from app.http.requests.api.manager.chatController.SendMessageRequest import SendMessageRequest
from app.http.requests.api.manager.chatController.GetSessionsRequest import GetSessionsRequest
from app.http.requests.api.manager.chatController.GetMessagesRequest import GetMessagesRequest
from app.http.requests.api.manager.chatController.DeleteSessionRequest import DeleteSessionRequest
from app.http.requests.api.manager.chatController.StreamMessageRequest import StreamMessageRequest
from app.dto.services.requests.manager.chat.GetSessionsDTO import GetSessionsDTO
from app.dto.services.requests.manager.chat.CreateSessionDTO import CreateSessionDTO
from app.dto.services.requests.manager.chat.GetMessagesDTO import GetMessagesDTO
from app.dto.services.requests.manager.chat.SendMessageDTO import SendMessageDTO
from app.dto.services.requests.manager.chat.DeleteSessionDTO import DeleteSessionDTO
import json
import asyncio
import threading
from queue import Queue
from datetime import datetime, timezone
from typing import Any, Dict, List

# Max rows per SQL result to include in narrative context (e.g. 90-day daily = 90 rows)
_MAX_ROWS_FOR_NARRATIVE = 120

NARRATIVE_SYSTEM_PROMPT = """Sen otel verileri üzerinden yöneticiye rapor sunan bir analistsin. Kullanıcının sorusu için sistem zaten veritabanından veriyi çekmiş; sana gelen "SQL çalıştırma sonuçları" bu çekilen verinin özetidir. Sen bu veriye bakıp yanıt üretiyorsun; kullanıcı sana ek veri sağlamaz.

Kurallar:
1. Yanıtı Türkçe, net ve profesyonel bir rapor diliyle yaz (yöneticiye sunum yapıyormuşsun gibi).
2. Verileri kullanarak somut sonuçlar ve varsa sayıları belirt. Sana gelen veri setinde hangi tarih aralığı / kaç satır varsa tümünü "eldeki veri" olarak kabul et ve ona göre yorumla.
3. Bazı sorgular HATA ile sonuçlanmışsa: O kısım için kısa ve net ifade et; kullanıcıdan asla "veri paylaşın", "veri gönderin" gibi taleplerde bulunma. Veri sistem tarafından çekilir; eksikse erişim sorunu olarak ifade et.
4. Tüm sorgular başarılıysa verileri kullanarak soruyu doğrudan ve güzelce cevapla.
5. Sadece metin yanıt ver; kod, JSON veya ekstra açıklama ekleme. Doğrudan rapor metnini yaz.
6. Asla teknik detay verme: Sütun adı (TOTAL_GUEST_IN_HOUSE, GUEST_IN_HOUSE vb.), None/null, "alan boş görünüyor", "sistemden çekilen X = None", SQL, hata kodu veya ham sistem çıktısı kullanma. Yanıt her zaman sadece iş diliyle, kullanıcının anlayacağı şekilde olsun; veri yoksa veya erişilemiyorsa bunu teknik jargonsuz kısaca belirt ve geç."""


def _format_sql_results_for_narrative(user_question: str, pipeline_result: Dict[str, Any]) -> str:
    """Build a single string: user question + SQL results summary for the narrative LLM."""
    parts = [f"Kullanıcı sorusu: {user_question}", "\n--- SQL çalıştırma sonuçları ---\n"]
    sql_results = pipeline_result.get("sql_results") or []
    if not sql_results:
        parts.append("Hiçbir alt istek veya SQL sonucu yok.\n")
    for sr in sql_results:
        req_desc = sr.get("request_description") or sr.get("clarified_description") or "İstek"
        parts.append(f"İstek: {req_desc}")
        exec_results = sr.get("sql_execution_results") or []
        for i, er in enumerate(exec_results):
            success = er.get("success", False)
            if success:
                data = er.get("data") or []
                total = len(data)
                rows = data[:_MAX_ROWS_FOR_NARRATIVE]
                if total > len(rows):
                    parts.append(f"  Sorgu {i+1}: Başarılı. Toplam {total} satır (aşağıda ilk {len(rows)} satır):")
                else:
                    parts.append(f"  Sorgu {i+1}: Başarılı. Toplam {total} satır:")
                parts.append(json.dumps(rows, ensure_ascii=False, default=str))
            else:
                parts.append(f"  Sorgu {i+1}: HATA - {er.get('error', 'Bilinmeyen hata')}")
        if sr.get("error"):
            parts.append(f"  İstek hatası: {sr['error']}")
        parts.append("")
    return "\n".join(parts)


class ChatController(Controller):
    def __init__(self, chatService: ChatService, ragService: RAGService = None):
        self.chatService = chatService
        self.ragService = ragService

    def getSessions(
        self,
        currentUser: dict,
        request: GetSessionsRequest,
        db: Session = Depends(get_db)
    ):
        return self.chatService.getSessions(GetSessionsDTO(currentUser, request))

    def getSession(self, sessionId: int, currentUser: dict):
        companyId = currentUser.get("company_id")
        userType = currentUser.get("user_type", "manager")
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
        currentUser: dict,
        request: CreateSessionRequest,
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

    async def streamMessage(
        self,
        request: StreamMessageRequest,
        currentUser: dict,
        db: Session = Depends(get_db)
    ):
        """Run RAG+SQL pipeline, then stream narrative answer via SSE (executive report style)."""
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

        userMessage = ChatMessage(
            chat_session_id=request.sessionId,
            role=ChatRole.USER,
            content=request.content
        )
        db.add(userMessage)
        db.commit()

        # RAG + SQL gen + execute (same as test module)
        testController = TestController(
            EmbeddingService(),
            OpenAIChatService(),
            self.ragService,
            db
        )
        try:
            pipelineResult = await testController.run_sql_pipeline(
                companyId, request.content, db
            )
        except Exception as e:
            async def error_gen():
                yield f"data: {json.dumps({'error': f'Veri analizi başlatılamadı: {str(e)}', 'done': True})}\n\n"
            return StreamingResponse(
                error_gen(),
                media_type="text/event-stream",
                headers={"Cache-Control": "no-cache", "Connection": "keep-alive", "X-Accel-Buffering": "no"}
            )

        narrativeContext = _format_sql_results_for_narrative(request.content, pipelineResult)
        narrativeMessages = [
            {"role": "system", "content": NARRATIVE_SYSTEM_PROMPT},
            {"role": "user", "content": narrativeContext}
        ]
        try:
            from app.utils.websocket_events import notifyMessageSentToLLM
            notifyMessageSentToLLM(request.sessionId, {"sessionId": request.sessionId})
        except Exception:
            pass
        chatService = OpenAIChatService()

        async def generate():
            try:
                fullContent = ""
                stream = chatService.chat_stream(narrativeMessages, temperature=0.5)
                
                # Use a queue to bridge sync stream and async generator
                chunkQueue = Queue()
                streamError = None
                
                def processStream():
                    """Process OpenAI stream synchronously and put chunks in queue"""
                    nonlocal streamError
                    try:
                        for chunk in stream:
                            if chunk.choices and len(chunk.choices) > 0:
                                delta = chunk.choices[0].delta
                                if hasattr(delta, 'content') and delta.content:
                                    chunkQueue.put(delta.content)
                        chunkQueue.put(None)  # Signal end of stream
                    except Exception as e:
                        streamError = str(e)
                        chunkQueue.put(None)
                
                # Start stream processing in background thread
                streamThread = threading.Thread(target=processStream, daemon=True)
                streamThread.start()
                
                # Yield chunks as they arrive
                while True:
                    try:
                        # Wait for chunk with timeout
                        contentChunk = await asyncio.get_event_loop().run_in_executor(
                            None, lambda: chunkQueue.get(timeout=0.1)
                        )
                        
                        if contentChunk is None:
                            # End of stream
                            break
                        
                        fullContent += contentChunk
                        # Send chunk via SSE immediately
                        data = json.dumps({"content": contentChunk, "done": False})
                        yield f"data: {data}\n\n"
                        
                    except Exception:
                        # Timeout - check if thread is still alive
                        if not streamThread.is_alive():
                            # Thread finished, check for remaining items
                            try:
                                while True:
                                    contentChunk = chunkQueue.get_nowait()
                                    if contentChunk is None:
                                        break
                                    fullContent += contentChunk
                                    data = json.dumps({"content": contentChunk, "done": False})
                                    yield f"data: {data}\n\n"
                            except:
                                pass
                            break
                        continue
                
                if streamError:
                    errorData = json.dumps({"error": streamError, "done": True})
                    yield f"data: {errorData}\n\n"
                    return
                
                # Save assistant message to database using a new session
                dbSession = SessionLocal()
                try:
                    # Get session from database
                    dbSessionObj = dbSession.query(ChatSession).filter(ChatSession.id == request.sessionId).first()
                    
                    assistantMessage = ChatMessage(
                        chat_session_id=request.sessionId,
                        role=ChatRole.ASSISTANT,
                        content=fullContent
                    )
                    dbSession.add(assistantMessage)
                    
                    # Update session title if not set
                    if dbSessionObj and not dbSessionObj.title:
                        newTitle = request.content[:50] + "..." if len(request.content) > 50 else request.content
                        dbSessionObj.title = newTitle
                    
                    # Event listener will automatically update updated_at on commit
                    dbSession.commit()
                    dbSession.refresh(assistantMessage)
                    
                    messageId = assistantMessage.id
                except Exception as dbError:
                    dbSession.rollback()
                    messageId = None
                finally:
                    dbSession.close()
                
                # Send done event
                data = json.dumps({"content": "", "done": True, "message_id": messageId})
                yield f"data: {data}\n\n"
                
            except Exception as e:
                errorData = json.dumps({"error": str(e), "done": True})
                yield f"data: {errorData}\n\n"

        return StreamingResponse(
            generate(),
            media_type="text/event-stream",
            headers={
                "Cache-Control": "no-cache",
                "Connection": "keep-alive",
                "X-Accel-Buffering": "no"
            }
        )
