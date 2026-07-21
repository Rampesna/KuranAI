"""
Chat role definitions
"""


class ChatRole:
    USER = "user"
    ASSISTANT = "assistant"
    SYSTEM = "system"

    @staticmethod
    def is_valid(chat_role: str) -> bool:
        return chat_role in [ChatRole.USER, ChatRole.ASSISTANT, ChatRole.SYSTEM]
