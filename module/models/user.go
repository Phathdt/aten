package models

import "github.com/phathdt/service-context/core"

type User struct {
	core.SQLModel
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) Mask() {
	u.Password = ""
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreate struct {
	core.SQLModel
	Email    string `json:"email"`
	Password string `json:"password"`
}
