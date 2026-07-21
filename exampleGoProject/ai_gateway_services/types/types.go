package types

// Message tek mesaj (role + content); tüm provider'lar için ortak.
type Message struct {
	Role    string `json:"role"`    // user | assistant
	Content string `json:"content"`
}
