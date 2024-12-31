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
	"log/slog"
	"net/http"
)

type API struct {
	Router *chi.Mux
	Server *http.Server
	DB     *pgxpool.Pool
	Cfg    *shared.Config
}

func NewAPI(environment string) (*API, error) {
	ctx := context.Background()

	cfg, err := shared.LoadConfig(environment)
	if err != nil {
		slog.Info("failed to load config")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	conn, err := database.NewPGConnection(ctx, cfg)
	if err != nil {
		slog.Info("failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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

	SetupMetadata(r, conn)

	r.Route("/users", func(r chi.Router) {
		userController.Setup(r)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.API.Port),
		Handler: r,
	}

	return &API{
		Router: r,
		Server: server,
		DB:     conn,
		Cfg:    cfg,
	}, nil
}

func (api *API) Start() error {
	slog.Info(fmt.Sprintf("Starting server on port %d", api.Cfg.API.Port))
	return api.Server.ListenAndServe()
}

func (api *API) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")
	if err := api.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	api.DB.Close()
	return nil
}

func SetupMetadata(r chi.Router, conn *pgxpool.Pool) {
	r.Get("/live", func(w http.ResponseWriter, r *http.Request) {
		util.OkMessage(w, "ready")
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		err := conn.Ping(r.Context())
		if err != nil {
			slog.Error("Failed to ping database")
			util.InternalServerError(w, "unready")
			return
		}
		util.OkMessage(w, "ready")
	})
}
