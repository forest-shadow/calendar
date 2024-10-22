package application

import (
	"context"
	"fmt"

	"github.com/forest-shadow/calendar/internal/config"
	router "github.com/forest-shadow/calendar/internal/controllers/http"
	"github.com/forest-shadow/calendar/internal/database"
	"github.com/forest-shadow/calendar/internal/logger"
	"github.com/forest-shadow/calendar/internal/transport/http"
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server
	logger     logger.Logger
	db         *database.DB
}

func newApp(ctx context.Context) (*App, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	logger, err := logger.NewLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	db, err := database.NewDB(&cfg.DB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	eventsDomain := buildEventsDomain(db.Connection, logger)
	errChan := eventsDomain.eventsService.StartEventJobs(ctx)
	go func() {
		defer close(errChan)
		for err := range errChan {
			if err != nil {
				logger.Errorf("error during event jobs: %w", err)
				return
			}
		}
	}()

	router := router.NewRouter(logger, eventsDomain.eventsService)
	httpServer, err := http.NewServer(&cfg.HTTP, logger, router)
	if err != nil {
		return nil, fmt.Errorf("failed to create http server: %w", err)
	}

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		logger:     logger,
		db:         db,
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

func (app *App) shutdown() {
	if err := app.db.Close(); err != nil {
		app.logger.Error("close db connection: ", err.Error())
	}
	app.logger.Info("DB connection: closed")

	if err := app.httpServer.Stop(); err != nil {
		app.logger.Errorf("failed to stop http server: %w", err)
	}
	app.logger.Info("Appication successfully shutted down")
}

func Run(ctx context.Context) error {
	app, err := newApp(ctx)
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
