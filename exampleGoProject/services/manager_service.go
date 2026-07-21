package services

import (
	"context"
	"errors"

	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/models"
	managerRequest "github.com/talhacan/gateway/requests/manager"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ManagerService managers tablosu için business logic; parametre olarak request alır (DTO yok).
type ManagerService struct {
	tokenService *TokenService
}

// NewManagerService ManagerService oluşturur.
func NewManagerService(tokenService *TokenService) *ManagerService {
	return &ManagerService{tokenService: tokenService}
}

// Me JWT'den gelen manager ID ile mevcut manager'ı döner (handler tek iş: WithAuthManager → Me).
func (managerService *ManagerService) Me(ctx context.Context, managerID uint) *core.ServiceResponse {
	return managerService.GetByID(ctx, &managerRequest.ManagerGetByIDRequest{ID: managerID})
}

// GetByID id ile manager döner (request ile).
func (managerService *ManagerService) GetByID(ctx context.Context, req *managerRequest.ManagerGetByIDRequest) *core.ServiceResponse {
	var manager models.Manager
	if findErr := database.DB.WithContext(ctx).First(&manager, req.ID).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"Manager not found",
				404,
				nil,
			)
		}
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Manager fetched successfully",
		200,
		manager,
	)
}

// Index sayfalı listeleme (request ile). PageSize <= 0 ise limit yok (tümü).
func (managerService *ManagerService) Index(ctx context.Context, req *managerRequest.ManagerIndexRequest) *core.ServiceResponse {
	pageSize := req.PageSize
	if pageSize < 0 {
		pageSize = 0
	}
	pageIndex := req.PageIndex
	if pageIndex < 0 {
		pageIndex = 0
	}
	dbQuery := database.DB.WithContext(ctx).Model(&models.Manager{})
	if req.Keyword != "" {
		keywordPattern := "%" + req.Keyword + "%"
		dbQuery = dbQuery.Where("name ILIKE ? OR email ILIKE ? OR username ILIKE ?", keywordPattern, keywordPattern, keywordPattern)
	}
	if req.OrderBy != "" && req.OrderDirection != "" {
		dbQuery = dbQuery.Order(req.OrderBy + " " + req.OrderDirection)
	}
	var total int64
	dbQuery.Count(&total)
	offset := pageIndex * pageSize
	if pageSize > 0 {
		dbQuery = dbQuery.Offset(offset).Limit(pageSize)
	} else {
		dbQuery = dbQuery.Offset(offset)
	}
	var managers []models.Manager
	if findErr := dbQuery.Find(&managers).Error; findErr != nil {
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Managers fetched successfully",
		200,
		map[string]interface{}{
			"list":  managers,
			"total": total,
		},
	)
}

// Create yeni manager oluşturur (request ile).
func (managerService *ManagerService) Create(ctx context.Context, req *managerRequest.ManagerCreateRequest) *core.ServiceResponse {
	var exist models.Manager
	if findErr := database.DB.WithContext(ctx).Where("email = ?", req.Email).First(&exist).Error; findErr == nil {
		return core.NewServiceResponse(
			false,
			"Email already exists",
			400,
			nil,
		)
	}
	if findErr := database.DB.WithContext(ctx).Where("username = ?", req.Username).First(&exist).Error; findErr == nil {
		return core.NewServiceResponse(
			false,
			"Username already exists",
			400,
			nil,
		)
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to hash password",
			500,
			nil,
		)
	}
	manager := models.Manager{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Username: req.Username,
		Status:   1,
	}
	if createErr := database.DB.WithContext(ctx).Create(&manager).Error; createErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to create manager",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Manager created successfully",
		201,
		manager,
	)
}

// Update manager günceller (request ile).
func (managerService *ManagerService) Update(ctx context.Context, req *managerRequest.ManagerUpdateRequest) *core.ServiceResponse {
	var manager models.Manager
	if findErr := database.DB.WithContext(ctx).First(&manager, req.ID).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"Manager not found",
				404,
				nil,
			)
		}
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	if req.Email != "" {
		manager.Email = req.Email
	}
	if req.Name != "" {
		manager.Name = req.Name
	}
	if req.Username != "" {
		manager.Username = req.Username
	}
	if req.Status != nil {
		manager.Status = *req.Status
	}
	if saveErr := database.DB.WithContext(ctx).Save(&manager).Error; saveErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to update manager",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Manager updated successfully",
		200,
		manager,
	)
}

// Delete manager siler (request ile).
func (managerService *ManagerService) Delete(ctx context.Context, req *managerRequest.ManagerDeleteRequest) *core.ServiceResponse {
	result := database.DB.WithContext(ctx).Delete(&models.Manager{}, req.ID)
	if result.Error != nil {
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	if result.RowsAffected == 0 {
		return core.NewServiceResponse(
			false,
			"Manager not found",
			404,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Manager deleted successfully",
		200,
		nil,
	)
}

// Login credential + password ile giriş; request içinde DeviceInfo ile token üretimi serviste yapılır.
func (managerService *ManagerService) Login(ctx context.Context, req *managerRequest.LoginRequest) *core.ServiceResponse {
	var manager models.Manager
	if findErr := database.DB.WithContext(ctx).Where("email = ? OR username = ?", req.Credential, req.Credential).First(&manager).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"Invalid credentials",
				401,
				nil,
			)
		}
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	if manager.Status != 1 {
		return core.NewServiceResponse(
			false,
			"Account is inactive",
			403,
			nil,
		)
	}
	if compareErr := bcrypt.CompareHashAndPassword([]byte(manager.Password), []byte(req.Password)); compareErr != nil {
		return core.NewServiceResponse(
			false,
			"Invalid credentials",
			401,
			nil,
		)
	}
	return managerService.tokenService.IssueTokens(ctx, manager.ID, models.TokenableTypeManager, req.DeviceInfo)
}
