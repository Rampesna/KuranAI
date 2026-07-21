#!/usr/bin/env python3
"""
Update imports in all Python files to match new structure
"""
import re
from pathlib import Path

# Base path
BASE_PATH = Path(__file__).parent / "app_new"

# Import mappings (old -> new)
IMPORT_MAPPINGS = {
    # Core
    r'from app\.core\.': 'from app.Core.',
    r'import app\.core\.': 'import app.Core.',

    # Enums (formerly types)
    r'from app\.types\.': 'from app.Enums.',
    r'import app\.types\.': 'import app.Enums.',

    # DTOs
    r'from app\.dto\.services\.': 'from app.DataTransferObjects.Sqlalchemy.',
    r'import app\.dto\.services\.': 'import app.DataTransferObjects.Sqlalchemy.',

    # Http Controllers
    r'from app\.http\.controllers\.': 'from app.Http.Controllers.',
    r'import app\.http\.controllers\.': 'import app.Http.Controllers.',

    # Http Middleware
    r'from app\.http\.middlewares\.': 'from app.Http.Middleware.',
    r'import app\.http\.middlewares\.': 'import app.Http.Middleware.',

    # Http Requests
    r'from app\.http\.requests\.': 'from app.Http.Requests.',
    r'import app\.http\.requests\.': 'import app.Http.Requests.',

    # Services
    r'from app\.services\.sqlalchemy\.': 'from app.Services.Sqlalchemy.',
    r'import app\.services\.sqlalchemy\.': 'import app.Services.Sqlalchemy.',
    r'from app\.services\.common\.': 'from app.Services.Common.',
    r'import app\.services\.common\.': 'import app.Services.Common.',
    r'from app\.services\.openai\.': 'from app.Services.OpenAI.',
    r'import app\.services\.openai\.': 'import app.Services.OpenAI.',
    r'from app\.services\.intelor_services\.': 'from app.Services.IntelorServices.',
    r'import app\.services\.intelor_services\.': 'import app.Services.IntelorServices.',

    # Models
    r'from app\.models\.': 'from app.Models.',
    r'import app\.models\.': 'import app.Models.',

    # Utils
    r'from app\.utils\.': 'from app.Utils.',
    r'import app\.utils\.': 'import app.Utils.',

    # Routers
    r'from app\.routers\.': 'from app.Routers.',
    r'import app\.routers\.': 'import app.Routers.',
}

def update_file_imports(file_path):
    """Update imports in a single file"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Apply all import mappings
        for old_pattern, new_pattern in IMPORT_MAPPINGS.items():
            content = re.sub(old_pattern, new_pattern, content)

        # Only write if content changed
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        return False
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def process_directory(directory):
    """Process all Python files in directory recursively"""
    updated_count = 0
    total_count = 0

    for py_file in directory.rglob('*.py'):
        if '__pycache__' in str(py_file):
            continue

        total_count += 1
        if update_file_imports(py_file):
            updated_count += 1
            print(f"Updated: {py_file.relative_to(BASE_PATH)}")

    return updated_count, total_count

def main():
    print("Starting import updates...")
    print(f"Base path: {BASE_PATH}\n")

    updated, total = process_directory(BASE_PATH)

    print(f"\n{'='*50}")
    print(f"Import updates complete!")
    print(f"{'='*50}")
    print(f"Updated {updated} out of {total} files")

if __name__ == "__main__":
    main()
