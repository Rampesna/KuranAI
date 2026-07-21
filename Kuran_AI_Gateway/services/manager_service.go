package services

import (
	"context"
	"errors"

	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ManagerService manager tablosu için business logic.
type ManagerService struct {
	tokenService *TokenService
}

// NewManagerService ManagerService oluşturur.
func NewManagerService(tokenService *TokenService) *ManagerService {
	return &ManagerService{tokenService: tokenService}
}

// ManagerLoginRequest manager login girdisi.
type ManagerLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login manager giriş yapar.
func (s *ManagerService) Login(ctx context.Context, req *ManagerLoginRequest) *core.ServiceResponse {
	var manager models.Manager
	if err := database.DB.WithContext(ctx).Where("email = ?", req.Email).First(&manager).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "Invalid credentials", 401, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}
	if manager.Status != 1 {
		return core.NewServiceResponse(false, "Account is inactive", 403, nil)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(manager.Password), []byte(req.Password)); err != nil {
		return core.NewServiceResponse(false, "Invalid credentials", 401, nil)
	}
	return s.tokenService.IssueTokens(ctx, manager.ID, models.TokenableTypeManager, nil)
}

// Me JWT'den gelen manager ID ile mevcut yöneticiyi döner.
func (s *ManagerService) Me(ctx context.Context, managerID uint) *core.ServiceResponse {
	var manager models.Manager
	if err := database.DB.WithContext(ctx).First(&manager, managerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "Manager not found", 404, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}
	return core.NewServiceResponse(true, "Manager fetched successfully", 200, manager)
}
