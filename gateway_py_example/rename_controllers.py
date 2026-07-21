#!/usr/bin/env python3
"""
Update controller imports after renaming files to PascalCase
"""
import re
from pathlib import Path

BASE_PATH = Path(__file__).parent / "app_new"

CONTROLLER_RENAMES = {
    # Auth Controllers
    r'\.authController import': '.AuthController import',
    r'/authController\.': '/AuthController.',

    # Chat Controllers
    r'\.chatController import': '.ChatController import',
    r'/chatController\.': '/ChatController.',

    # Agent Controllers
    r'\.agentController import': '.AgentController import',
    r'/agentController\.': '/AgentController.',

    # Test Controllers
    r'\.testController import': '.TestController import',
    r'/testController\.': '/TestController.',
}

def update_file(file_path):
    """Update controller references in a file"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        for old_pattern, new_pattern in CONTROLLER_RENAMES.items():
            content = re.sub(old_pattern, new_pattern, content)

        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        return False
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def main():
    print("Updating controller imports...")

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
