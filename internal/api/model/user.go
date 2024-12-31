package model

import "guilhermefaleiros/go-service-template/internal/application/usecase"

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (r CreateUserRequest) ToUseCaseInput() usecase.CreateUserInput {
	return usecase.CreateUserInput{
		Name:  r.Name,
		Email: r.Email,
		Phone: r.Phone,
	}
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

func NewCreateUserResponse(output *usecase.CreateUserOutput) CreateUserResponse {
	return CreateUserResponse{
		ID: output.ID,
	}
}

type RetrieveUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func NewRetrieveUserResponse(output *usecase.RetrieveUserOutput) RetrieveUserResponse {
	return RetrieveUserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
		Phone: output.Phone,
	}
}
