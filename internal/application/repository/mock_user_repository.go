package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
	"guilhermefaleiros/go-service-template/internal/domain/entity"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}
