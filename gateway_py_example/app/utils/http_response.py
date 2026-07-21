from fastapi import Response, status
import json


def httpResponse(
    isSuccess: bool = False,
    message: str = "",
    statusCode: int = 400,
    data=None
) -> Response:
    responseData = {
        "isSuccess": isSuccess,
        "message": message,
        "statusCode": statusCode,
        "data": data
    }
    
    return Response(
        content=json.dumps(responseData, default=str),
        status_code=statusCode,
        media_type="application/json"
    )
