package application

import (
	"fmt"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/transport/http"
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server

	// logger           *logger.Logger
    // db               *database.Database
}

func newApp() (*App, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	httpServer, err := http.NewServer(&cfg.HTTP)
	if err != nil {
		return nil, fmt.Errorf("failed to create http server: %w", err)
	}
	return &App{
		cfg: cfg,
		httpServer: httpServer,
	}, nil
}

func (app *App) start() error {
	if err := app.httpServer.Start(&app.cfg.HTTP); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func Run() error {
	app, err := newApp()
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	return app.start()
}
