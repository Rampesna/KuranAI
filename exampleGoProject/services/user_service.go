package services

import (
	"context"
	"errors"

	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/models"
	userRequest "github.com/talhacan/gateway/requests/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService users tablosu için business logic; parametre olarak request alır (DTO yok).
type UserService struct {
	tokenService *TokenService
}

// NewUserService UserService oluşturur.
func NewUserService(tokenService *TokenService) *UserService {
	return &UserService{tokenService: tokenService}
}

// Me JWT'den gelen user ID ile mevcut kullanıcıyı döner (handler tek iş: WithAuthUser → Me).
func (userService *UserService) Me(ctx context.Context, userID uint) *core.ServiceResponse {
	return userService.GetByID(ctx, &userRequest.UserGetByIDRequest{ID: userID})
}

// Register kullanıcı kaydı; e-posta ve username kontrolü sonrası Create çağrılır.
func (userService *UserService) Register(ctx context.Context, req *userRequest.RegisterRequest) *core.ServiceResponse {
	var exist models.User
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
	createReq := &userRequest.UserCreateRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Username,
	}
	return userService.Create(ctx, createReq)
}

// GetByID id ile kullanıcı döner (request ile).
func (userService *UserService) GetByID(ctx context.Context, req *userRequest.UserGetByIDRequest) *core.ServiceResponse {
	var user models.User
	if findErr := database.DB.WithContext(ctx).First(&user, req.ID).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"User not found",
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
		"User fetched successfully",
		200,
		user,
	)
}

// Index sayfalı listeleme (request ile). PageSize <= 0 ise limit yok (tümü).
func (userService *UserService) Index(ctx context.Context, req *userRequest.UserIndexRequest) *core.ServiceResponse {
	pageSize := req.PageSize
	if pageSize < 0 {
		pageSize = 0
	}
	pageIndex := req.PageIndex
	if pageIndex < 0 {
		pageIndex = 0
	}
	dbQuery := database.DB.WithContext(ctx).Model(&models.User{})
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
	var users []models.User
	if findErr := dbQuery.Find(&users).Error; findErr != nil {
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Users fetched successfully",
		200,
		map[string]interface{}{
			"list":  users,
			"total": total,
		},
	)
}

// Create yeni kullanıcı oluşturur (request ile).
func (userService *UserService) Create(ctx context.Context, req *userRequest.UserCreateRequest) *core.ServiceResponse {
	var exist models.User
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
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Username: req.Username,
		Status:   1,
	}
	if createErr := database.DB.WithContext(ctx).Create(&user).Error; createErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to create user",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"User created successfully",
		201,
		user,
	)
}

// Update kullanıcı günceller (request ile).
func (userService *UserService) Update(ctx context.Context, req *userRequest.UserUpdateRequest) *core.ServiceResponse {
	var user models.User
	if findErr := database.DB.WithContext(ctx).First(&user, req.ID).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"User not found",
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
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if saveErr := database.DB.WithContext(ctx).Save(&user).Error; saveErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to update user",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"User updated successfully",
		200,
		user,
	)
}

// Delete kullanıcı siler (request ile).
func (userService *UserService) Delete(ctx context.Context, req *userRequest.UserDeleteRequest) *core.ServiceResponse {
	result := database.DB.WithContext(ctx).Delete(&models.User{}, req.ID)
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
			"User not found",
			404,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"User deleted successfully",
		200,
		nil,
	)
}

// Login credential + password ile giriş; request içinde DeviceInfo ile token üretimi serviste yapılır.
func (userService *UserService) Login(ctx context.Context, req *userRequest.LoginRequest) *core.ServiceResponse {
	var user models.User
	if findErr := database.DB.WithContext(ctx).Where("email = ? OR username = ?", req.Credential, req.Credential).First(&user).Error; findErr != nil {
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
	if user.Status != 1 {
		return core.NewServiceResponse(
			false,
			"Account is inactive",
			403,
			nil,
		)
	}
	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); compareErr != nil {
		return core.NewServiceResponse(
			false,
			"Invalid credentials",
			401,
			nil,
		)
	}
	return userService.tokenService.IssueTokens(ctx, user.ID, models.TokenableTypeUser, req.DeviceInfo)
}
