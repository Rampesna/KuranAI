package core

// ServiceResponse tüm servis metodlarının döndüreceği standart yapı.
type ServiceResponse struct {
	IsSuccess  bool        `json:"isSuccess"`
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data,omitempty"`
}

// NewServiceResponse yeni bir ServiceResponse oluşturur.
func NewServiceResponse(isSuccess bool, message string, statusCode int, data interface{}) *ServiceResponse {
	return &ServiceResponse{
		IsSuccess:  isSuccess,
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}
