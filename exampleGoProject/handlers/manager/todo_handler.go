package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/middleware"
	"github.com/talhacan/gateway/requests/todo"
	"github.com/talhacan/gateway/services"
)

// TodoHandler manager todo handler; Request (ManagerID core tarafından set edilir) → ToDTO → TodoService.
type TodoHandler struct {
	todoService *services.TodoService
}

// NewTodoHandler TodoHandler oluşturur.
func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// Index TodoService.Index; ManagerID core.WithValidatedQueryAndManager ile request'e set edilir.
func (todoHandler *TodoHandler) Index(ginContext *gin.Context) {
	core.WithValidatedQueryAndManager(ginContext, string(middleware.ContextKeyManagerID), func(todoIndexRequest *todo.TodoIndexRequest) *core.ServiceResponse {
		return todoHandler.todoService.Index(ginContext.Request.Context(), todoIndexRequest)
	})
}

// GetByID TodoService.GetByID.
func (todoHandler *TodoHandler) GetByID(ginContext *gin.Context) {
	core.WithValidatedUriAndManager(ginContext, string(middleware.ContextKeyManagerID), func(todoGetByIDRequest *todo.TodoGetByIDRequest) *core.ServiceResponse {
		return todoHandler.todoService.GetByID(ginContext.Request.Context(), todoGetByIDRequest)
	})
}

// Create TodoService.Create.
func (todoHandler *TodoHandler) Create(ginContext *gin.Context) {
	core.WithValidatedJSONAndManager(ginContext, string(middleware.ContextKeyManagerID), func(todoCreateRequest *todo.TodoCreateRequest) *core.ServiceResponse {
		return todoHandler.todoService.Create(ginContext.Request.Context(), todoCreateRequest)
	})
}

// Update TodoService.Update.
func (todoHandler *TodoHandler) Update(ginContext *gin.Context) {
	core.WithValidatedUriAndJSONAndManager(ginContext, string(middleware.ContextKeyManagerID), func(todoUpdateRequest *todo.TodoUpdateRequest) *core.ServiceResponse {
		return todoHandler.todoService.Update(ginContext.Request.Context(), todoUpdateRequest)
	})
}

// Delete TodoService.Delete.
func (todoHandler *TodoHandler) Delete(ginContext *gin.Context) {
	core.WithValidatedUriAndManager(ginContext, string(middleware.ContextKeyManagerID), func(todoDeleteRequest *todo.TodoDeleteRequest) *core.ServiceResponse {
		return todoHandler.todoService.Delete(ginContext.Request.Context(), todoDeleteRequest)
	})
}
