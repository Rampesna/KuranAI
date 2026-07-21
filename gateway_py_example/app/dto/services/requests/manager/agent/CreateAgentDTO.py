from typing import Dict
from app.http.requests.api.manager.agentController.CreateAgentRequest import CreateAgentRequest


class CreateAgentDTO:
    def __init__(self, currentUser: Dict, request: CreateAgentRequest):
        self._companyId = currentUser.get("company_id")
        self._name = request.name
        self._description = request.description
        self._scheduleType = request.schedule_type
        self._scheduleTime = request.schedule_time
        self._scheduleDay = request.schedule_day
        self._scheduleDate = request.schedule_date
        self._reports = request.reports
        self._status = request.status or "active"

    def getCompanyId(self) -> int:
        return self._companyId

    def getName(self) -> str:
        return self._name

    def getDescription(self) -> str:
        return self._description

    def getScheduleType(self) -> str:
        return self._scheduleType

    def getScheduleTime(self) -> str:
        return self._scheduleTime

    def getScheduleDay(self) -> str:
        return self._scheduleDay

    def getScheduleDate(self) -> int:
        return self._scheduleDate

    def getReports(self) -> list:
        return self._reports

    def getStatus(self) -> str:
        return self._status

    def toDict(self) -> Dict:
        return {
            "name": self._name,
            "description": self._description,
            "schedule_type": self._scheduleType,
            "schedule_time": self._scheduleTime,
            "schedule_day": self._scheduleDay,
            "schedule_date": self._scheduleDate,
            "reports": self._reports,
            "status": self._status
        }
