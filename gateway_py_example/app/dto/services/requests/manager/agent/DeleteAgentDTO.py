from typing import Dict
from app.http.requests.api.manager.agentController.DeleteAgentRequest import DeleteAgentRequest


class DeleteAgentDTO:
    def __init__(self, currentUser: Dict, request: DeleteAgentRequest):
        self._companyId = currentUser.get("company_id")
        self._agentId = request.agentId

    def getCompanyId(self) -> int:
        return self._companyId

    def getAgentId(self) -> int:
        return self._agentId
