from typing import Dict
from app.http.requests.api.manager.agentController.GetAgentRequest import GetAgentRequest


class GetAgentDTO:
    def __init__(self, currentUser: Dict, request: GetAgentRequest):
        self._companyId = currentUser.get("company_id")
        self._agentId = request.agentId

    def getCompanyId(self) -> int:
        return self._companyId

    def getAgentId(self) -> int:
        return self._agentId
