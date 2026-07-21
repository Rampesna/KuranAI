from openai import OpenAI
from app.core.config import settings
from typing import List
import os


class EmbeddingService:
    _client = None
    _model = None
    
    def __init__(self):
        if EmbeddingService._client is None:
            # Remove proxies from env to avoid httpx version conflict
            old_proxies = os.environ.pop("HTTP_PROXY", None)
            old_https_proxies = os.environ.pop("HTTPS_PROXY", None)
            try:
                # Initialize OpenAI client
                EmbeddingService._client = OpenAI(
                    api_key=settings.OPENAI_API_KEY,
                    timeout=60.0
                )
            finally:
                # Restore proxies if they existed
                if old_proxies:
                    os.environ["HTTP_PROXY"] = old_proxies
                if old_https_proxies:
                    os.environ["HTTPS_PROXY"] = old_https_proxies
            EmbeddingService._model = settings.OPENAI_EMBEDDING_MODEL
        self.client = EmbeddingService._client
        self.model = EmbeddingService._model

    def embed(self, text: str) -> List[float]:
        response = self.client.embeddings.create(
            model=self.model,
            input=text
        )
        return response.data[0].embedding
