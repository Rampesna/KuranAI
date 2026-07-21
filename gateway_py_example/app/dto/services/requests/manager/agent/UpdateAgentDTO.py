from typing import Dict, Optional
from app.http.requests.api.manager.agentController.UpdateAgentRequest import UpdateAgentRequest


class UpdateAgentDTO:
    def __init__(self, currentUser: Dict, request: UpdateAgentRequest):
        self._companyId = currentUser.get("company_id")
        self._agentId = request.agentId
        self._name = request.name
        self._description = request.description
        self._scheduleType = request.schedule_type
        self._scheduleTime = request.schedule_time
        self._scheduleDay = request.schedule_day
        self._scheduleDate = request.schedule_date
        self._reports = request.reports
        self._status = request.status

    def getCompanyId(self) -> int:
        return self._companyId

    def getAgentId(self) -> int:
        return self._agentId

    def getName(self) -> Optional[str]:
        return self._name

    def getDescription(self) -> Optional[str]:
        return self._description

    def getScheduleType(self) -> Optional[str]:
        return self._scheduleType

    def getScheduleTime(self) -> Optional[str]:
        return self._scheduleTime

    def getScheduleDay(self) -> Optional[str]:
        return self._scheduleDay

    def getScheduleDate(self) -> Optional[int]:
        return self._scheduleDate

    def getReports(self) -> Optional[list]:
        return self._reports

    def getStatus(self) -> Optional[str]:
        return self._status

    def toDict(self) -> Dict:
        data = {}
        if self._name is not None:
            data["name"] = self._name
        if self._description is not None:
            data["description"] = self._description
        if self._scheduleType is not None:
            data["schedule_type"] = self._scheduleType
        if self._scheduleTime is not None:
            data["schedule_time"] = self._scheduleTime
        if self._scheduleDay is not None:
            data["schedule_day"] = self._scheduleDay
        if self._scheduleDate is not None:
            data["schedule_date"] = self._scheduleDate
        if self._reports is not None:
            data["reports"] = self._reports
        if self._status is not None:
            data["status"] = self._status
        return data
