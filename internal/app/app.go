package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"resume-generator/api"
	repo "resume-generator/internal/infrastructure/repository/postgresql"
	"resume-generator/internal/pkg/config"
	"resume-generator/internal/pkg/logger"
	"resume-generator/internal/pkg/postgres"
	"resume-generator/internal/usecase"
	"time"
)

type App struct {
	Config     *config.Config
	Logger     *zap.Logger
	DB         *postgres.PostgresDB
	server     *http.Server
	AnimalType usecase.User
}

func NewApp(cfg *config.Config) (*App, error) {
	// l init
	l, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// postgres init
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		Logger: l,
		DB:     db,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration("2s")
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	// repositories initialization

	user := repo.NewUserRepo(a.DB)

	// use case initialization

	userUseCase := usecase.NewUserUseCase(user)

	// api init
	handler := api.NewRoute(api.RouteOption{
		ContextTimeout: contextTimeout,
		Logger:         a.Logger,
		Config:         a.Config,
		User:           userUseCase,
	})

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
