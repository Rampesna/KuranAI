#!/usr/bin/env python3
"""
Refactor script to reorganize gateway_py project structure to match Laravel architecture
"""
import os
import shutil
from pathlib import Path

# Base path
BASE_PATH = Path(__file__).parent / "app"

# Mapping old paths to new paths
RESTRUCTURE_MAP = {
    # Core layer
    "core": "Core",

    # Enums from types
    "types": "Enums",

    # DTOs
    "dto/services": "DataTransferObjects/Sqlalchemy",

    # Http layer
    "http/controllers": "Http/Controllers",
    "http/middlewares": "Http/Middleware",
    "http/requests": "Http/Requests",

    # Services
    "services/sqlalchemy": "Services/Sqlalchemy",
    "services/common": "Services/Common",
    "services/openai": "Services/OpenAI",
    "services/intelor_services": "Services/IntelorServices",

    # Models stays same (will just capitalize)
    "models": "Models",

    # Utils stays same (will just capitalize)
    "utils": "Utils",
}

def create_init_file(directory):
    """Create __init__.py in directory if it doesn't exist"""
    init_file = directory / "__init__.py"
    if not init_file.exists():
        init_file.write_text("")
        print(f"Created {init_file}")

def ensure_directory(path):
    """Ensure directory exists"""
    path.mkdir(parents=True, exist_ok=True)
    create_init_file(path)

def copy_with_structure(src, dst):
    """Copy file or directory from src to dst, creating necessary structure"""
    if src.is_file():
        ensure_directory(dst.parent)
        shutil.copy2(src, dst)
        print(f"Copied: {src} -> {dst}")
    elif src.is_dir():
        ensure_directory(dst)
        for item in src.iterdir():
            if item.name == "__pycache__" or item.name == ".DS_Store":
                continue
            new_dst = dst / item.name
            copy_with_structure(item, new_dst)

def main():
    print("Starting refactoring...")
    print(f"Base path: {BASE_PATH}")

    # Create backup directory
    backup_dir = BASE_PATH.parent / "app_backup"
    if backup_dir.exists():
        print(f"Removing old backup...")
        shutil.rmtree(backup_dir)

    print(f"Creating backup at {backup_dir}...")
    shutil.copytree(BASE_PATH, backup_dir)
    print("Backup created!")

    # Create new structure directory
    new_structure_dir = BASE_PATH.parent / "app_new"
    if new_structure_dir.exists():
        print(f"Removing old new_structure...")
        shutil.rmtree(new_structure_dir)

    new_structure_dir.mkdir()
    create_init_file(new_structure_dir)

    # Copy main.py
    main_py = BASE_PATH / "main.py"
    if main_py.exists():
        shutil.copy2(main_py, new_structure_dir / "main.py")
        print(f"Copied main.py")

    # Restructure according to map
    for old_path, new_path in RESTRUCTURE_MAP.items():
        src = BASE_PATH / old_path
        dst = new_structure_dir / new_path

        if src.exists():
            print(f"\nRestructuring: {old_path} -> {new_path}")
            copy_with_structure(src, dst)
        else:
            print(f"Warning: {src} does not exist")

    # Copy routers
    routers_src = BASE_PATH / "routers"
    routers_dst = new_structure_dir / "Routers"
    if routers_src.exists():
        print(f"\nCopying routers...")
        copy_with_structure(routers_src, routers_dst)

    print("\n" + "="*50)
    print("Restructuring complete!")
    print("="*50)
    print(f"\nOld structure backed up to: {backup_dir}")
    print(f"New structure created at: {new_structure_dir}")
    print("\nNext steps:")
    print("1. Review the new structure")
    print("2. Update imports in all files")
    print("3. Replace app/ with app_new/")
    print("4. Test the application")

if __name__ == "__main__":
    main()
