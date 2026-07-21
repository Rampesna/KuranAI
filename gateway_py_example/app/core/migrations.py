"""
Database migration and seeding system
"""
from sqlalchemy import text, inspect
from sqlalchemy.orm import Session
from datetime import datetime, timedelta
from app.core.database import engine, SessionLocal, Base
from app.core.config import settings
from app.models import (
    Company, Admin, Manager, Employee, ChatSession, ChatMessage,
    SQLGenerationRAGVector, Setting, UserToken, VectorEmbedding,
    Agent, AgentReport
)
from app.services.common.password import PasswordService
import os
import subprocess
from pathlib import Path


def drop_all_tables():
    """Drop all existing tables"""
    print("🗑️  Dropping all existing tables...")
    with engine.connect() as conn:
        # Disable foreign key checks temporarily
        conn.execute(text("SET session_replication_role = 'replica';"))
        
        # Get all table names
        inspector = inspect(engine)
        tables = inspector.get_table_names()
        
        # Drop tables in reverse dependency order
        for table in reversed(tables):
            try:
                conn.execute(text(f'DROP TABLE IF EXISTS "{table}" CASCADE;'))
                print(f"   ✓ Dropped table: {table}")
            except Exception as e:
                print(f"   ⚠️  Error dropping table {table}: {e}")
        
        # Re-enable foreign key checks
        conn.execute(text("SET session_replication_role = 'origin';"))
        conn.commit()
    
    print("✅ All tables dropped")


def run_migrations():
    """Create all tables if they don't exist"""
    print("🔄 Running database migrations...")
    
    # Create all tables
    Base.metadata.create_all(bind=engine)
    print("✅ Database migrations completed")


def import_sql_file(file_path: str, description: str):
    """Import SQL file into database using psql for better compatibility"""
    sql_file = Path(file_path)
    
    if not sql_file.exists():
        print(f"⚠️  SQL file not found: {file_path}, skipping...")
        return
    
    print(f"📥 Importing {description} from {sql_file.name}...")
    
    try:
        # Use psql for importing large SQL files (better for complex INSERT statements)
        
        # Extract connection info from DATABASE_URL
        db_url = settings.DATABASE_URL
        # Format: postgresql://user:password@host:port/dbname
        import re
        match = re.match(r'postgresql://([^:]+):([^@]+)@([^:]+):(\d+)/(.+)', db_url)
        if match:
            db_user, db_pass, db_host, db_port, db_name = match.groups()
            
            # Set PGPASSWORD environment variable
            env = os.environ.copy()
            env['PGPASSWORD'] = db_pass
            
            # Use psql to import the file
            cmd = [
                'psql',
                '-h', db_host,
                '-p', db_port,
                '-U', db_user,
                '-d', db_name,
                '-f', str(sql_file),
                '-q'  # Quiet mode
            ]
            
            result = subprocess.run(
                cmd,
                env=env,
                capture_output=True,
                text=True,
                timeout=300  # 5 minute timeout
            )
            
            if result.returncode == 0:
                print(f"✅ {description} imported successfully using psql")
            else:
                # If psql fails, try with Python SQLAlchemy as fallback
                print(f"   ⚠️  psql import failed, trying Python fallback...")
                print(f"   ⚠️  psql error: {result.stderr[:200]}")
                
                # Fallback: try executing with SQLAlchemy
                with open(sql_file, 'r', encoding='utf-8') as f:
                    sql_content = f.read()
                
                with engine.connect() as conn:
                    trans = conn.begin()
                    try:
                        conn.execute(text(sql_content))
                        trans.commit()
                        print(f"✅ {description} imported successfully (Python fallback)")
                    except Exception as e:
                        trans.rollback()
                        error_str = str(e).lower()
                        if "syntax error" in error_str:
                            print(f"   ⚠️  SQL syntax error (file may need manual review): {str(e)[:200]}")
                            print(f"   ⚠️  Skipping import due to syntax errors")
                        else:
                            raise
        else:
            raise ValueError("Could not parse DATABASE_URL")
            
    except subprocess.TimeoutExpired:
        print(f"❌ Timeout importing {description} (file too large)")
    except FileNotFoundError:
        # psql not available, use Python fallback
        print(f"   ⚠️  psql not found, using Python fallback...")
        with open(sql_file, 'r', encoding='utf-8') as f:
            sql_content = f.read()
        
        with engine.connect() as conn:
            trans = conn.begin()
            try:
                conn.execute(text(sql_content))
                trans.commit()
                print(f"✅ {description} imported successfully (Python fallback)")
            except Exception as e:
                trans.rollback()
                error_str = str(e).lower()
                if "syntax error" in error_str:
                    print(f"   ⚠️  SQL syntax error: {str(e)[:200]}")
                    print(f"   ⚠️  Skipping import due to syntax errors")
                else:
                    raise
    except Exception as e:
        print(f"❌ Error importing {description}: {e}")
        raise


def seed_database():
    """Seed initial data"""
    print("🌱 Seeding database...")
    
    db = SessionLocal()
    try:
        # Check if data already exists
        existing_company = db.query(Company).filter(Company.slug == "intelor-demo").first()
        if existing_company:
            company_id = existing_company.id
            # Check if agents already exist
            existing_agents = db.query(Agent).filter(Agent.company_id == company_id).count()
            if existing_agents > 0:
                print("⚠️  Seed data already exists, skipping...")
                return
            # Company exists but no agents, so seed agents only
            print("⚠️  Company exists but agents missing, seeding agents only...")
        else:
            company_id = None
        
        password_hash = PasswordService.hash("admin123")
        
        # 1. Create demo company (if not exists)
        if company_id is None:
            company = Company(
                name="Intelor Demo Hotel",
                slug="intelor-demo",
                description="Demo hotel for Intelor AI Analytics Chat system",
                is_active=True,
                database_type="Oracle",
                db_ip="95.70.179.78",
                db_port="1521",
                db_username="PROISV",
                db_password="PROISV",
                db_service_name=os.getenv("ORACLE_SERVICE_NAME") or None,
                db_sid=os.getenv("ORACLE_SID") or "OPERA",
                db_schema="OPERA",  # Schema prefix for tables (e.g., OPERA.REP_MANAGER)
                settings={
                    "timezone": "Europe/Istanbul",
                    "language": "tr",
                    "features": {"chat": True, "analytics": True}
                }
            )
            db.add(company)
            db.flush()  # Get the company ID
            company_id = company.id
            company_name = company.name
        else:
            company = existing_company
            company_name = company.name
            # Update database connection info if not set
            if not company.database_type:
                company.database_type = "Oracle"
                company.db_ip = "95.70.179.78"
                company.db_port = "1521"
                company.db_username = "PROISV"
                company.db_password = "PROISV"
                company.db_service_name = os.getenv("ORACLE_SERVICE_NAME") or company.db_service_name
                company.db_sid = os.getenv("ORACLE_SID") or "OPERA"
                company.db_schema = "OPERA"
                db.flush()
            # Always fix db_sid if null for Oracle demo company (e.g. after earlier migrate)
            elif company.database_type == "Oracle" and not company.db_sid:
                company.db_sid = os.getenv("ORACLE_SID") or "OPERA"
                db.flush()
        
        # 2. Create admin user (if not exists)
        if not db.query(Admin).filter(Admin.email == "admin@intelor.com").first():
            admin = Admin(
                company_id=company_id,
                email="admin@intelor.com",
                username="admin",
                password_hash=password_hash,
                first_name="System",
                last_name="Administrator",
                is_active=True,
                is_super_admin=True
            )
            db.add(admin)
        else:
            admin = db.query(Admin).filter(Admin.email == "admin@intelor.com").first()
        
        # 3. Create manager user (if not exists)
        if not db.query(Manager).filter(Manager.email == "manager@intelor.com").first():
            manager = Manager(
                company_id=company_id,
                email="manager@intelor.com",
                username="manager",
                password_hash=password_hash,
                first_name="General",
                last_name="Manager",
                is_active=True
            )
            db.add(manager)
        else:
            manager = db.query(Manager).filter(Manager.email == "manager@intelor.com").first()
        
        # 4. Create employee user (if not exists)
        if not db.query(Employee).filter(Employee.email == "employee@intelor.com").first():
            employee = Employee(
                company_id=company_id,
                email="employee@intelor.com",
                username="employee",
                password_hash=password_hash,
                first_name="Demo",
                last_name="Employee",
                is_active=True
            )
            db.add(employee)
        else:
            employee = db.query(Employee).filter(Employee.email == "employee@intelor.com").first()
        
        # 5. Create settings entry for prompt analysis (if not exists)
        setting = db.query(Setting).filter(Setting.key == "prompt_analysis_system_prompt").first()
        if not setting:
            setting = Setting(
            company_id=company_id,
            key="prompt_analysis_system_prompt",
            value="""You are a SQL query analysis expert. Analyze the user's prompt and determine:
1. Does the user want ONE thing or MULTIPLE things?
2. For each thing they want, how many SQL queries are needed to get the answer?

IMPORTANT: Return ONLY valid JSON in this exact format:
{
  "is_single_request": true/false,
  "request_count": number,
  "requests": [
    {
      "index": 1,
      "description": "What the user wants for this request",
      "sql_count": number,
      "sql_descriptions": ["First SQL does...", "Second SQL does..."]
    }
  ],
  "requires_multiple_sql": true/false
}

Rules:
- If user asks "Show me revenue and occupancy", that's 2 requests (revenue + occupancy)
- If user asks "Show me this month's revenue", that's 1 request but might need 1 SQL
- If user asks "Compare this month vs last month revenue", that's 1 request but needs 2 SQLs (current month + last month)
- Be precise and return ONLY JSON, no explanations.""",
            description="System prompt used by the RAG service to analyze user prompts and determine the number of SQL queries needed. This prompt instructs the LLM to break down user requests into structured JSON format.",
                category="rag",
                is_active=True
            )
            db.add(setting)
        
        # 6. Create agents with reports
        # Agent 1: Günlük Gelir Raporu
        agent1 = Agent(
            company_id=company_id,
            name="Günlük Gelir Raporu",
            description="Her gün sabah 08:00'de toplam gelir raporunu hazırlar",
            schedule_type="daily",
            schedule_time="08:00",
            schedule="Her gün 08:00",
            status="active",
            last_run_at=datetime.now() - timedelta(days=1) + timedelta(hours=8),
            next_run_at=datetime.now() + timedelta(hours=8)
        )
        db.add(agent1)
        db.flush()
        
        for report_name in ["Toplam Gelir", "Oda Bazlı Gelir", "Hizmet Gelirleri"]:
            db.add(AgentReport(agent_id=agent1.id, report_name=report_name))
        
        # Agent 2: Doluluk Oranı Analizi
        agent2 = Agent(
            company_id=company_id,
            name="Doluluk Oranı Analizi",
            description="Her gün akşam 18:00'de doluluk oranı ve tahminleri gönderir",
            schedule_type="daily",
            schedule_time="18:00",
            schedule="Her gün 18:00",
            status="active",
            last_run_at=datetime.utcnow() - timedelta(days=1) + timedelta(hours=18),
            next_run_at=datetime.utcnow() + timedelta(hours=18)
        )
        db.add(agent2)
        db.flush()
        
        for report_name in ["Doluluk Oranı", "Gelecek Hafta Tahmini", "Yoğun Günler"]:
            db.add(AgentReport(agent_id=agent2.id, report_name=report_name))
        
        # Agent 3: Haftalık Performans Özeti
        agent3 = Agent(
            company_id=company_id,
            name="Haftalık Performans Özeti",
            description="Her Pazartesi sabah 09:00'de haftalık performans raporu hazırlar",
            schedule_type="weekly",
            schedule_time="09:00",
            schedule_day="Pazartesi",
            schedule="Her Pazartesi 09:00",
            status="paused",
            last_run_at=datetime.utcnow() - timedelta(days=7) + timedelta(hours=9),
            next_run_at=datetime.utcnow() + timedelta(days=7) + timedelta(hours=9)
        )
        db.add(agent3)
        db.flush()
        
        for report_name in ["Haftalık Gelir", "Ortalama ADR", "RevPAR", "Misafir Memnuniyeti"]:
            db.add(AgentReport(agent_id=agent3.id, report_name=report_name))
        
        # Agent 4: Aylık Personel Raporu
        agent4 = Agent(
            company_id=company_id,
            name="Aylık Personel Raporu",
            description="Her ayın 1'inde personel performans raporunu hazırlar",
            schedule_type="monthly",
            schedule_time="10:00",
            schedule_date=1,
            schedule="Her ayın 1'i 10:00",
            status="active",
            last_run_at=datetime.utcnow() - timedelta(days=30) + timedelta(hours=10),
            next_run_at=datetime.utcnow() + timedelta(days=30) - timedelta(days=1) + timedelta(hours=10)
        )
        db.add(agent4)
        db.flush()
        
        for report_name in ["Personel Performansı", "Çalışma Saatleri", "Başarı Metrikleri"]:
            db.add(AgentReport(agent_id=agent4.id, report_name=report_name))
        
        db.commit()
        print("✅ Database seeding completed")
        print(f"   - Company: {company.name} (ID: {company_id})")
        print(f"   - Admin: {admin.email}")
        print(f"   - Manager: {manager.email}")
        print(f"   - Employee: {employee.email}")
        print(f"   - Setting: {setting.key}")
        print(f"   - Agents: 4 agents created")
        
    except Exception as e:
        db.rollback()
        print(f"❌ Error seeding database: {e}")
        raise
    finally:
        db.close()


def run_full_migration():
    """Run full migration: drop all tables, recreate, seed, and import SQL files"""
    print("🚀 Starting full database migration...")
    print("=" * 60)
    
    # Step 1: Drop all tables
    drop_all_tables()
    print()
    
    # Step 2: Create all tables
    run_migrations()
    print()
    
    # Step 3: Seed initial data
    seed_database()
    print()
    
    # Step 4: Import SQL files (if they exist)
    # Try multiple possible locations (Docker container and local development)
    possible_paths = [
        Path(__file__).parent.parent.parent,  # gateway_py/ (Docker: /app/)
        Path(__file__).parent.parent.parent.parent,  # Project root (local dev)
        Path.cwd(),  # Current working directory
        Path.cwd().parent,  # Parent of current directory
        Path("/app"),  # Docker container app directory
    ]
    
    sql_rag_vectors_path = None
    vector_embeddings_path = None
    
    for base_path in possible_paths:
        if (base_path / "sql_generation_rag_vectors.sql").exists():
            sql_rag_vectors_path = base_path / "sql_generation_rag_vectors.sql"
        if (base_path / "vector_embeddings.sql").exists():
            vector_embeddings_path = base_path / "vector_embeddings.sql"
        if sql_rag_vectors_path and vector_embeddings_path:
            break
    
    # Import sql_generation_rag_vectors.sql
    if sql_rag_vectors_path:
        import_sql_file(str(sql_rag_vectors_path), "SQL Generation RAG Vectors")
    else:
        print("⚠️  sql_generation_rag_vectors.sql not found, skipping...")
    print()
    
    # Import vector_embeddings.sql
    if vector_embeddings_path:
        import_sql_file(str(vector_embeddings_path), "Vector Embeddings")
    else:
        print("⚠️  vector_embeddings.sql not found, skipping...")
    print()
    
    print("=" * 60)
    print("✅ Full migration completed successfully!")


if __name__ == "__main__":
    import sys
    
    if len(sys.argv) > 1 and sys.argv[1] == "full":
        # Run full migration (drop all, recreate, seed, import SQL)
        run_full_migration()
    else:
        # Run normal migration (create tables if not exist, seed if not exists)
        print("Running migrations and seeding...")
        run_migrations()
        seed_database()
        print("✅ All done!")
