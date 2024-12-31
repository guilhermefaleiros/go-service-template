package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"guilhermefaleiros/go-service-template/internal/api/controller"
	"guilhermefaleiros/go-service-template/internal/api/util"
	"guilhermefaleiros/go-service-template/internal/application/usecase"
	"guilhermefaleiros/go-service-template/internal/infrastructure/database"
	"guilhermefaleiros/go-service-template/internal/infrastructure/repository"
	"guilhermefaleiros/go-service-template/internal/shared"
	"log"
	"net/http"
)

func StartAPI(environment string) {
	ctx := context.Background()
	cfg, err := shared.LoadConfig(environment)
	if err != nil {
		panic(err)
	}

	conn, err := database.NewPGConnection(ctx, cfg)

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	userRepo := repository.NewPGUserRepository(conn)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	retrieveUserUseCase := usecase.NewRetrieveUserUseCase(userRepo)

	userController := controller.NewUserController(createUserUseCase, retrieveUserUseCase)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	SetupHealthChecks(r, conn)

	r.Route("/users", func(r chi.Router) {
		userController.Setup(r)
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.API.Port), r)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func SetupHealthChecks(r chi.Router, conn *pgxpool.Pool) {
	r.Get("/live", func(w http.ResponseWriter, r *http.Request) {
		util.OkMessage(w, "ready")
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		err := conn.Ping(r.Context())
		if err != nil {
			util.InternalServerError(w, "unready")
			return
		}
		util.OkMessage(w, "ready")
	})
}
