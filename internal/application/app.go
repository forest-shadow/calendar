package application

import (
	"context"
	"fmt"

	"github.com/forest-shadow/calendar/internal/config"
	router "github.com/forest-shadow/calendar/internal/controllers/http"
	"github.com/forest-shadow/calendar/internal/database"
	"github.com/forest-shadow/calendar/internal/jobs"
	logger "github.com/forest-shadow/calendar/internal/logger"
	http "github.com/forest-shadow/calendar/internal/transport/http"
)

type App struct {
	cfg              *config.Config
	httpServer       *http.Server
	logger           logger.Logger
	db               *database.DB
	eventNotifierJob *jobs.EventNotifierJob
	eventCleanerJob  *jobs.EventCleanerJob
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
	appLogger := logger.With("component", "app")

	db, err := database.NewDB(&cfg.DB, appLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	eventsDomain := buildEventsDomain(db.Connection, appLogger)

	eventNotifierJob := jobs.NewEventNotifier(ctx, cfg, eventsDomain.eventsService, logger)
	eventCleanerJob := jobs.NewEventCleanerJob(ctx, cfg, eventsDomain.eventsService, logger)

	router := router.NewRouter(appLogger, eventsDomain.eventsService)
	httpServer, err := http.NewServer(&cfg.HTTP, appLogger, router)
	if err != nil {
		return nil, fmt.Errorf("failed to create http server: %w", err)
	}

	return &App{
		cfg:              cfg,
		httpServer:       httpServer,
		logger:           appLogger,
		db:               db,
		eventNotifierJob: eventNotifierJob,
		eventCleanerJob:  eventCleanerJob,
	}, nil
}

func (app *App) startCronJobs() chan error {
	errCh := make(chan error, 2)

	app.eventNotifierJob.Start(errCh)
	app.eventCleanerJob.Start(errCh)

	return errCh
}

func (app *App) start() error {
	httpConfig := app.cfg.HTTP
	if err := app.httpServer.Start(&httpConfig); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	errChan := app.startCronJobs()

	go func() {
		defer close(errChan)
		for err := range errChan {
			if err != nil {
				app.logger.Errorf("error during event jobs: %w", err)
				return
			}
		}
	}()
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
