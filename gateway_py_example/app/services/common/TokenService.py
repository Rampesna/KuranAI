"""
Token management service - handles token storage in DB and Redis
"""
from sqlalchemy.orm import Session
from datetime import datetime, timedelta, timezone
from typing import Optional, Dict, Any
from app.models.user_token import UserToken
from app.core.RedisClient import redis_set, redis_get, redis_delete, redis_exists
from app.core.config import settings
from app.services.common.jwt import JWTService
import json


class TokenService:
    def __init__(self, db: Session):
        self.db = db
        self.jwtService = JWTService()

    def saveToken(
        self,
        userId: int,
        userType: str,
        accessToken: str,
        refreshToken: str,
        accessTokenExpiresAt: datetime,
        refreshTokenExpiresAt: datetime,
        deviceInfo: Optional[Dict[str, Any]] = None
    ) -> UserToken:
        """Save token to database and Redis"""
        # Save to database
        userToken = UserToken(
            user_id=userId,
            user_type=userType,
            access_token=accessToken,
            refresh_token=refreshToken,
            access_token_expires_at=accessTokenExpiresAt,
            refresh_token_expires_at=refreshTokenExpiresAt,
            device_info=deviceInfo,
            is_active=True,
            last_used_at=datetime.now(timezone.utc)
        )
        self.db.add(userToken)
        self.db.commit()
        self.db.refresh(userToken)

        # Save to Redis with TTL
        self._saveToRedis(userToken)

        return userToken

    def _saveToRedis(self, userToken: UserToken):
        """Save token info to Redis"""
        # Redis key format: token:{access_token}
        redisKey = f"token:{userToken.access_token}"
        
        # Calculate TTL (time until access token expires)
        now = datetime.now(timezone.utc)
        # Ensure both datetimes are timezone-aware
        if userToken.access_token_expires_at.tzinfo is None:
            expiresAt = userToken.access_token_expires_at.replace(tzinfo=timezone.utc)
        else:
            expiresAt = userToken.access_token_expires_at
        ttlSeconds = int((expiresAt - now).total_seconds())
        
        # Store token info as JSON
        tokenData = {
            "id": userToken.id,
            "user_id": userToken.user_id,
            "user_type": userToken.user_type,
            "is_active": userToken.is_active,
            "expires_at": userToken.access_token_expires_at.isoformat() if userToken.access_token_expires_at else None,
            "logged_out_at": userToken.logged_out_at.isoformat() if userToken.logged_out_at else None,
        }
        
        # Save to Redis with TTL
        if ttlSeconds > 0:
            redis_set(redisKey, json.dumps(tokenData), ttl=ttlSeconds)
        
        # Also save refresh token if exists
        if userToken.refresh_token:
            refreshKey = f"refresh_token:{userToken.refresh_token}"
            if userToken.refresh_token_expires_at:
                if userToken.refresh_token_expires_at.tzinfo is None:
                    refreshExpiresAt = userToken.refresh_token_expires_at.replace(tzinfo=timezone.utc)
                else:
                    refreshExpiresAt = userToken.refresh_token_expires_at
                refreshTtl = int((refreshExpiresAt - now).total_seconds())
            else:
                refreshTtl = None
            if refreshTtl and refreshTtl > 0:
                redis_set(refreshKey, json.dumps(tokenData), ttl=refreshTtl)

    def isTokenValid(self, accessToken: str) -> bool:
        """Check if token is valid (exists, active, not expired, not logged out)"""
        # First check Redis (fast)
        redisKey = f"token:{accessToken}"
        tokenDataStr = redis_get(redisKey)
        
        if tokenDataStr:
            try:
                tokenData = json.loads(tokenDataStr)
                # Check if token is active and not logged out
                if not tokenData.get("is_active", True):
                    return False
                if tokenData.get("logged_out_at"):
                    return False
                # Check expiry
                expiresAtStr = tokenData.get("expires_at")
                if expiresAtStr:
                    expiresAt = datetime.fromisoformat(expiresAtStr.replace('Z', '+00:00'))
                    if expiresAt < datetime.now(timezone.utc):
                        return False
                return True
            except:
                pass
        
        # If not in Redis, check database (fallback)
        userToken = self.db.query(UserToken).filter(
            UserToken.access_token == accessToken,
            UserToken.deleted_at.is_(None)
        ).first()
        
        if not userToken:
            return False
        
        # Check if active
        if not userToken.is_active:
            return False
        
        # Check if logged out
        if userToken.logged_out_at:
            return False
        
        # Check expiry
        now = datetime.now(timezone.utc)
        expiresAt = userToken.access_token_expires_at
        if expiresAt.tzinfo is None:
            expiresAt = expiresAt.replace(tzinfo=timezone.utc)
        if expiresAt < now:
            return False
        
        # Update Redis (cache it)
        self._saveToRedis(userToken)
        
        return True

    def deactivateToken(self, accessToken: str) -> bool:
        """Deactivate token (logout)"""
        # Update database
        userToken = self.db.query(UserToken).filter(
            UserToken.access_token == accessToken,
            UserToken.deleted_at.is_(None)
        ).first()
        
        if not userToken:
            return False
        
        userToken.is_active = False
        userToken.logged_out_at = datetime.now(timezone.utc)
        self.db.commit()
        
        # Update Redis
        redisKey = f"token:{accessToken}"
        tokenDataStr = redis_get(redisKey)
        if tokenDataStr:
            try:
                tokenData = json.loads(tokenDataStr)
                tokenData["is_active"] = False
                tokenData["logged_out_at"] = userToken.logged_out_at.isoformat()
                # Update with remaining TTL
                expiresAt = datetime.fromisoformat(tokenData["expires_at"].replace('Z', '+00:00'))
                now = datetime.now(timezone.utc)
                ttl = int((expiresAt - now).total_seconds())
                if ttl > 0:
                    redis_set(redisKey, json.dumps(tokenData), ttl=ttl)
            except:
                pass
        
        return True

    def validateTokenAndUpdateUsage(self, accessToken: str) -> Optional[Dict[str, Any]]:
        """
        Validate token (JWT signature, expiry, Redis/DB check, logout status) and update last_used_at.
        Returns JWT payload if valid, None otherwise.
        This method checks Redis FIRST for performance, then falls back to DB.
        """
        # 1. Validate JWT signature and expiry
        payload = self.jwtService.validateToken(accessToken)
        if not payload:
            return None  # Invalid JWT or expired signature
        
        # 2. Check Redis first (fast path)
        redisKey = f"token:{accessToken}"
        tokenDataStr = redis_get(redisKey)
        
        if tokenDataStr:
            try:
                tokenData = json.loads(tokenDataStr)
                # Check if token is active and not logged out
                if not tokenData.get("is_active", True):
                    return None  # Token is inactive
                if tokenData.get("logged_out_at"):
                    return None  # Token was logged out
                # Check expiry
                expiresAtStr = tokenData.get("expires_at")
                if expiresAtStr:
                    expiresAt = datetime.fromisoformat(expiresAtStr.replace('Z', '+00:00'))
                    if expiresAt < datetime.now(timezone.utc):
                        return None  # Token expired
                
                # Token is valid in Redis, update last_used_at in DB (async or background if performance critical)
                userTokenDbId = tokenData.get("id")
                if userTokenDbId:
                    userToken = self.db.query(UserToken).filter(UserToken.id == userTokenDbId).first()
                    if userToken:
                        userToken.last_used_at = datetime.now(timezone.utc)
                        self.db.commit()
                
                return payload  # Token is valid
            except Exception:
                pass  # Fall through to DB check
        
        # 3. If not in Redis, check database (fallback)
        userToken = self.db.query(UserToken).filter(
            UserToken.access_token == accessToken,
            UserToken.deleted_at.is_(None)
        ).first()
        
        if not userToken:
            return None  # Token not found in DB
        
        # Check if active
        if not userToken.is_active:
            return None  # Token is inactive
        
        # Check if logged out
        if userToken.logged_out_at:
            return None  # Token was logged out
        
        # Check expiry
        now = datetime.now(timezone.utc)
        expiresAt = userToken.access_token_expires_at
        if expiresAt.tzinfo is None:
            expiresAt = expiresAt.replace(tzinfo=timezone.utc)
        if expiresAt < now:
            return None  # Token expired
        
        # Update last_used_at and re-cache in Redis
        userToken.last_used_at = datetime.now(timezone.utc)
        self.db.commit()
        self._saveToRedis(userToken)  # Re-cache with updated info and TTL
        
        return payload  # Token is valid

    def updateLastUsed(self, accessToken: str):
        """Update last_used_at timestamp"""
        userToken = self.db.query(UserToken).filter(
            UserToken.access_token == accessToken,
            UserToken.deleted_at.is_(None)
        ).first()
        
        if userToken:
            userToken.last_used_at = datetime.now(timezone.utc)
            self.db.commit()

    def getUserTokens(self, userId: int, userType: str, activeOnly: bool = True) -> list:
        """Get all tokens for a user (for device management)"""
        query = self.db.query(UserToken).filter(
            UserToken.user_id == userId,
            UserToken.user_type == userType,
            UserToken.deleted_at.is_(None)
        )
        
        if activeOnly:
            query = query.filter(UserToken.is_active == True)
        
        return query.order_by(UserToken.last_used_at.desc()).all()

    def deactivateUserTokens(self, userId: int, userType: str, excludeToken: Optional[str] = None):
        """Deactivate all tokens for a user (except the current one)"""
        query = self.db.query(UserToken).filter(
            UserToken.user_id == userId,
            UserToken.user_type == userType,
            UserToken.is_active == True,
            UserToken.deleted_at.is_(None)
        )
        
        if excludeToken:
            query = query.filter(UserToken.access_token != excludeToken)
        
        tokens = query.all()
        for token in tokens:
            token.is_active = False
            token.logged_out_at = datetime.now(timezone.utc)
            # Also update Redis
            redisKey = f"token:{token.access_token}"
            redis_delete(redisKey)
        
        self.db.commit()
