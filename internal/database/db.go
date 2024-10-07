package database

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/logger"
)

type DB struct {
	connection *sql.DB
	logger     logger.Logger
	Close      func() error
}

func NewDB(cfg *config.DB, logger logger.Logger) (*DB, error) {
	connCfg, err := pgx.ParseURI(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("parse URI: %w", err)
	}
	db := stdlib.OpenDB(connCfg)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("check connection: %w", err)
	}

	logger.Info("db connection success")

	return &DB{
		connection: db,
		logger:     logger,
		Close:      db.Close,
	}, nil
}
