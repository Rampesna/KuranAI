"""
User type definitions
"""


class UserType:
    ADMIN = "admin"
    MANAGER = "manager"
    EMPLOYEE = "employee"

    @staticmethod
    def is_valid(user_type: str) -> bool:
        return user_type in [UserType.ADMIN, UserType.MANAGER, UserType.EMPLOYEE]
