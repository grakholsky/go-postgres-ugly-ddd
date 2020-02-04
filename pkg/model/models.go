package model

import (
	"time"
)

type (
	Role struct {
		ID   string `json:"id" gorm:"primary_key"`
		Name string `json:"name"`
	}

	User struct {
		ID        string    `json:"id" gorm:"primary_key"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Password  string    `json:"-"`
		Salt      []byte    `json:"-"`
		Avatar    *string   `json:"avatar"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserRole struct {
		ID     string `json:"id" gorm:"primary_key"`
		UserID string `json:"user_id"`
		RoleID string `json:"role_id"`
	}
)
