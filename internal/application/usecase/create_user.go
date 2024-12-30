package usecase

import (
	"context"
	"errors"
	"guilhermefaleiros/go-service-template/internal/application/repository"
	"guilhermefaleiros/go-service-template/internal/domain/entity"
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
	existentUser, _ := u.repository.FindByEmail(ctx, input.Email)
	if existentUser != nil {
		return nil, ErrUserAlreadyExists
	}
	newUser := entity.NewUser(input.Name, input.Email, input.Phone)
	err := u.repository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &CreateUserOutput{
		ID: newUser.ID,
	}, nil
}

func NewCreateUserUseCase(repository repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repository: repository}
}
