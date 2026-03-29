package app

import (
	"database/sql"

	"github.com/forest-shadow/calendar/internal/config"
)

type CompositionRoot struct {
	config config.Config
	db     *sql.DB
}

func NewCompositionRoot(config config.Config, db *sql.DB) *CompositionRoot {
	return &CompositionRoot{
		config: config,
		db:     db,
	}
}
