package models

// User rappresenta un utente della piattaforma
type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"` // "teacher" o "student"
}
