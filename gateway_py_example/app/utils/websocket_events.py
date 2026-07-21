"""
Notify Node.js WebSocket server (Socket.IO) on chat events.
Backend calls HTTP POST to websocket_server; it broadcasts to connected clients.
"""
import os
import logging

logger = logging.getLogger(__name__)

BASE_URL = os.getenv("WEBSOCKET_SERVER_URL", "http://websocket_server:4575")


def _post(path: str, json_body: dict) -> None:
    try:
        import httpx
        with httpx.Client(timeout=3.0) as client:
            r = client.post(f"{BASE_URL}{path}", json=json_body)
            if r.status_code != 200:
                logger.warning("WebSocket server %s returned %s", path, r.status_code)
    except Exception as e:
        logger.warning("WebSocket server notify failed %s: %s", path, e)


def notifySessionCreated(sessionId: int) -> None:
    _post("/events/onChatSessionCreated", {"session_id": sessionId})


def notifySessionDeleted(sessionId: int) -> None:
    _post("/events/onChatSessionDeleted", {"session_id": sessionId})


def notifyMessageSentToLLM(sessionId: int, data: dict = None) -> None:
    body = {"session_id": sessionId}
    if data:
        body.update(data)
    _post("/events/onChatMessageSentToLLM", body)
