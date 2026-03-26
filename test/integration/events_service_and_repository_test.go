package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/forest-shadow/calendar/internal/events"
	testutil "github.com/forest-shadow/calendar/test/util"
)

type RepositoryTestSuite struct {
	suite.Suite
	testDB *testutil.DB
	repo   *events.Repository
	ctx    context.Context
}

func (s *RepositoryTestSuite) SetupSuite() {
	s.testDB = testutil.NewTestDB(s.T())
	s.testDB.RunMigrations(s.T(), "../../migrations")
	s.repo = events.NewEventsRepository(s.testDB.DB)
	s.ctx = context.Background()
}

func (s *RepositoryTestSuite) TearDownSuite() {
	s.testDB.Close(s.T())
}

func (s *RepositoryTestSuite) TestEventsCRUD() {
	t := s.T()

	testCases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "Create and retrieve event",
			fn: func(t *testing.T) {
				event := events.Event{
					ID:        uuid.New(),
					Title:     "Test Event",
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC().Add(time.Hour),
					UserID:    uuid.New(),
				}

				// Create
				err := s.repo.Create(s.ctx, event)
				require.NoError(t, err)

				// Retrieve
				retrieved, err := s.repo.Get(s.ctx, event.ID)
				require.NoError(t, err)
				require.Equal(t, event.Title, retrieved.Title)
			},
		},
		{
			name: "Update event",
			fn: func(t *testing.T) {
				// Create initial event
				event := events.Event{
					ID:        uuid.New(),
					Title:     "Initial Title",
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC().Add(time.Hour),
					UserID:    uuid.New(),
				}
				err := s.repo.Create(s.ctx, event)
				require.NoError(t, err)

				// Update
				updateData := events.UpdateEvent{
					"title": "Updated Title",
				}
				err = s.repo.Update(s.ctx, event.ID, updateData)
				require.NoError(t, err)

				// Verify
				updated, err := s.repo.Get(s.ctx, event.ID)
				require.NoError(t, err)
				require.Equal(t, "Updated Title", updated.Title)
			},
		},
		{
			name: "Delete event",
			fn: func(t *testing.T) {
				event := events.Event{
					ID:        uuid.New(),
					Title:     "To Be Deleted",
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC().Add(time.Hour),
					UserID:    uuid.New(),
				}
				err := s.repo.Create(s.ctx, event)
				require.NoError(t, err)

				// Delete
				err = s.repo.Delete(s.ctx, event.ID.String())
				require.NoError(t, err)

				// Verify deletion
				_, err = s.repo.Get(s.ctx, event.ID)
				require.Error(t, err) // Should return error as event is deleted
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.fn)
	}
}

func (s *RepositoryTestSuite) TestGetRecentEvents() {
	t := s.T()

	// Create events with different notification times
	now := time.Now().UTC()
	events := []events.Event{
		{
			ID:           uuid.New(),
			Title:        "Soon to notify",
			StartTime:    now.Add(10 * time.Minute),
			EndTime:      now.Add(1 * time.Hour),
			UserID:       uuid.New(),
			NotifyBefore: ptr("PT15M"), // 15 minutes before
		},
		{
			ID:           uuid.New(),
			Title:        "Later notification",
			StartTime:    now.Add(1 * time.Hour),
			EndTime:      now.Add(2 * time.Hour),
			UserID:       uuid.New(),
			NotifyBefore: ptr("PT30M"), // 30 minutes before
		},
	}

	// Create events
	for _, event := range events {
		err := s.repo.Create(s.ctx, event)
		require.NoError(t, err)
	}

	// Test GetRecentEvents
	recentEvents, err := s.repo.GetRecentEvents(s.ctx, now)
	require.NoError(t, err)
	require.NotNil(t, recentEvents)
	require.Len(t, *recentEvents, 1) // Should only get the first event
}

func ptr(s string) *string {
	return &s
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
