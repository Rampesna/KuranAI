package user

// UpdateUsernameRequest kullanıcı adı güncelleme.
type UpdateUsernameRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
}

// UpdatePasswordRequest şifre güncelleme.
// CurrentPassword social login kullanıcılarında boş olabilir (şifreleri yok).
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
