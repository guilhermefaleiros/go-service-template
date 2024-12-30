package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Status    UserStatus
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

func NewUser(name, email, phone string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Status:    UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
