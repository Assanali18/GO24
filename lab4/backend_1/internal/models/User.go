package models

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required" validate:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
}
