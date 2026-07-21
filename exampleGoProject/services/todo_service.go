package services

import (
	"context"
	"errors"

	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/models"
	"github.com/talhacan/gateway/requests/todo"
	"gorm.io/gorm"
)

// TodoService todos tablosu için business logic; parametre olarak request alır (DTO yok).
type TodoService struct{}

// NewTodoService TodoService oluşturur.
func NewTodoService() *TodoService {
	return &TodoService{}
}

// GetAll tüm todo'ları döner (request ile UserID/ManagerID filtre).
func (todoService *TodoService) GetAll(ctx context.Context, req *todo.TodoIndexRequest) *core.ServiceResponse {
	dbQuery := database.DB.WithContext(ctx).Model(&models.Todo{})
	if req.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *req.UserID)
	}
	if req.ManagerID != nil {
		dbQuery = dbQuery.Where("manager_id = ?", *req.ManagerID)
	}
	var todos []models.Todo
	if findErr := dbQuery.Find(&todos).Error; findErr != nil {
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Todos fetched successfully",
		200,
		todos,
	)
}

// GetByID id ile todo döner; request içinde UserID/ManagerID ile sahiplik kontrolü.
func (todoService *TodoService) GetByID(ctx context.Context, req *todo.TodoGetByIDRequest) *core.ServiceResponse {
	var todoModel models.Todo
	if findErr := database.DB.WithContext(ctx).First(&todoModel, req.ID).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"Todo not found",
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
	if req.UserID != nil && (todoModel.UserID == nil || *todoModel.UserID != *req.UserID) {
		return core.NewServiceResponse(
			false,
			"Todo not found",
			404,
			nil,
		)
	}
	if req.ManagerID != nil && (todoModel.ManagerID == nil || *todoModel.ManagerID != *req.ManagerID) {
		return core.NewServiceResponse(
			false,
			"Todo not found",
			404,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Todo fetched successfully",
		200,
		todoModel,
	)
}

// Index sayfalı listeleme (request ile).
func (todoService *TodoService) Index(ctx context.Context, req *todo.TodoIndexRequest) *core.ServiceResponse {
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := req.PageIndex
	if pageIndex < 0 {
		pageIndex = 0
	}
	dbQuery := database.DB.WithContext(ctx).Model(&models.Todo{})
	if req.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *req.UserID)
	}
	if req.ManagerID != nil {
		dbQuery = dbQuery.Where("manager_id = ?", *req.ManagerID)
	}
	if req.Keyword != "" {
		keywordPattern := "%" + req.Keyword + "%"
		dbQuery = dbQuery.Where("title ILIKE ? OR body ILIKE ?", keywordPattern, keywordPattern)
	}
	if req.OrderBy != "" && req.OrderDirection != "" {
		dbQuery = dbQuery.Order(req.OrderBy + " " + req.OrderDirection)
	}
	var total int64
	dbQuery.Count(&total)
	offset := pageIndex * pageSize
	var todos []models.Todo
	if findErr := dbQuery.Offset(offset).Limit(pageSize).Find(&todos).Error; findErr != nil {
		return core.NewServiceResponse(
			false,
			"Database error",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Todos fetched successfully",
		200,
		map[string]interface{}{
			"list":  todos,
			"total": total,
		},
	)
}

// Create yeni todo oluşturur (request ile UserID/ManagerID core tarafından set edilir).
func (todoService *TodoService) Create(ctx context.Context, req *todo.TodoCreateRequest) *core.ServiceResponse {
	todoModel := models.Todo{
		Title:     req.Title,
		Body:      req.Body,
		UserID:    req.UserID,
		ManagerID: req.ManagerID,
	}
	if createErr := database.DB.WithContext(ctx).Create(&todoModel).Error; createErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to create todo",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Todo created successfully",
		201,
		todoModel,
	)
}

// Update todo günceller (request ile sahiplik kontrolü).
func (todoService *TodoService) Update(ctx context.Context, req *todo.TodoUpdateRequest) *core.ServiceResponse {
	var todoModel models.Todo
	if findErr := database.DB.WithContext(ctx).First(&todoModel, req.ID).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(
				false,
				"Todo not found",
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
	if req.UserID != nil && (todoModel.UserID == nil || *todoModel.UserID != *req.UserID) {
		return core.NewServiceResponse(
			false,
			"Todo not found",
			404,
			nil,
		)
	}
	if req.ManagerID != nil && (todoModel.ManagerID == nil || *todoModel.ManagerID != *req.ManagerID) {
		return core.NewServiceResponse(
			false,
			"Todo not found",
			404,
			nil,
		)
	}
	if req.Title != "" {
		todoModel.Title = req.Title
	}
	if req.Body != "" {
		todoModel.Body = req.Body
	}
	if req.Done != nil {
		todoModel.Done = *req.Done
	}
	if saveErr := database.DB.WithContext(ctx).Save(&todoModel).Error; saveErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to update todo",
			500,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Todo updated successfully",
		200,
		todoModel,
	)
}

// Delete todo siler (request ile sahiplik kontrolü).
func (todoService *TodoService) Delete(ctx context.Context, req *todo.TodoDeleteRequest) *core.ServiceResponse {
	dbQuery := database.DB.WithContext(ctx).Where("id = ?", req.ID)
	if req.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *req.UserID)
	}
	if req.ManagerID != nil {
		dbQuery = dbQuery.Where("manager_id = ?", *req.ManagerID)
	}
	result := dbQuery.Delete(&models.Todo{})
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
			"Todo not found",
			404,
			nil,
		)
	}
	return core.NewServiceResponse(
		true,
		"Todo deleted successfully",
		200,
		nil,
	)
}
