package usecase

import (
	"context"
	"errors"
	"fmt"
	"guilhermefaleiros/go-service-template/internal/application/repository"
	"guilhermefaleiros/go-service-template/internal/domain/entity"
	"log/slog"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type CreateUserInput struct {
	Name  string
	Email string
	Phone string
}

type CreateUserOutput struct {
	ID string
}

type CreateUserUseCase struct {
	repository repository.UserRepository
}

func (u *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	existentUser, err := u.repository.FindByEmail(ctx, input.Email)
	if existentUser != nil {
		slog.Warn(fmt.Sprintf("user already exists with email: %s", existentUser.Email))
		return nil, ErrUserAlreadyExists
	}
	newUser := entity.NewUser(input.Name, input.Email, input.Phone)
	err = u.repository.Create(ctx, newUser)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create user: %v", err))
		return nil, err
	}
	slog.Info(fmt.Sprintf("user created with id: %s", newUser.ID))
	return &CreateUserOutput{
		ID: newUser.ID,
	}, nil
}

func NewCreateUserUseCase(repository repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repository: repository}
}
