#!/usr/bin/env python3
"""
Update http_response usages to use new Traits.HttpResponse
"""
import re
from pathlib import Path

BASE_PATH = Path(__file__).parent / "app_new"

def update_file(file_path):
    """Update http_response references in a file"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Update imports
        content = re.sub(
            r'from app\.Utils\.http_response import http_response',
            'from app.Traits.HttpResponse import httpResponse',
            content
        )

        # Update function calls - convert parameters to camelCase
        # This is a simple replacement - may need manual review
        content = re.sub(r'\bhttp_response\(', 'httpResponse(', content)
        content = re.sub(r'\bis_success=', 'isSuccess=', content)
        content = re.sub(r'\bstatus_code=', 'statusCode=', content)

        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        return False
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def main():
    print("Updating http_response to httpResponse...")

    updated_count = 0
    for py_file in BASE_PATH.rglob('*.py'):
        if '__pycache__' in str(py_file):
            continue

        if update_file(py_file):
            updated_count += 1
            print(f"Updated: {py_file.relative_to(BASE_PATH)}")

    print(f"\nUpdated {updated_count} files")

if __name__ == "__main__":
    main()
