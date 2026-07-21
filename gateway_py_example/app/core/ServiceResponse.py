from typing import TypeVar, Generic, Optional

T = TypeVar('T')


class ServiceResponse(Generic[T]):
    def __init__(
        self,
        isSuccess: bool,
        message: str,
        statusCode: int,
        data: Optional[T] = None
    ):
        self._is_success = isSuccess
        self._message = message
        self._status_code = statusCode
        self._data = data

    def isSuccess(self) -> bool:
        return self._is_success

    def getIsSuccess(self) -> bool:
        return self._is_success

    def getMessage(self) -> str:
        return self._message

    def getStatusCode(self) -> int:
        return self._status_code

    def getData(self) -> Optional[T]:
        return self._data

    def to_dict(self) -> dict:
        return {
            "isSuccess": self._is_success,
            "message": self._message,
            "statusCode": self._status_code,
            "data": self._data
        }
