from openai import OpenAI
from app.core.config import settings
from typing import List, Dict, Optional
import os


class ChatService:
    _client = None
    _model = None
    
    def __init__(self):
        if ChatService._client is None:
            # Remove proxies from env to avoid httpx version conflict
            old_proxies = os.environ.pop("HTTP_PROXY", None)
            old_https_proxies = os.environ.pop("HTTPS_PROXY", None)
            try:
                # Initialize OpenAI client
                ChatService._client = OpenAI(
                    api_key=settings.OPENAI_API_KEY,
                    timeout=60.0
                )
            finally:
                # Restore proxies if they existed
                if old_proxies:
                    os.environ["HTTP_PROXY"] = old_proxies
                if old_https_proxies:
                    os.environ["HTTPS_PROXY"] = old_https_proxies
            ChatService._model = settings.OPENAI_MODEL
        self.client = ChatService._client
        self.model = ChatService._model

    def chat(self, messages: List[Dict[str, str]], temperature: float = 0.7) -> str:
        response = self.client.chat.completions.create(
            model=self.model,
            messages=messages,
            temperature=temperature
        )
        return response.choices[0].message.content
    
    def chat_stream(self, messages: List[Dict[str, str]], temperature: float = 0.7):
        """Stream chat completion from OpenAI"""
        stream = self.client.chat.completions.create(
            model=self.model,
            messages=messages,
            temperature=temperature,
            stream=True
        )
        return stream
