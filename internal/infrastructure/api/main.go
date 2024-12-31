package api

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"guilhermefaleiros/go-service-template/internal/application/usecase"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/controller"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api/util"
	"guilhermefaleiros/go-service-template/internal/infrastructure/database"
	"guilhermefaleiros/go-service-template/internal/infrastructure/observability"
	"guilhermefaleiros/go-service-template/internal/infrastructure/repository"
	"guilhermefaleiros/go-service-template/internal/shared"
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
		slog.Error("failed to load config")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	observability.InitMeterProvider(cfg)
	observability.InitTracer(cfg)

	conn, err := database.NewPGConnection(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userRepo := repository.NewPGUserRepository(conn)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	retrieveUserUseCase := usecase.NewRetrieveUserUseCase(userRepo)
	userController := controller.NewUserController(createUserUseCase, retrieveUserUseCase)

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(otelecho.Middleware(cfg.App.Name))

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
	slog.Info(fmt.Sprintf("Starting server on port %d", api.Cfg.API.Port))
	return api.Server.ListenAndServe()
}

func (api *API) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")
	if err := api.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	api.DB.Close()
	observability.ShutdownTracerProvider(otel.GetTracerProvider().(*sdktrace.TracerProvider))
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
			slog.Error("Failed to ping database")
			return util.InternalServerError(c, "unready")
		}
		return util.OkMessage(c, "ready")
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
