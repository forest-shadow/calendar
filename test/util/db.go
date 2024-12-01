package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/database"
)

type DB struct {
	DB        *database.DB
	Container testcontainers.Container
	dbURI     string
}

// NewTestDB creates a new Postgres container and returns database connection
func NewTestDB(t *testing.T) *DB {
	ctx := context.Background()

	// Create postgres container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get container connection details
	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)
	hostIP, err := container.Host(ctx)
	require.NoError(t, err)

	// Create database connection
	dbURI := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		"test",
		"test",
		hostIP,
		mappedPort.Port(),
		"testdb",
	)

	// Initialize DB connection
	db, err := database.NewDB(&config.DB{URI: dbURI}, zap.NewNop().Sugar())
	require.NoError(t, err)

	return &DB{
		DB:        db,
		Container: container,
		dbURI:     dbURI,
	}
}

// Close closes database connection and removes container
func (tdb *DB) Close(t *testing.T) {
	if tdb.DB != nil {
		err := tdb.DB.Close()
		require.NoError(t, err)
	}

	if tdb.Container != nil {
		err := tdb.Container.Terminate(context.Background())
		require.NoError(t, err)
	}
}

// RunMigrations runs all migrations from the specified directory
func (tdb *DB) RunMigrations(t *testing.T, migrationsPath string) {
	// Implementation depends on your migration tool
	// Example with goose:

	db, err := goose.OpenDBWithDriver("postgres", tdb.dbURI)
	require.NoError(t, err)
	defer db.Close()

	err = goose.Up(db, migrationsPath)
	require.NoError(t, err)
}
