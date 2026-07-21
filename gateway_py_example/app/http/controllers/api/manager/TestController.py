from fastapi import Depends
from sqlalchemy.orm import Session
from app.core.database import get_db, SessionLocal
from app.core.controller import Controller
from app.utils.http_response import httpResponse
from app.services.openai.EmbeddingService import EmbeddingService
from app.services.openai.ChatService import ChatService
from app.services.intelor_services.RAGService import RAGService
from app.http.requests.api.manager.authController.TestOpenaiRequest import TestOpenAIRequest
from app.models.company import Company
import asyncio
import json
from concurrent.futures import ThreadPoolExecutor
from typing import Dict, Any, List, Optional, Tuple
from datetime import datetime, timezone, timedelta
import os
import re
import oracledb


class TestController(Controller):
    def __init__(
        self,
        embeddingService: EmbeddingService,
        chatService: ChatService,
        ragService: RAGService,
        db: Session
    ):
        self.embeddingService = embeddingService
        self.chatService = chatService
        self.ragService = ragService
        self.db = db

    def _loadREPManagerContext(self, company_id: int) -> str:
        """Load context for the given company. Uses company_contexts/company_{company_id}.json if present; else continues without (no error)."""
        try:
            # Path to the TOON context file
            contextPath = os.path.join(
                os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))),
                "New_Oracle_Schemas",
                "REP_MANAGER",
                "REP_MANAGER_TOON_Context.txt"
            )
            
            with open(contextPath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Parse the custom format (key,value pairs)
            lines = content.split('\n')
            contextText = "REP_MANAGER Table Context:\n\n"
            
            tableName = ""
            tableDesc = ""
            primaryKey = []
            columns = []
            currentColumn = None
            columnInfo = {}
            inColumn = False
            indentLevel = 0
            commonCalculations = []
            viewDefinitions = []
            canonicalColumnNames: List[str] = []
            columnNameRules: Dict[str, str] = {}
            inColumnNameRules = False
            
            for line in lines:
                line = line.rstrip()
                if not line:
                    continue
                
                # Detect indent level
                indentLevel = len(line) - len(line.lstrip())
                line = line.strip()
                
                if ',' in line:
                    parts = line.split(',', 1)
                    key = parts[0].strip()
                    value = parts[1].strip() if len(parts) > 1 else ""
                    
                    # Table metadata
                    if key == "table_name":
                        tableName = value
                        inColumnNameRules = False
                        contextText += f"Table Name: {value}\n"
                    elif key == "table_description":
                        tableDesc = value
                        contextText += f"Description: {value}\n\n"
                    elif key.startswith("canonical_column_names"):
                        inColumnNameRules = False
                        if value:
                            canonicalColumnNames = [c.strip() for c in value.split(",") if c.strip()]
                    elif key == "column_name_rules":
                        inColumnNameRules = True
                    elif inColumnNameRules and indentLevel == 2:
                        if key == "critical":
                            columnNameRules["critical"] = value
                        elif "do_not" in key:
                            columnNameRules["do_not"] = value
                        elif key == "do":
                            columnNameRules["do"] = value
                    elif indentLevel == 0 and key != "column_name_rules":
                        inColumnNameRules = False
                    if key == "primary_key":
                        contextText += "Primary Key: "
                    elif key.startswith("columns[") and value:
                        # Primary key columns
                        pkCols = [c.strip() for c in value.split(',')]
                        primaryKey.extend(pkCols)
                        contextText += ", ".join(pkCols) + "\n"
                    elif key == "description" and indentLevel == 2 and "Primary Key" in contextText:
                        contextText += f"  {value}\n\n"
                    elif key == "columns":
                        contextText += "Available Columns:\n"
                    elif key == "common_calculations[" and value:
                        # Parse common calculations
                        calcParts = value.split(',')
                        commonCalculations.extend([c.strip() for c in calcParts])
                    elif key.startswith("view_definitions[") and value:
                        # Parse view definitions (metrics)
                        defParts = value.split(',')
                        viewDefinitions.extend([d.strip() for d in defParts])
                    elif key and key.isupper() and "_" in key and indentLevel == 2:
                        # New column definition
                        if currentColumn and columnInfo:
                            columns.append(columnInfo)
                        currentColumn = key
                        columnInfo = {"name": key, "type": "", "description": "", "sample": ""}
                        inColumn = True
                    elif inColumn and currentColumn:
                        if key == "type" and indentLevel == 4:
                            columnInfo["type"] = value
                        elif key == "description" and indentLevel == 4:
                            columnInfo["description"] = value
                        elif key == "sample_value" and indentLevel == 6:
                            columnInfo["sample"] = value
            
            # Add last column
            if currentColumn and columnInfo:
                columns.append(columnInfo)
            
            # Load company-specific context JSON: company_contexts/company_{company_id}.json (no error if missing)
            # base_dir = gateway_py project root (6 levels up from this file)
            if not canonicalColumnNames:
                base_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))))
                company_json_path = os.path.join(base_dir, "company_contexts", f"company_{company_id}.json")
                try:
                    if os.path.isfile(company_json_path):
                        with open(company_json_path, "r", encoding="utf-8") as jf:
                            j = json.load(jf)
                        canonicalColumnNames = j.get("canonical_column_names") or []
                        r = j.get("column_name_rules") or {}
                        if isinstance(r, dict):
                            columnNameRules = {k: (v if isinstance(v, str) else "") for k, v in r.items()}
                except Exception:
                    pass
            
            # Build and inject ALLOWED COLUMNS block – ONLY this list, nothing else
            if canonicalColumnNames:
                canonicalBlock = "\n🚨 ALLOWED COLUMN NAMES – YOU MUST USE ONLY THESE 🚨\n\n"
                canonicalBlock += "The ONLY column names that exist are the ones in the list below. Use ONLY these. "
                canonicalBlock += "If a name is not in this list, it does not exist – never use it.\n\n"
                canonicalBlock += "Allowed columns (" + str(len(canonicalColumnNames)) + "):\n"
                canonicalBlock += ", ".join(canonicalColumnNames) + "\n\n"
                canonicalBlock += "Request → column (use ONLY these from the list): "
                canonicalBlock += "children / çocuk → CHILDREN_IN_HOUSE. Food revenue / yemek geliri → FOOD_BEV_REVENUE. Room revenue → ROOM_REVENUE. "
                canonicalBlock += "Kanal / channel / satış kanalı → SOURCE_ROOMS_TODAY, COMPANY_ROOMS_TODAY, AGENT_ROOMS_TODAY, GROUP_ROOM. Revenue by channel: COMPANY_ROOM_REVENUE, AGENT_ROOM_REVENUE, GROUP_ROOM_REV (exact name, not GROUP_ROOM_REVENUE). No SOURCE_ROOM_REVENUE or COMMISSION column; use TOTAL_REVENUE, COMPANY_ROOM_REVENUE, AGENT_ROOM_REVENUE, GROUP_ROOM_REV only. "
                canonicalBlock += "Segment (grup, kurumsal, bireysel) → GROUP_ROOM, COMPANY_ROOMS_TODAY, FIT_GUESTS_TODAY and GROUP_ROOM_REV, COMPANY_ROOM_REVENUE, FIT_ROOM_REVENUE. "
                canonicalBlock += "İptal / cancellation → CANCEL_ROOMS, CANCEL_RESERVATION, LATE_CANCEL_ROOMS, LATE_CANCEL_RESERVATION. No-show → NOSHOW_ROOM, NOSHOW_PRS, NOSHOW_RESERVATIONS (exact; not NO_SHOW_ROOMS). "
                canonicalBlock += "LOS / konaklama süresi → FIT_MEM_LOS_NIGHTS, FIT_MEM_LOS_RESV, BLK_MEM_LOS_NIGHTS, BLK_MEM_LOS_RESV. "
                canonicalBlock += "Otelde kaç kişi / in-house guests / persons in house → GUEST_IN_HOUSE only (do NOT use INHOUSE_PRS, OCC_PRS – they do not exist). "
                canonicalBlock += "Kârlılık / revenue → ROOM_REVENUE, TOTAL_REVENUE, COMPANY_ROOM_REVENUE, AGENT_ROOM_REVENUE, GROUP_ROOM_REV.\n\n"
                marker = f"Description: {tableDesc}\n\n"
                idx = contextText.find(marker)
                if idx != -1:
                    insertAt = idx + len(marker)
                    contextText = contextText[:insertAt] + canonicalBlock + contextText[insertAt:]
            
            # Format columns for context (limit to 150 most important columns)
            contextText += "\nColumn Details (showing key columns):\n"
            
            # Prioritize important columns
            importantColumns = [
                "BUSINESS_DATE", "ROOM_CLASS", "RESORT", "OCC_ROOM", "PHYSICAL_ROOM",
                "ROOM_REVENUE", "FOOD_BEV_REVENUE", "OTHER_REVENUE", "TOTAL_REVENUE",
                "CHILDREN_IN_HOUSE", "ADULTS_IN_HOUSE", "GUEST_IN_HOUSE",
                "ADR_ROOM_WO_COMP_HOUSE", "PER_OCCUPANCY",
                "ARRIVAL_ROOM", "DEPARTURE_ROOM", "COMP_ROOM", "HOUSE_USE_ROOM"
            ]
            
            # Add important columns first
            addedColumns = set()
            for colName in importantColumns:
                for col in columns:
                    if col["name"] == colName and colName not in addedColumns:
                        colType = col.get("type", "")
                        colDesc = col.get("description", "")
                        colSample = col.get("sample", "")
                        contextText += f"\n- {colName} ({colType}): {colDesc}"
                        if colSample:
                            contextText += f" [Example: {colSample}]"
                        addedColumns.add(colName)
            
            # Add remaining columns (up to 150 total)
            remainingCount = 150 - len(addedColumns)
            for col in columns:
                if remainingCount <= 0:
                    break
                if col["name"] not in addedColumns:
                    colName = col.get("name", "")
                    colType = col.get("type", "")
                    colDesc = col.get("description", "")
                    colSample = col.get("sample", "")
                    contextText += f"\n- {colName} ({colType}): {colDesc}"
                    if colSample:
                        contextText += f" [Example: {colSample}]"
                    addedColumns.add(colName)
                    remainingCount -= 1
            
            contextText += f"\n\nTotal Columns: {len(columns)}"
            
            # Add calculated metrics section - CRITICAL for LLM to understand
            contextText += "\n\n🚨🚨🚨 CALCULATED METRICS (NOT DIRECT COLUMNS) 🚨🚨🚨"
            contextText += "\n\n⚠️ CRITICAL WARNING: The following metrics are NOT direct columns in the table!"
            contextText += "\n⚠️ You MUST calculate them using the formulas below. DO NOT try to SELECT them as columns!"
            contextText += "\n\n"
            
            if commonCalculations:
                for calc in commonCalculations:
                    if calc:
                        if "ADR" in calc.upper() or "REVPAR" in calc.upper() or "RevPar" in calc:
                            contextText += f"🚨 {calc}\n"
                            if "ADR" in calc.upper():
                                contextText += "   SQL Calculation: CASE WHEN OCC_ROOM = 0 THEN 0 ELSE ROOM_REVENUE / OCC_ROOM END\n"
                            if "REVPAR" in calc.upper() or "RevPar" in calc:
                                contextText += "   SQL Calculation: CASE WHEN PHYSICAL_ROOM = 0 THEN 0 ELSE ROOM_REVENUE / PHYSICAL_ROOM END\n"
                        else:
                            contextText += f"- {calc}\n"
            else:
                # Fallback to known calculations from View_Query_Example.sql
                contextText += "🚨 Occupancy % = (OCC_ROOM / PHYSICAL_ROOM) * 100\n"
                contextText += "   SQL Calculation: CASE WHEN PHYSICAL_ROOM = 0 THEN 0 ELSE (OCC_ROOM / PHYSICAL_ROOM) * 100 END\n\n"
                contextText += "🚨 ADR (Average Daily Rate) = ROOM_REVENUE / OCC_ROOM\n"
                contextText += "   SQL Calculation: CASE WHEN OCC_ROOM = 0 THEN 0 ELSE ROOM_REVENUE / OCC_ROOM END\n"
                contextText += "   ⚠️ DO NOT SELECT ADR as a column - it doesn't exist!\n\n"
                contextText += "🚨 RevPar (Revenue per Available Room) = ROOM_REVENUE / PHYSICAL_ROOM\n"
                contextText += "   SQL Calculation: CASE WHEN PHYSICAL_ROOM = 0 THEN 0 ELSE ROOM_REVENUE / PHYSICAL_ROOM END\n"
                contextText += "   ⚠️ DO NOT SELECT REVPAR or RevPar as columns - they don't exist!\n"
            
            contextText += "\n\n=== AVAILABLE METRICS (from view definitions) ==="
            if viewDefinitions:
                contextText += "\nThe following metrics are commonly requested (some are calculated, some are direct columns):\n"
                for metric in viewDefinitions[:20]:  # Limit to first 20
                    if metric:
                        contextText += f"- {metric}\n"
            
            contextText += "\n\n=== IMPORTANT NOTES ==="
            contextText += "\n- Use BUSINESS_DATE column for date filtering (DATE type)"
            contextText += "\n- Primary Key: " + ", ".join(primaryKey)
            contextText += "\n- This table contains daily aggregated hotel operational metrics"
            contextText += "\n- Date filtering example: WHERE BUSINESS_DATE >= TO_DATE('2025-12-01', 'YYYY-MM-DD') AND BUSINESS_DATE < TO_DATE('2026-01-01', 'YYYY-MM-DD')"
            contextText += "\n- When calculating ADR or RevPar, use the formulas above. DO NOT try to select ADR or REVPAR as columns - they don't exist!"
            contextText += "\n- Always handle division by zero: CASE WHEN denominator = 0 THEN 0 ELSE numerator / denominator END"
            contextText += "\n\n⚠️ CRITICAL: REVENUE CALCULATIONS ⚠️"
            contextText += "\n- ROOM_REVENUE, FOOD_BEV_REVENUE, OTHER_REVENUE, TOTAL_REVENUE values in the table are in the original currency"
            contextText += "\n- If currency conversion is needed, multiply by exchange rate (e.g., room_revenue * exchange_rate)"
            contextText += "\n- For MTD (Month-to-Date) calculations: Use TRUNC(BUSINESS_DATE, 'MONTH') to get month start"
            contextText += "\n- For date ranges: Use BETWEEN with inclusive start and exclusive end (e.g., >= start_date AND < end_date + 1)"
            contextText += "\n- When aggregating revenue metrics, SUM the values: SUM(ROOM_REVENUE) / SUM(PHYSICAL_ROOM) for RevPar"
            contextText += "\n- Note: View calculations typically use 'dün'e kadar' (yesterday) for MTD, not including today"
            
            return contextText
            
        except Exception as e:
            return f"Error loading context: {str(e)}"

    def _getCurrentDateInfo(self) -> Dict[str, Any]:
        """Get current date information for relative date resolution"""
        now = datetime.now(timezone.utc)
        return {
            "today": now.strftime("%Y-%m-%d"),
            "current_year": now.year,
            "current_month": now.month,
            "current_day": now.day,
            "current_week_start": (now - timedelta(days=now.weekday())).strftime("%Y-%m-%d"),
            "current_month_start": now.replace(day=1).strftime("%Y-%m-%d"),
            "current_year_start": now.replace(month=1, day=1).strftime("%Y-%m-%d"),
            "yesterday": (now - timedelta(days=1)).strftime("%Y-%m-%d"),
            "last_month_start": ((now.replace(day=1) - timedelta(days=1)).replace(day=1)).strftime("%Y-%m-%d"),
            "last_month_end": (now.replace(day=1) - timedelta(days=1)).strftime("%Y-%m-%d"),
            "last_year_start": now.replace(year=now.year-1, month=1, day=1).strftime("%Y-%m-%d"),
            "last_year_end": now.replace(year=now.year-1, month=12, day=31).strftime("%Y-%m-%d")
        }

    def _clarifyDatesInRequest(self, requestDescription: str, dateInfo: Dict[str, Any]) -> str:
        """Clarify relative dates in request description with actual dates"""
        clarified = requestDescription
        
        # Replace relative date terms with actual dates
        dateReplacements = {
            "bugün": f"{dateInfo['today']} (bugün)",
            "bu gün": f"{dateInfo['today']} (bugün)",
            "today": f"{dateInfo['today']} (today)",
            "bu ay": f"{dateInfo['current_month_start']} ile {dateInfo['today']} arası (bu ay)",
            "this month": f"{dateInfo['current_month_start']} to {dateInfo['today']} (this month)",
            "bu yıl": f"{dateInfo['current_year_start']} ile {dateInfo['today']} arası (bu yıl)",
            "this year": f"{dateInfo['current_year_start']} to {dateInfo['today']} (this year)",
            "geçen ay": f"{dateInfo['last_month_start']} ile {dateInfo['last_month_end']} arası (geçen ay)",
            "last month": f"{dateInfo['last_month_start']} to {dateInfo['last_month_end']} (last month)",
            "geçen yıl": f"{dateInfo['last_year_start']} ile {dateInfo['last_year_end']} arası (geçen yıl)",
            "last year": f"{dateInfo['last_year_start']} to {dateInfo['last_year_end']} (last year)",
            "dün": f"{dateInfo['yesterday']} (dün)",
            "yesterday": f"{dateInfo['yesterday']} (yesterday)"
        }
        
        for term, replacement in dateReplacements.items():
            if term.lower() in clarified.lower():
                clarified = re.sub(re.escape(term), replacement, clarified, flags=re.IGNORECASE)
        
        return clarified

    def _processSingleRequestSync(
        self,
        req: Dict[str, Any],
        companyId: int,
        db: Session,
        repManagerContext: str,
        dateInfo: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Process a single request: analyze dates, generate SQL with REP_MANAGER context"""
        requestDescription = req.get("description", "")
        requestIndex = req.get("index", 0)
        sqlCount = req.get("sql_queries_needed", 1)
        
        try:
            # 1. Clarify dates in request description
            clarifiedDescription = self._clarifyDatesInRequest(requestDescription, dateInfo)
            
            # 2. Get company database schema information
            connectionInfo = self._getCompanyDatabaseConnection(companyId, db)
            schemaPrefix = ""
            if connectionInfo and connectionInfo.get("schema"):
                schemaPrefix = connectionInfo["schema"] + "."
            
            # 3. Generate SQL query using LLM with REP_MANAGER context
            sqlGenerationMessages = [
                {
                    "role": "system",
                    "content": """You are a SQL expert for Oracle databases (OperaPMS). 
Generate ONLY valid Oracle SQL SELECT queries based on the user's request and the provided REP_MANAGER table context.

🚨 COLUMNS: USE ONLY THE ALLOWED LIST 🚨
- Use ONLY column names from the "ALLOWED COLUMN NAMES" list in the table context. Never use any other name.
- If a name is not in that list, it does not exist. Map requests to columns from that list only.
- Children → CHILDREN_IN_HOUSE. Food revenue → FOOD_BEV_REVENUE. Room revenue → ROOM_REVENUE. Kanal revenue: COMPANY_ROOM_REVENUE, AGENT_ROOM_REVENUE, GROUP_ROOM_REV (exact; no SOURCE_ROOM_REVENUE or GROUP_ROOM_REVENUE). No COMMISSION column. Segment → GROUP_ROOM, COMPANY_ROOMS_TODAY, FIT_GUESTS_TODAY. İptal → CANCEL_ROOMS, CANCEL_RESERVATION. No-show → NOSHOW_ROOM, NOSHOW_PRS, NOSHOW_RESERVATIONS (not NO_SHOW_ROOMS). LOS → FIT_MEM_LOS_NIGHTS, FIT_MEM_LOS_RESV, BLK_MEM_LOS_NIGHTS. Room category → ROOM_CLASS (not ROOM_CATEGORY).
- Otelde kaç kişi / in-house guests / persons in house → GUEST_IN_HOUSE only. Do NOT use INHOUSE_PRS, OCC_PRS, INHOUSE_GUESTS – they do not exist. Use GUEST_IN_HOUSE.

🚨 CRITICAL: CALCULATED METRICS ARE NOT COLUMNS 🚨
- ADR (Average Daily Rate) is NOT a column! You MUST calculate it as: CASE WHEN OCC_ROOM = 0 THEN 0 ELSE ROOM_REVENUE / OCC_ROOM END
- RevPar/REVPAR is NOT a column! You MUST calculate it as: CASE WHEN PHYSICAL_ROOM = 0 THEN 0 ELSE ROOM_REVENUE / PHYSICAL_ROOM END
- Occupancy % is NOT a column! You MUST calculate it as: CASE WHEN PHYSICAL_ROOM = 0 THEN 0 ELSE (OCC_ROOM / PHYSICAL_ROOM) * 100 END
- NEVER try to SELECT ADR, REVPAR, or RevPar as columns – they are calculated. Use the formulas provided. All other columns must be from the allowed list only.

IMPORTANT DATE INFORMATION:
- Today's date: """ + dateInfo['today'] + """
- Current year: """ + str(dateInfo['current_year']) + """
- Current month: """ + str(dateInfo['current_month']) + """
- When user mentions "bugün", "today", "bu ay", "this month", "bu yıl", "this year", use the actual dates provided in the request.

Rules:
0. Use ONLY column names from the allowed list in the table context. Never use a name that is not in that list.
1. Return ONLY the SQL query/queries, nothing else
2. MUST be SELECT queries only
3. ⚠️ CRITICAL: NEVER use ROWNUM with aggregate functions (SUM, AVG, COUNT, MAX, MIN, etc.)!
   - If using aggregate functions (SUM, AVG, COUNT, etc.), DO NOT add ROWNUM <= 100
   - Aggregate functions already return a single row, so ROWNUM would incorrectly limit the data BEFORE aggregation
   - ROWNUM should ONLY be used when selecting individual rows without aggregation
   - Example WRONG: SELECT SUM(ROOM_REVENUE) FROM REP_MANAGER WHERE ... AND ROWNUM <= 100
   - Example CORRECT: SELECT SUM(ROOM_REVENUE) FROM REP_MANAGER WHERE ...
4. NO comments or explanations
5. Query must be valid Oracle SQL syntax. Every non-aggregated column in SELECT must appear in GROUP BY. No trailing comma. Subqueries: matching parentheses.
6. Use proper Oracle date functions (e.g., TRUNC for dates, TO_DATE for date strings)
7. Use proper Oracle string functions if needed
8. Filter by BUSINESS_DATE column for date ranges
9. If multiple SQL queries are needed for this request, return them separated by semicolons (;)
10. For ADR: Calculate as ROOM_REVENUE / OCC_ROOM (with NULLIF or CASE to handle division by zero)
11. For RevPar: Calculate as ROOM_REVENUE / PHYSICAL_ROOM (with NULLIF or CASE to handle division by zero)
12. When calculating metrics (ADR, RevPar, Occupancy %), you MUST use aggregate functions (SUM) to get ALL data, never limit with ROWNUM
13. ⚠️ SCHEMA PREFIX: """ + (f"ALWAYS use schema prefix '{schemaPrefix.rstrip('.')}' when referencing tables. Example: Use {schemaPrefix}REP_MANAGER instead of REP_MANAGER" if schemaPrefix else "No schema prefix required. Use table names directly (e.g., REP_MANAGER)") + """
13. ⚠️ IMPORTANT: When comparing values (e.g., "this month vs last month", "this year vs last year"), ALWAYS return BOTH the current and previous values along with the comparison result!
   - Example: If comparing RevPar, return: current_revpar, previous_revpar, percentage_change, change_direction
   - This allows users to see the actual values, not just the percentage change
   - Use conditional aggregation (CASE WHEN in SUM) to calculate both periods in a single query, NOT CROSS JOIN!
   - Example CORRECT query structure:
     SELECT 
       CASE WHEN SUM(curr_physical) = 0 THEN 0 ELSE SUM(curr_revenue) / SUM(curr_physical) END AS current_revpar,
       CASE WHEN SUM(prev_physical) = 0 THEN 0 ELSE SUM(prev_revenue) / SUM(prev_physical) END AS previous_revpar,
       CASE WHEN SUM(prev_physical) = 0 THEN NULL ELSE ((SUM(curr_revenue)/SUM(curr_physical) - SUM(prev_revenue)/SUM(prev_physical)) / (SUM(prev_revenue)/SUM(prev_physical))) * 100 END AS percentage_change,
       CASE WHEN SUM(curr_revenue)/SUM(curr_physical) > SUM(prev_revenue)/SUM(prev_physical) THEN 'INCREASE' ELSE 'DECREASE' END AS change_direction
     FROM (
       SELECT
         SUM(CASE WHEN BUSINESS_DATE >= current_start AND BUSINESS_DATE < current_end THEN ROOM_REVENUE ELSE 0 END) AS curr_revenue,
         SUM(CASE WHEN BUSINESS_DATE >= current_start AND BUSINESS_DATE < current_end THEN PHYSICAL_ROOM ELSE 0 END) AS curr_physical,
         SUM(CASE WHEN BUSINESS_DATE >= previous_start AND BUSINESS_DATE < previous_end THEN ROOM_REVENUE ELSE 0 END) AS prev_revenue,
         SUM(CASE WHEN BUSINESS_DATE >= previous_start AND BUSINESS_DATE < previous_end THEN PHYSICAL_ROOM ELSE 0 END) AS prev_physical
       FROM REP_MANAGER
       WHERE BUSINESS_DATE >= previous_start AND BUSINESS_DATE < current_end
     )
   - NEVER use CROSS JOIN for period comparisons - it's inefficient and wrong!

Table Context:
""" + repManagerContext
                },
                {
                    "role": "user",
                    "content": f"""User Request: {clarifiedDescription}

Generate a valid Oracle SQL SELECT query/queries for this request. 
If this request requires {sqlCount} SQL query/queries, generate all of them."""
                }
            ]
            
            # Call LLM to generate SQL
            llmResponse = self.chatService.chat(sqlGenerationMessages, temperature=0.1)
            
            # Extract SQL from response (remove markdown code blocks if present)
            generatedSql = llmResponse.strip()
            if generatedSql.startswith("```sql"):
                generatedSql = generatedSql[7:]
            elif generatedSql.startswith("```"):
                generatedSql = generatedSql[3:]
            if generatedSql.endswith("```"):
                generatedSql = generatedSql[:-3]
            generatedSql = generatedSql.strip()
            
            # Split multiple SQL queries if separated by semicolons
            sqlQueries = [q.strip() for q in generatedSql.split(';') if q.strip()]
            if not sqlQueries:
                sqlQueries = [generatedSql]
            
            # 3. Execute SQL queries with retry mechanism
            sqlExecutionResults = []
            for sqlQuery in sqlQueries:
                executionResult = self._executeSqlWithRetry(
                    sqlQuery, 
                    companyId, 
                    db, 
                    sqlGenerationMessages,
                    maxRetries=3
                )
                sqlExecutionResults.append(executionResult)
            
            return {
                "request_index": requestIndex,
                "request_description": requestDescription,
                "clarified_description": clarifiedDescription,
                "sql_queries": sqlQueries,
                "sql_count": len(sqlQueries),
                "sql_execution_results": sqlExecutionResults,
                "error": None
            }
            
        except Exception as e:
            return {
                "request_index": requestIndex,
                "request_description": requestDescription,
                "clarified_description": None,
                "sql_queries": [],
                "sql_count": 0,
                "sql_execution_results": [],
                "error": str(e)
            }

    def _getCompanyDatabaseConnection(self, companyId: int, db: Session) -> Optional[Dict[str, Any]]:
        """Get company database connection information"""
        try:
            company = db.query(Company).filter(Company.id == companyId).first()
            if not company or not company.database_type:
                return None
            
            return {
                "database_type": company.database_type,
                "host": company.db_ip,
                "port": company.db_port,
                "username": company.db_username,
                "password": company.db_password,
                "service_name": company.db_service_name,
                "sid": company.db_sid,
                "schema": company.db_schema
            }
        except Exception as e:
            return None

    def _executeOracleSql(self, sql: str, connectionInfo: Dict[str, Any]) -> Tuple[bool, Any, Optional[str]]:
        """Execute Oracle SQL query and return results"""
        connection = None
        lastError = None

        # Oracle oracledb does not accept trailing semicolon (causes ORA-00933).
        # DataGrip and other tools allow it; strip it before execute.
        sql = sql.strip()
        if sql.endswith(";"):
            sql = sql[:-1].rstrip()

        try:
            # Try multiple connection methods in order of preference
            
            # Method 1: Try with service_name if provided
            if connectionInfo.get("service_name"):
                try:
                    params = oracledb.ConnectParams(
                        host=connectionInfo["host"],
                        port=int(connectionInfo["port"]),
                        user=connectionInfo["username"],
                        password=connectionInfo["password"],
                        service_name=connectionInfo["service_name"]
                    )
                    connection = oracledb.connect(params=params)
                except Exception as e1:
                    lastError = str(e1)
                    # Continue to next method
            
            # Method 2: Try with SID if provided (and Method 1 didn't work)
            if connection is None and connectionInfo.get("sid"):
                try:
                    params = oracledb.ConnectParams(
                        host=connectionInfo["host"],
                        port=int(connectionInfo["port"]),
                        user=connectionInfo["username"],
                        password=connectionInfo["password"],
                        sid=connectionInfo["sid"]
                    )
                    connection = oracledb.connect(params=params)
                except Exception as e2:
                    lastError = str(e2)
                    # Continue to next method
            
            # Method 3: Try connection string format (username/password@host:port)
            if connection is None:
                try:
                    connectionString = f"{connectionInfo['username']}/{connectionInfo['password']}@{connectionInfo['host']}:{connectionInfo['port']}"
                    connection = oracledb.connect(connectionString)
                except Exception as e3:
                    lastError = str(e3)
                    # Continue to next method
            
            # Method 4: Try with makedsn
            if connection is None:
                try:
                    dsn = oracledb.makedsn(
                        host=connectionInfo["host"],
                        port=int(connectionInfo["port"])
                    )
                    connection = oracledb.connect(
                        user=connectionInfo["username"],
                        password=connectionInfo["password"],
                        dsn=dsn
                    )
                except Exception as e4:
                    lastError = str(e4)
                    # Continue to next method
            
            # Method 5: Try with ConnectParams without service_name/SID
            if connection is None:
                try:
                    params = oracledb.ConnectParams(
                        host=connectionInfo["host"],
                        port=int(connectionInfo["port"]),
                        user=connectionInfo["username"],
                        password=connectionInfo["password"]
                    )
                    connection = oracledb.connect(params=params)
                except Exception as e5:
                    lastError = str(e5)
            
            # If all methods failed, raise error
            if connection is None:
                errorMsg = f"Failed to connect to Oracle database at {connectionInfo['host']}:{connectionInfo['port']}. Last error: {lastError}"
                if not connectionInfo.get("service_name") and not connectionInfo.get("sid"):
                    errorMsg += " Note: Oracle connection typically requires service_name or SID. Please configure db_service_name or db_sid in company settings."
                return False, None, errorMsg
            
            # Execute SQL query
            cursor = connection.cursor()
            try:
                cursor.execute(sql)
                
                # Fetch results
                columns = [desc[0] for desc in cursor.description] if cursor.description else []
                rows = cursor.fetchall()
                
                # Convert to list of dictionaries
                results = []
                for row in rows:
                    rowDict = {}
                    for i, col in enumerate(columns):
                        value = row[i]
                        # Convert datetime and other types to string for JSON serialization
                        if isinstance(value, datetime):
                            value = value.isoformat()
                        elif hasattr(value, '__str__'):
                            value = str(value)
                        rowDict[col] = value
                    results.append(rowDict)
                
                return True, results, None
                
            finally:
                cursor.close()
                connection.close()
                
        except Exception as e:
            errorMsg = str(e)
            # Add more context to connection errors
            if "cannot connect" in errorMsg.lower() or "listener" in errorMsg.lower():
                if not connectionInfo.get("service_name") and not connectionInfo.get("sid"):
                    errorMsg += " (Oracle connection may require service_name or SID. Please configure db_service_name or db_sid in company settings.)"
            return False, None, errorMsg

    def _executeSqlWithRetry(
        self, 
        sql: str, 
        companyId: int, 
        db: Session, 
        sqlGenerationMessages: List[Dict[str, str]],
        maxRetries: int = 3
    ) -> Dict[str, Any]:
        """Execute SQL with retry mechanism and LLM-based error correction"""
        connectionInfo = self._getCompanyDatabaseConnection(companyId, db)
        
        if not connectionInfo:
            return {
                "sql": sql,
                "success": False,
                "data": None,
                "error": "Company database connection information not found",
                "retry_count": 0
            }
        
        lastError = None
        currentSql = sql
        
        for attempt in range(1, maxRetries + 1):
            try:
                # Try to execute SQL
                if connectionInfo["database_type"].lower() == "oracle":
                    success, data, error = self._executeOracleSql(currentSql, connectionInfo)
                else:
                    return {
                        "sql": currentSql,
                        "success": False,
                        "data": None,
                        "error": f"Unsupported database type: {connectionInfo['database_type']}",
                        "retry_count": attempt - 1
                    }
                
                if success:
                    return {
                        "sql": currentSql,
                        "success": True,
                        "data": data,
                        "error": None,
                        "retry_count": attempt - 1
                    }
                else:
                    lastError = error
                    # If this is not the last attempt, try to fix SQL with LLM
                    if attempt < maxRetries:
                        hint = ""
                        if "ORA-00904" in (error or "") or "invalid identifier" in (error or "").lower():
                            hint = "\n\nUse ONLY columns from the allowed list. Revenue by channel: COMPANY_ROOM_REVENUE, AGENT_ROOM_REVENUE, GROUP_ROOM_REV (not GROUP_ROOM_REVENUE or SOURCE_ROOM_REVENUE). No COMMISSION column. Segment → GROUP_ROOM, COMPANY_ROOMS_TODAY, FIT_GUESTS_TODAY. İptal → CANCEL_ROOMS, CANCEL_RESERVATION. LOS → FIT_MEM_LOS_NIGHTS, FIT_MEM_LOS_RESV. Otelde kaç kişi / in-house guests → GUEST_IN_HOUSE only; do NOT use INHOUSE_PRS, OCC_PRS – they do not exist. Do not use SOURCE_ROOM_REVENUE, GROUP_ROOM_REVENUE, COMMISSION, ROOM_CATEGORY, NO_SHOW_ROOMS, COMP_PHYSICAL_ROOM, MARKET_GROUP, CHANNEL – not in list. Physical rooms → PHYSICAL_ROOM. No-show → NOSHOW_ROOM, NOSHOW_PRS, NOSHOW_RESERVATIONS. Room category → ROOM_CLASS."
                        elif "ORA-00933" in (error or "") or "command not properly ended" in (error or "").lower():
                            hint = "\n\nCheck SQL syntax: Every non-aggregated column in SELECT must be in GROUP BY; no trailing comma before FROM/WHERE); subquery must have closing paren. Use ONLY allowed column names. Room category → ROOM_CLASS."
                        elif "ORA-01722" in (error or "") or "invalid number" in (error or "").lower():
                            hint = "\n\nInvalid number: ensure numeric columns are not compared to strings; use TO_NUMBER or correct type. Use ONLY allowed column names."
                        elif "ORA-12702" in (error or "") or "NLS" in (error or ""):
                            hint = "\n\nInvalid NLS: use simple TO_CHAR with format 'YYYY-MM-DD' or 'DD-MON-YYYY'; avoid NLS_DATE_LANGUAGE or complex NLS parameters. Use ONLY allowed column names."
                        fixMessages = sqlGenerationMessages.copy()
                        fixMessages.append({
                            "role": "user",
                            "content": f"""The previous SQL query failed with an error. Please fix it.

Previous SQL:
{currentSql}

Error:
{error}
{hint}

Please generate a corrected SQL query that fixes this error."""
                        })
                        
                        llmResponse = self.chatService.chat(fixMessages, temperature=0.1)
                        
                        # Extract fixed SQL
                        fixedSql = llmResponse.strip()
                        if fixedSql.startswith("```sql"):
                            fixedSql = fixedSql[7:]
                        elif fixedSql.startswith("```"):
                            fixedSql = fixedSql[3:]
                        if fixedSql.endswith("```"):
                            fixedSql = fixedSql[:-3]
                        fixedSql = fixedSql.strip()
                        
                        currentSql = fixedSql
                        # Update the last user message in sqlGenerationMessages for next retry
                        sqlGenerationMessages = fixMessages
                    else:
                        # Last attempt failed
                        return {
                            "sql": currentSql,
                            "success": False,
                            "data": None,
                            "error": f"Failed after {maxRetries} attempts. Last error: {error}",
                            "retry_count": maxRetries
                        }
                        
            except Exception as e:
                lastError = str(e)
                if attempt < maxRetries:
                    # Connection error - try to fix
                    fixMessages = sqlGenerationMessages.copy()
                    fixMessages.append({
                        "role": "user",
                        "content": f"""The SQL execution failed with a connection or execution error. Please review and fix the SQL if needed.

SQL:
{currentSql}

Error:
{str(e)}

Please generate a corrected SQL query if the error is SQL-related, or confirm the SQL is correct if it's a connection issue."""
                    })
                    
                    llmResponse = self.chatService.chat(fixMessages, temperature=0.1)
                    fixedSql = llmResponse.strip()
                    if fixedSql.startswith("```sql"):
                        fixedSql = fixedSql[7:]
                    elif fixedSql.startswith("```"):
                        fixedSql = fixedSql[3:]
                    if fixedSql.endswith("```"):
                        fixedSql = fixedSql[:-3]
                    fixedSql = fixedSql.strip()
                    
                    currentSql = fixedSql
                    sqlGenerationMessages = fixMessages
                else:
                    return {
                        "sql": currentSql,
                        "success": False,
                        "data": None,
                        "error": f"Failed after {maxRetries} attempts. Last error: {str(e)}",
                        "retry_count": maxRetries
                    }
        
        return {
            "sql": currentSql,
            "success": False,
            "data": None,
            "error": f"Failed after {maxRetries} attempts. Last error: {lastError}",
            "retry_count": maxRetries
        }

    async def run_sql_pipeline(
        self, company_id: int, user_question: str, db: Session
    ) -> Dict[str, Any]:
        """
        Run full RAG + prompt analysis + SQL generation + execution pipeline.
        Returns dict with prompt_analysis, sql_results, date_info (no HTTP wrapper).
        Used by both testOpenai and ChatController.streamMessage.
        """
        dateInfo = self._getCurrentDateInfo()
        repManagerContext = self._loadREPManagerContext(company_id)

        userMessageWithDateInfo = f"""{user_question}

IMPORTANT: Today's date is {dateInfo['today']} ({dateInfo['current_year']}-{dateInfo['current_month']:02d}-{dateInfo['current_day']:02d}).
When analyzing dates, use this actual date. If user mentions "bugün", "today", "bu ay", "this month", "bu yıl", "this year",
clarify the exact date range in the request description."""

        promptAnalysis = self.ragService.analyzeUserPrompt(
            userMessageWithDateInfo, companyId=company_id, db=db
        )
        sqlResults = []

        if promptAnalysis.get("request_count", 0) > 0:
            requests = promptAnalysis.get("requests", [])
            for idx, req in enumerate(requests, 1):
                if "index" not in req:
                    req["index"] = idx

            loop = asyncio.get_event_loop()

            def processWithNewSession(req):
                threadDb = SessionLocal()
                try:
                    return self._processSingleRequestSync(
                        req, company_id, threadDb, repManagerContext, dateInfo
                    )
                finally:
                    threadDb.close()

            with ThreadPoolExecutor(max_workers=min(len(requests), 10)) as executor:
                futures = [
                    loop.run_in_executor(executor, processWithNewSession, req)
                    for req in requests
                ]
                results = await asyncio.gather(*futures, return_exceptions=True)

                for result in results:
                    if isinstance(result, Exception):
                        sqlResults.append({
                            "request_index": 0,
                            "request_description": "",
                            "clarified_description": None,
                            "sql_queries": [],
                            "sql_count": 0,
                            "sql_execution_results": [],
                            "error": f"Processing failed: {str(result)}"
                        })
                    else:
                        sqlResults.append(result)

        sqlResults.sort(key=lambda x: x.get("request_index", 0))
        return {
            "prompt_analysis": promptAnalysis,
            "sql_results": sqlResults,
            "date_info": {
                "today": dateInfo["today"],
                "current_year": dateInfo["current_year"],
                "current_month": dateInfo["current_month"]
            }
        }

    async def testOpenai(self, request: TestOpenAIRequest, currentUser: dict, db: Session = Depends(get_db)):
        companyId = currentUser.get("company_id", 1)
        resultData = await self.run_sql_pipeline(companyId, request.message, db)
        return httpResponse(
            isSuccess=True,
            message="Analysis and SQL generation successful",
            statusCode=200,
            data=resultData
        )
