package usecase

import (
	"context"
	"guilhermefaleiros/go-service-template/internal/application/repository"
)

type RetrieveUserOutput struct {
	ID    string
	Name  string
	Email string
	Phone string
}

type RetrieveUserUseCase struct {
	repository repository.UserRepository
}

func (r *RetrieveUserUseCase) Execute(ctx context.Context, id string) (*RetrieveUserOutput, error) {
	user, err := r.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &RetrieveUserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}, nil
}

func NewRetrieveUserUseCase(repository repository.UserRepository) *RetrieveUserUseCase {
	return &RetrieveUserUseCase{repository: repository}
}
