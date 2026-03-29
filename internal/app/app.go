package app

import (
	"context"
	"database/sql"
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
	db               *sql.DB
	eventNotifierJob *jobs.EventNotifierJob
	eventCleanerJob  *jobs.EventCleanerJob
}

func newApp(ctx context.Context) (*App, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}

	logger, err := logger.NewLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}
	appLogger := logger.With("component", "app")

	db, err := database.NewDB(&cfg.DB, appLogger)
	if err != nil {
		return nil, fmt.Errorf("create db: %w", err)
	}

	eventsDomain := buildEventsDomain(db, appLogger)

	eventNotifierJob, err := jobs.NewEventNotifier(ctx, cfg, eventsDomain.eventsService)
	if err != nil {
		return nil, fmt.Errorf("create event notifier: %w", err)
	}
	eventCleanerJob := jobs.NewEventCleanerJob(ctx, cfg, eventsDomain.eventsService, logger)

	handlers := router.NewHandlers(appLogger, eventsDomain.eventsService)
	router := router.NewRouter(handlers)
	httpServer, err := http.NewServer(&cfg.HTTP, appLogger, router)
	if err != nil {
		return nil, fmt.Errorf("create http server: %w", err)
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

func (app *App) start() error {
	httpConfig := app.cfg.HTTP
	if err := app.httpServer.Start(&httpConfig); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}

	if err := jobs.RunJobs([]jobs.Job{app.eventNotifierJob, app.eventCleanerJob}); err != nil {
		return fmt.Errorf("cron jobs: %w", err)
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
		app.logger.Errorf("stop http server: %w", err)
	}
	app.logger.Info("Appication successfully shutted down")
}

func Run(ctx context.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return fmt.Errorf("create app: %w", err)
	}

	if err := app.start(); err != nil {
		return fmt.Errorf("start app: %w", err)
	}

	defer app.shutdown()

	<-ctx.Done()

	return ctx.Err()
}
