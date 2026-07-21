package manager

// ManagerIndexRequest manager listesi için request.
type ManagerIndexRequest struct {
	PageIndex      int    `form:"pageIndex"`
	PageSize       int    `form:"pageSize"`
	Keyword        string `form:"keyword"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}

// ManagerGetByIDRequest manager detay için request.
type ManagerGetByIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// ManagerCreateRequest manager oluşturma için request.
type ManagerCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// ManagerUpdateRequest manager güncelleme için request.
type ManagerUpdateRequest struct {
	ID       uint  `uri:"id" binding:"required"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Status   *int   `json:"status"`
}

// ManagerDeleteRequest manager silme için request.
type ManagerDeleteRequest struct {
	ID uint `uri:"id" binding:"required"`
}
