package todo

// TodoIndexRequest todo listesi için request (UserID/ManagerID core tarafından set edilir).
type TodoIndexRequest struct {
	PageIndex      int    `form:"pageIndex"`
	PageSize       int    `form:"pageSize"`
	Keyword        string `form:"keyword"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
	UserID         *uint  // core.WithValidatedQueryAndUser/AndManager tarafından set edilir
	ManagerID      *uint
}

// SetUserID core.UserIDSetter implementasyonu.
func (r *TodoIndexRequest) SetUserID(id uint) { r.UserID = &id }

// SetManagerID core.ManagerIDSetter implementasyonu.
func (r *TodoIndexRequest) SetManagerID(id uint) { r.ManagerID = &id }

// TodoGetByIDRequest todo detay için request.
type TodoGetByIDRequest struct {
	ID        uint `uri:"id" binding:"required"`
	UserID    *uint
	ManagerID *uint
}

func (r *TodoGetByIDRequest) SetUserID(id uint)    { r.UserID = &id }
func (r *TodoGetByIDRequest) SetManagerID(id uint) { r.ManagerID = &id }

// TodoCreateRequest todo oluşturma için request.
type TodoCreateRequest struct {
	Title     string `json:"title" binding:"required"`
	Body      string `json:"body"`
	UserID    *uint
	ManagerID *uint
}

func (r *TodoCreateRequest) SetUserID(id uint)    { r.UserID = &id }
func (r *TodoCreateRequest) SetManagerID(id uint) { r.ManagerID = &id }

// TodoUpdateRequest todo güncelleme için request.
type TodoUpdateRequest struct {
	ID        uint   `uri:"id" binding:"required"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Done      *bool  `json:"done"`
	UserID    *uint
	ManagerID *uint
}

func (r *TodoUpdateRequest) SetUserID(id uint)    { r.UserID = &id }
func (r *TodoUpdateRequest) SetManagerID(id uint) { r.ManagerID = &id }

// TodoDeleteRequest todo silme için request.
type TodoDeleteRequest struct {
	ID        uint `uri:"id" binding:"required"`
	UserID    *uint
	ManagerID *uint
}

func (r *TodoDeleteRequest) SetUserID(id uint)    { r.UserID = &id }
func (r *TodoDeleteRequest) SetManagerID(id uint) { r.ManagerID = &id }
