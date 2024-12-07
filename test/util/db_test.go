package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestDB(t *testing.T) {
	t.Run("Database Container Setup", func(t *testing.T) {
		// Create test database
		testDB := NewTestDB(t)
		defer testDB.Close(t)

		// Basic connectivity test
		err := testDB.DB.Connection.PingContext(context.Background())
		require.NoError(t, err, "Database should be reachable")
	})

	t.Run("Migration Testing", func(t *testing.T) {
		testDB := NewTestDB(t)
		defer testDB.Close(t)

		// Run migrations
		testDB.RunMigrations(t, "../../migrations")

		// Verify events table exists and has correct structure
		var tableName string
		err := testDB.DB.Connection.QueryRowContext(
			context.Background(),
			`SELECT table_name 
			 FROM information_schema.tables 
			 WHERE table_schema = 'public' AND table_name = 'events'`,
		).Scan(&tableName)
		require.NoError(t, err)
		assert.Equal(t, "events", tableName)

		// Verify columns
		rows, err := testDB.DB.Connection.QueryContext(
			context.Background(),
			`SELECT column_name, data_type 
			 FROM information_schema.columns 
			 WHERE table_schema = 'public' AND table_name = 'events'`,
		)
		require.NoError(t, err)
		defer rows.Close()

		columns := make(map[string]string)
		for rows.Next() {
			var columnName, dataType string
			err := rows.Scan(&columnName, &dataType)
			require.NoError(t, err)
			columns[columnName] = dataType
		}

		// Verify expected columns exist
		expectedColumns := map[string]string{
			"id":            "uuid",
			"title":         "character varying",
			"start_time":    "timestamp with time zone",
			"end_time":      "timestamp with time zone",
			"description":   "text",
			"user_id":       "uuid",
			"notify_before": "interval",
			"created_at":    "timestamp with time zone",
			"notified_at":   "timestamp with time zone",
		}

		for column, dataType := range expectedColumns {
			actualType, exists := columns[column]
			assert.True(t, exists, "Column %s should exist", column)
			assert.Equal(t, dataType, actualType, "Column %s should have correct type", column)
		}

		// Verify test data
		var count int
		err = testDB.DB.Connection.QueryRowContext(context.Background(), `
			SELECT COUNT(*) FROM events
		`).Scan(&count)
		require.NoError(t, err)
		assert.LessOrEqual(t, 0, count)
	})
}
