package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser_Creation(t *testing.T) {
	user := NewUser("John Doe", "johndoegmail.com", "62982498044")

	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "johndoegmail.com", user.Email)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, UserStatusActive, user.Status)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUser_Deactivate(t *testing.T) {
	user := NewUser("John Doe", "johndoegmail.com", "62982498044")
	user.Deactivate()

	assert.Equal(t, UserStatusInactive, user.Status)
	assert.NotZero(t, user.UpdatedAt)
}
