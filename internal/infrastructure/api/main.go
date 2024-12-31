package api

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"guilhermefaleiros/go-service-template/internal/application/usecase"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/controller"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/util"
	"guilhermefaleiros/go-service-template/internal/infrastructure/database"
	"guilhermefaleiros/go-service-template/internal/infrastructure/repository"
	"guilhermefaleiros/go-service-template/internal/shared"
	"log"
	"log/slog"
	"net/http"
)

type API struct {
	Router *echo.Echo
	Server *http.Server
	DB     *pgxpool.Pool
	Cfg    *shared.Config
}

func NewAPI(environment string) (*API, error) {
	ctx := context.Background()

	cfg, err := shared.LoadConfig(environment)
	if err != nil {
		log.Println("failed to load config")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	conn, err := database.NewPGConnection(ctx, cfg)
	if err != nil {
		log.Println("failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userRepo := repository.NewPGUserRepository(conn)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	retrieveUserUseCase := usecase.NewRetrieveUserUseCase(userRepo)
	userController := controller.NewUserController(createUserUseCase, retrieveUserUseCase)

	e := echo.New()

	SetupMetadata(e, conn)
	userController.Setup(e)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.API.Port),
		Handler: e,
	}

	return &API{
		Router: e,
		Server: server,
		DB:     conn,
		Cfg:    cfg,
	}, nil
}

func (api *API) Start() error {
	log.Printf("Starting server on port %d\n", api.Cfg.API.Port)
	return api.Server.ListenAndServe()
}

func (api *API) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if err := api.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	api.DB.Close()
	return nil
}

func SetupMetadata(e *echo.Echo, conn *pgxpool.Pool) {
	e.GET("/live", func(c echo.Context) error {
		return util.OkMessage(c, "live")
	})

	e.GET("/ready", func(c echo.Context) error {
		ctx := c.Request().Context()
		err := conn.Ping(ctx)
		if err != nil {
			slog.Info("Failed to ping database")
			return util.InternalServerError(c, "unready")
		}
		return util.OkMessage(c, "ready")
	})
}
