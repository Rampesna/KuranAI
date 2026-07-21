package user

// UserIndexRequest user listesi (sayfalı) için request.
type UserIndexRequest struct {
	PageIndex      int    `form:"pageIndex"`
	PageSize       int    `form:"pageSize"`
	Keyword        string `form:"keyword"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}

// UserGetByIDRequest user detay için request (id URL'den gelir).
type UserGetByIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// UserCreateRequest user oluşturma için request.
type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// UserUpdateRequest user güncelleme için request.
type UserUpdateRequest struct {
	ID       uint  `uri:"id" binding:"required"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Status   *int   `json:"status"`
}

// UserDeleteRequest user silme için request (id URL'den).
type UserDeleteRequest struct {
	ID uint `uri:"id" binding:"required"`
}
