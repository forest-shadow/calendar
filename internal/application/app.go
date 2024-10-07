package application

import (
	"context"
	"fmt"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/logger"
	"github.com/forest-shadow/calendar/internal/transport/http"
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server
	logger     logger.Logger
}

func newApp() (*App, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	logger, err := logger.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	httpServer, err := http.NewServer(&cfg.HTTP, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create http server: %w", err)
	}

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (app *App) start() error {
	httpConfig := app.cfg.HTTP
	if err := app.httpServer.Start(&httpConfig); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	app.logger.Infof("Appication started at port: %v", httpConfig.Port)
	return nil
}

func (app *App) shutdown() error {
	if err := app.httpServer.Stop(); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}
	app.logger.Info("Appication successfully shutted down")

	return nil
}

func Run(ctx context.Context) error {
	app, err := newApp()
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	if err := app.start(); err != nil {
		return fmt.Errorf("error during start: %w", err)
	}

	defer app.shutdown()

	<-ctx.Done()

	return nil
}
