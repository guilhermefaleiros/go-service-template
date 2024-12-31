package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"guilhermefaleiros/go-service-template/internal/application/usecase"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/model"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/util"
	"net/http"
)

type UserController struct {
	createUserUseCase   *usecase.CreateUserUseCase
	retrieveUserUseCase *usecase.RetrieveUserUseCase
}

func (u *UserController) Setup(r chi.Router) {
	r.Post("/", u.CreateUser)
	r.Get("/{id}", u.GetByID)
}

func (u *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	output, err := u.retrieveUserUseCase.Execute(r.Context(), id)
	if err != nil {
		util.NotFound(w, "user not found")
		return
	}
	response := model.NewRetrieveUserResponse(output)
	util.Ok(w, response)
}

func (u *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.BadRequest(w, err.Error())
		return
	}
	output, err := u.createUserUseCase.Execute(r.Context(), request.ToUseCaseInput())
	if err != nil {
		util.BadRequest(w, err.Error())
		return
	}
	response := model.NewCreateUserResponse(output)
	util.Created(w, response)
}

func NewUserController(createUserUseCase *usecase.CreateUserUseCase, retrieveUserUseCase *usecase.RetrieveUserUseCase) *UserController {
	return &UserController{createUserUseCase: createUserUseCase, retrieveUserUseCase: retrieveUserUseCase}
}
