package user

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CompanyID string    `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email, companyID string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		CompanyID: companyID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
