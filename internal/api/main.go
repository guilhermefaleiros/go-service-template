package api

import (
	"context"
	"guilhermefaleiros/go-service-template/internal/application/usecase"
	"guilhermefaleiros/go-service-template/internal/infrastructure/database"
	"guilhermefaleiros/go-service-template/internal/infrastructure/repository"
	"guilhermefaleiros/go-service-template/internal/shared"
)

func StartAPI(environment string) {
	ctx := context.Background()
	cfg, err := shared.LoadConfig(environment)
	if err != nil {
		println("Error")
		panic(err)
	}

	conn, _ := database.NewPGConnection(ctx, cfg)
	repo := repository.NewPGUserRepository(conn)
	createUser := usecase.NewCreateUserUseCase(repo)

	result, err := createUser.Execute(ctx, usecase.CreateUserInput{
		Name:  "Guilherme",
		Email: "guilhermefaleiros2000@gmail.com",
	})
	if err != nil {
		println(err.Error())
		return
	}

	println(result.ID)

}
