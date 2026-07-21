from app.services.openai.ChatService import ChatService
from app.models.setting import Setting
from sqlalchemy.orm import Session
from typing import Dict, Any, Optional
import json


class RAGService:
    def __init__(self, chatService: ChatService):
        self.chatService = chatService

    def analyzeUserPrompt(self, userMessage: str, companyId: int = 1, db: Optional[Session] = None) -> Dict[str, Any]:
        # Get system prompt from settings table
        systemPromptContent = None
        
        if db:
            try:
                setting = db.query(Setting).filter(
                    Setting.company_id == companyId,
                    Setting.key == "prompt_analysis_system_prompt",
                    Setting.category == "rag",
                    Setting.is_active == True,
                    Setting.deleted_at.is_(None)
                ).first()
                
                if setting:
                    systemPromptContent = setting.value
            except Exception as e:
                # If error fetching from DB, use default
                pass
        
        # Fallback to default prompt if not found in settings
        if not systemPromptContent:
            systemPromptContent = """You are a prompt analysis system. Analyze the user's prompt and determine:
1. Is this a single request or multiple requests?
2. If single, how many SQL queries are needed?
3. If multiple, break down each request and identify SQL queries needed for each.
4. For each request, if it contains date-related information, clarify the exact date range needed.

IMPORTANT: If the request mentions relative dates like "bugün", "today", "bu ay", "this month", "bu yıl", "this year", 
the user message will include the actual current date. Use that information to clarify the exact date range in the description.

Return ONLY valid JSON in this format:
{
    "is_single_request": true/false,
    "request_count": number,
    "requires_multiple_sql": true/false,
    "requests": [
        {
            "index": number (starting from 1),
            "description": "clarified description of request with exact date ranges if applicable",
            "sql_queries_needed": number,
            "sql_descriptions": ["brief description of what each SQL query should return"]
        }
    ]
}"""

        messages = [
            {"role": "system", "content": systemPromptContent},
            {"role": "user", "content": userMessage}
        ]

        response = self.chatService.chat(messages, temperature=0.3)
        
        try:
            return json.loads(response)
        except json.JSONDecodeError:
            return {
                "is_single_request": True,
                "request_count": 1,
                "requests": [
                    {
                        "description": userMessage,
                        "sql_queries_needed": 1
                    }
                ]
            }
