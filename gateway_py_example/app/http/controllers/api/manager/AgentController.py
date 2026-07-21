from fastapi import Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.core.controller import Controller
from app.utils.http_response import httpResponse
from app.services.sqlalchemy.AgentService import AgentService
from app.http.requests.api.manager.agentController.CreateAgentRequest import CreateAgentRequest
from app.http.requests.api.manager.agentController.UpdateAgentRequest import UpdateAgentRequest
from app.dto.services.requests.manager.agent.ListAgentsDTO import ListAgentsDTO
from app.dto.services.requests.manager.agent.GetAgentDTO import GetAgentDTO
from app.dto.services.requests.manager.agent.CreateAgentDTO import CreateAgentDTO
from app.dto.services.requests.manager.agent.UpdateAgentDTO import UpdateAgentDTO
from app.dto.services.requests.manager.agent.DeleteAgentDTO import DeleteAgentDTO
from app.http.requests.api.manager.agentController.ListAgentsRequest import ListAgentsRequest
from app.http.requests.api.manager.agentController.GetAgentRequest import GetAgentRequest
from app.http.requests.api.manager.agentController.DeleteAgentRequest import DeleteAgentRequest
from typing import Optional


class AgentController(Controller):
    def __init__(self, agentService: AgentService):
        self.agentService = agentService

    def listAgents(self, request: ListAgentsRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.agentService.listAgents(ListAgentsDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def getAgent(self, request: GetAgentRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.agentService.getAgent(GetAgentDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def createAgent(self, request: CreateAgentRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.agentService.createAgent(CreateAgentDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def updateAgent(self, request: UpdateAgentRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.agentService.updateAgent(UpdateAgentDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )

    def deleteAgent(self, request: DeleteAgentRequest, currentUser: dict, db: Session = Depends(get_db)):
        response = self.agentService.deleteAgent(DeleteAgentDTO(currentUser, request))
        return httpResponse(
            isSuccess=response.getIsSuccess(),
            message=response.getMessage(),
            statusCode=response.getStatusCode(),
            data=response.getData()
        )
