package controller

import (
	"github.com/labstack/echo/v4"
	"guilhermefaleiros/go-service-template/internal/application/usecase"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/model"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/util"
)

type UserController struct {
	createUserUseCase   *usecase.CreateUserUseCase
	retrieveUserUseCase *usecase.RetrieveUserUseCase
}

func (u *UserController) Setup(e *echo.Echo) {
	e.POST("/users", u.CreateUser)
	e.GET("/users/:id", u.GetByID)
}

func (u *UserController) GetByID(c echo.Context) error {
	id := c.Param("id")
	output, err := u.retrieveUserUseCase.Execute(c.Request().Context(), id)
	if err != nil {
		return util.NotFound(c, "user not found")
	}
	response := model.NewRetrieveUserResponse(output)
	return util.Ok(c, response)
}

func (u *UserController) CreateUser(c echo.Context) error {
	var request model.CreateUserRequest
	if err := c.Bind(&request); err != nil {
		return util.BadRequest(c, "invalid request")
	}

	output, err := u.createUserUseCase.Execute(c.Request().Context(), request.ToUseCaseInput())
	if err != nil {
		return util.BadRequest(c, err.Error())
	}

	response := model.NewCreateUserResponse(output)
	return util.Created(c, response)
}

func NewUserController(createUserUseCase *usecase.CreateUserUseCase, retrieveUserUseCase *usecase.RetrieveUserUseCase) *UserController {
	return &UserController{createUserUseCase: createUserUseCase, retrieveUserUseCase: retrieveUserUseCase}
}
