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

func TestRetrieveUserUseCase_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockUserRepository)
	useCase := NewRetrieveUserUseCase(mockRepo)

	userID := "123"

	expectedUser := &entity.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Phone: "62982498044",
	}

	mockRepo.On("FindByID", mock.Anything, userID).Return(expectedUser, nil)

	output, err := useCase.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedUser.ID, output.ID)
	assert.Equal(t, expectedUser.Name, output.Name)
	assert.Equal(t, expectedUser.Email, output.Email)
	assert.Equal(t, expectedUser.Phone, output.Phone)

	mockRepo.AssertCalled(t, "FindByID", mock.Anything, userID)
	mockRepo.AssertExpectations(t)
}

func TestRetrieveUserUseCase_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockUserRepository)
	useCase := NewRetrieveUserUseCase(mockRepo)

	userID := "123"

	mockRepo.On("FindByID", ctx, userID).Return((*entity.User)(nil), errors.New("user not found"))

	output, err := useCase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "user not found")

	mockRepo.AssertCalled(t, "FindByID", ctx, userID)
	mockRepo.AssertExpectations(t)
}

func TestRetrieveUserUseCase_RepositoryError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockUserRepository)
	useCase := NewRetrieveUserUseCase(mockRepo)

	userID := "123"
	
	mockRepo.On("FindByID", ctx, userID).Return((*entity.User)(nil), errors.New("database error"))

	output, err := useCase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "database error")

	mockRepo.AssertCalled(t, "FindByID", ctx, userID)
	mockRepo.AssertExpectations(t)
}
