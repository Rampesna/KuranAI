package core

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/dto"
)

// UserIDSetter request'in auth sonrası UserID alması için.
type UserIDSetter interface {
	SetUserID(uint)
}

// ManagerIDSetter request'in auth sonrası ManagerID alması için.
type ManagerIDSetter interface {
	SetManagerID(uint)
}

// DeviceInfoSetter request'in Login/Refresh için DeviceInfo alması için.
type DeviceInfoSetter interface {
	SetDeviceInfo(*dto.DeviceInfo)
}

// WithValidatedJSON JSON body'yi bind eder, validate eder; hata varsa 422/400 yazar.
func WithValidatedJSON[T any](ginContext *gin.Context, responseCallback func(*T) *ServiceResponse) {
	var validatedRequest T
	if bindErr := ginContext.ShouldBindJSON(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedQuery query parametrelerini bind eder, validate eder.
func WithValidatedQuery[T any](ginContext *gin.Context, responseCallback func(*T) *ServiceResponse) {
	var validatedRequest T
	if bindErr := ginContext.ShouldBindQuery(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid query", 400, nil))
		return
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedUri URI parametrelerini bind eder, validate eder.
func WithValidatedUri[T any](ginContext *gin.Context, responseCallback func(*T) *ServiceResponse) {
	var validatedRequest T
	if bindErr := ginContext.ShouldBindUri(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedUriAndJSON URI sonra JSON bind eder (Update gibi).
func WithValidatedUriAndJSON[T any](ginContext *gin.Context, responseCallback func(*T) *ServiceResponse) {
	var validatedRequest T
	if uriBindErr := ginContext.ShouldBindUri(&validatedRequest); uriBindErr != nil {
		if serviceResponse := ValidationErrorResponse(uriBindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if jsonBindErr := ginContext.ShouldBindJSON(&validatedRequest); jsonBindErr != nil {
		if serviceResponse := ValidationErrorResponse(jsonBindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithAuthID JWT'den context'e yazılmış ID'yi key ile alır; geçerliyse responseCallback(id) çağrılır.
func WithAuthID(ginContext *gin.Context, contextKey string, responseCallback func(id uint) *ServiceResponse) {
	idValue, exists := ginContext.Get(contextKey)
	if !exists {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	id, ok := idValue.(uint)
	if !ok {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	serviceResponse := responseCallback(id)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithAuthToken Authorization: Bearer <token> header'ını okur; geçerliyse responseCallback(accessToken) çağrılır.
func WithAuthToken(ginContext *gin.Context, responseCallback func(accessToken string) *ServiceResponse) {
	authHeader := ginContext.GetHeader("Authorization")
	if authHeader == "" {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Authorization header required", 401, nil))
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Invalid authorization format", 401, nil))
		return
	}
	accessToken := strings.TrimSpace(parts[1])
	if accessToken == "" {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Token required", 401, nil))
		return
	}
	serviceResponse := responseCallback(accessToken)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

func getAuthIDFromContext(ginContext *gin.Context, contextKey string) *uint {
	idValue, exists := ginContext.Get(contextKey)
	if !exists {
		return nil
	}
	id, ok := idValue.(uint)
	if !ok {
		return nil
	}
	return &id
}

// WithValidatedQueryAndUser Query bind + validate; UserID context'ten alınır, request'e set edilir.
func WithValidatedQueryAndUser[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	userID := getAuthIDFromContext(ginContext, contextKey)
	if userID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if bindErr := ginContext.ShouldBindQuery(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid query", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(UserIDSetter); ok {
		setter.SetUserID(*userID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedQueryAndManager Query bind + validate; ManagerID context'ten alınır, request'e set edilir.
func WithValidatedQueryAndManager[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	managerID := getAuthIDFromContext(ginContext, contextKey)
	if managerID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if bindErr := ginContext.ShouldBindQuery(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid query", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(ManagerIDSetter); ok {
		setter.SetManagerID(*managerID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedJSONAndUser JSON bind + validate; UserID context'ten alınır, request'e set edilir.
func WithValidatedJSONAndUser[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	userID := getAuthIDFromContext(ginContext, contextKey)
	if userID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if bindErr := ginContext.ShouldBindJSON(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(UserIDSetter); ok {
		setter.SetUserID(*userID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedJSONAndManager JSON bind + validate; ManagerID context'ten alınır, request'e set edilir.
func WithValidatedJSONAndManager[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	managerID := getAuthIDFromContext(ginContext, contextKey)
	if managerID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if bindErr := ginContext.ShouldBindJSON(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(ManagerIDSetter); ok {
		setter.SetManagerID(*managerID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedUriAndUser Uri bind + validate; UserID context'ten alınır, request'e set edilir.
func WithValidatedUriAndUser[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	userID := getAuthIDFromContext(ginContext, contextKey)
	if userID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if bindErr := ginContext.ShouldBindUri(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(UserIDSetter); ok {
		setter.SetUserID(*userID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedUriAndManager Uri bind + validate; ManagerID context'ten alınır, request'e set edilir.
func WithValidatedUriAndManager[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	managerID := getAuthIDFromContext(ginContext, contextKey)
	if managerID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if bindErr := ginContext.ShouldBindUri(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(ManagerIDSetter); ok {
		setter.SetManagerID(*managerID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedUriAndJSONAndUser Uri + JSON bind; UserID context'ten alınır, request'e set edilir.
func WithValidatedUriAndJSONAndUser[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	userID := getAuthIDFromContext(ginContext, contextKey)
	if userID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if uriBindErr := ginContext.ShouldBindUri(&validatedRequest); uriBindErr != nil {
		if serviceResponse := ValidationErrorResponse(uriBindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if jsonBindErr := ginContext.ShouldBindJSON(&validatedRequest); jsonBindErr != nil {
		if serviceResponse := ValidationErrorResponse(jsonBindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(UserIDSetter); ok {
		setter.SetUserID(*userID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// WithValidatedUriAndJSONAndManager Uri + JSON bind; ManagerID context'ten alınır, request'e set edilir.
func WithValidatedUriAndJSONAndManager[T any](ginContext *gin.Context, contextKey string, responseCallback func(*T) *ServiceResponse) {
	managerID := getAuthIDFromContext(ginContext, contextKey)
	if managerID == nil {
		ginContext.JSON(http.StatusUnauthorized, NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var validatedRequest T
	if uriBindErr := ginContext.ShouldBindUri(&validatedRequest); uriBindErr != nil {
		if serviceResponse := ValidationErrorResponse(uriBindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if jsonBindErr := ginContext.ShouldBindJSON(&validatedRequest); jsonBindErr != nil {
		if serviceResponse := ValidationErrorResponse(jsonBindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if setter, ok := any(&validatedRequest).(ManagerIDSetter); ok {
		setter.SetManagerID(*managerID)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}

// DeviceInfoFromContext gin context'ten DeviceInfo üretir.
func DeviceInfoFromContext(ginContext *gin.Context) *dto.DeviceInfo {
	return &dto.DeviceInfo{
		UserAgent:  ginContext.GetHeader("User-Agent"),
		IPAddress:  ginContext.ClientIP(),
		DeviceInfo: ginContext.GetHeader("X-Device-Info"),
		DeviceType: ginContext.GetHeader("X-Device-Type"),
	}
}

// WithValidatedJSONAndDeviceInfo JSON bind + validate; DeviceInfo context'ten alınır, request'e set edilir.
func WithValidatedJSONAndDeviceInfo[T any](ginContext *gin.Context, responseCallback func(*T) *ServiceResponse) {
	var validatedRequest T
	if bindErr := ginContext.ShouldBindJSON(&validatedRequest); bindErr != nil {
		if serviceResponse := ValidationErrorResponse(bindErr); serviceResponse != nil {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			return
		}
		ginContext.JSON(http.StatusBadRequest, NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	deviceInfo := DeviceInfoFromContext(ginContext)
	if setter, ok := any(&validatedRequest).(DeviceInfoSetter); ok {
		setter.SetDeviceInfo(deviceInfo)
	}
	serviceResponse := responseCallback(&validatedRequest)
	ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
}
