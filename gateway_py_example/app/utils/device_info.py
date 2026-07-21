"""
Utility to extract device information from request headers
"""
from fastapi import Request
from typing import Dict, Any
import re


def extract_device_info(request: Request) -> Dict[str, Any]:
    """
    Extract device information from request headers (Netflix/Google style)
    """
    user_agent = request.headers.get("user-agent", "")
    ip_address = request.headers.get("x-forwarded-for") or request.headers.get("x-real-ip") or request.client.host if request.client else None
    
    device_info: Dict[str, Any] = {
        "user_agent": user_agent,
        "ip_address": ip_address,
    }
    
    # Parse user agent to extract device details
    if user_agent:
        # Platform detection
        if re.search(r'iPhone|iPad|iPod', user_agent, re.IGNORECASE):
            device_info["platform"] = "iOS"
            device_info["device_type"] = "mobile" if "iPhone" in user_agent or "iPod" in user_agent else "tablet"
            # Extract iOS version
            ios_match = re.search(r'OS (\d+[._]\d+)', user_agent)
            if ios_match:
                device_info["os_version"] = ios_match.group(1).replace('_', '.')
            # Extract device name
            device_match = re.search(r'iPhone|iPad|iPod', user_agent)
            if device_match:
                device_info["device_name"] = device_match.group(0)
        elif re.search(r'Android', user_agent, re.IGNORECASE):
            device_info["platform"] = "Android"
            device_info["device_type"] = "mobile"
            # Extract Android version
            android_match = re.search(r'Android (\d+[.\d]*)', user_agent)
            if android_match:
                device_info["os_version"] = android_match.group(1)
            # Extract device name
            device_match = re.search(r';\s*([^;)]+)\s*\)', user_agent)
            if device_match:
                device_info["device_name"] = device_match.group(1).strip()
        elif re.search(r'Windows', user_agent, re.IGNORECASE):
            device_info["platform"] = "Windows"
            device_info["device_type"] = "desktop"
            # Extract Windows version
            windows_match = re.search(r'Windows NT (\d+[.\d]*)', user_agent)
            if windows_match:
                device_info["os_version"] = windows_match.group(1)
        elif re.search(r'Mac OS X|macOS', user_agent, re.IGNORECASE):
            device_info["platform"] = "macOS"
            device_info["device_type"] = "desktop"
            # Extract macOS version
            mac_match = re.search(r'Mac OS X (\d+[._]\d+)', user_agent)
            if mac_match:
                device_info["os_version"] = mac_match.group(1).replace('_', '.')
        elif re.search(r'Linux', user_agent, re.IGNORECASE):
            device_info["platform"] = "Linux"
            device_info["device_type"] = "desktop"
        
        # Browser detection
        if re.search(r'Chrome', user_agent, re.IGNORECASE) and not re.search(r'Edg|OPR', user_agent):
            device_info["browser"] = "Chrome"
            chrome_match = re.search(r'Chrome/(\d+)', user_agent)
            if chrome_match:
                device_info["browser_version"] = chrome_match.group(1)
        elif re.search(r'Safari', user_agent, re.IGNORECASE) and not re.search(r'Chrome', user_agent):
            device_info["browser"] = "Safari"
            safari_match = re.search(r'Version/(\d+[.\d]*)', user_agent)
            if safari_match:
                device_info["browser_version"] = safari_match.group(1)
        elif re.search(r'Firefox', user_agent, re.IGNORECASE):
            device_info["browser"] = "Firefox"
            firefox_match = re.search(r'Firefox/(\d+[.\d]*)', user_agent)
            if firefox_match:
                device_info["browser_version"] = firefox_match.group(1)
        elif re.search(r'Edg', user_agent, re.IGNORECASE):
            device_info["browser"] = "Edge"
        elif re.search(r'OPR', user_agent, re.IGNORECASE):
            device_info["browser"] = "Opera"
        
        # Mobile app detection (if custom header exists)
        app_version = request.headers.get("x-app-version")
        if app_version:
            device_info["app_version"] = app_version
            device_info["device_type"] = "mobile"  # Assume mobile if app version exists
    
    # Set defaults if not detected
    if "platform" not in device_info:
        device_info["platform"] = "Unknown"
    if "device_type" not in device_info:
        device_info["device_type"] = "desktop"
    if "device_name" not in device_info:
        device_info["device_name"] = "Unknown Device"
    if "browser" not in device_info:
        device_info["browser"] = "Unknown"
    
    return device_info
