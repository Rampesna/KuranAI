#!/usr/bin/env python3
"""
Update ServiceResponse method calls to camelCase
"""
import re
from pathlib import Path

BASE_PATH = Path(__file__).parent / "app_new"

def update_file(file_path):
    """Update ServiceResponse method calls in a file"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Update ServiceResponse constructor parameters
        content = re.sub(r'\bis_success=', 'isSuccess=', content)
        content = re.sub(r'\bstatus_code=', 'statusCode=', content)

        # Update method calls
        content = re.sub(r'\.is_success\(\)', '.isSuccess()', content)
        content = re.sub(r'\.get_message\(\)', '.getMessage()', content)
        content = re.sub(r'\.get_status_code\(\)', '.getStatusCode()', content)
        content = re.sub(r'\.get_data\(\)', '.getData()', content)
        content = re.sub(r'\.to_dict\(\)', '.toDict()', content)

        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        return False
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def main():
    print("Updating ServiceResponse method calls...")

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
