package repository

import (
	"context"
	"guilhermefaleiros/go-service-template/internal/domain/entity"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}
