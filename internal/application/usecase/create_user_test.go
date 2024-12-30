package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"guilhermefaleiros/go-service-template/internal/application/repository"
	"guilhermefaleiros/go-service-template/internal/domain/entity"
	"testing"
)

func TestCreateUserUseCase_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockUserRepository)
	NewCreateUserUseCase(mockRepo)
	useCase := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Phone: "62982498044",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

	output, err := useCase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

	mockRepo.AssertCalled(t, "FindByEmail", input.Email)
	mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("*entity.User"))
	mockRepo.AssertExpectations(t)
}

func TestCreateUserUseCase_UserAlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockUserRepository)
	useCase := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Phone: "62982498044",
	}

	existentUser := &entity.User{ID: "123", Name: "John Doe", Email: input.Email}
	mockRepo.On("FindByEmail", input.Email).Return(existentUser, nil)

	output, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyExists, err)
	assert.Nil(t, output)

	mockRepo.AssertCalled(t, "FindByEmail", input.Email)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestCreateUserUseCase_CreateFails(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockUserRepository)
	useCase := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Phone: "62982498044",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

	output, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Nil(t, output)

	mockRepo.AssertCalled(t, "FindByEmail", input.Email)
	mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("*entity.User"))
	mockRepo.AssertExpectations(t)
}
