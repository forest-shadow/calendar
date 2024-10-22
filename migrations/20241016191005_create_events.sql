-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    description TEXT,
    user_id UUID NOT NULL,
    notify_before INTERVAL DAY TO SECOND, -- ISO 8601 interval format
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notified_at TIMESTAMPTZ
);

-- Insert test data
INSERT INTO events (id, title, start_time, end_time, description, user_id, notify_before)
VALUES
    (gen_random_uuid(), 'Team Meeting', NOW() + INTERVAL '00:00:20', NOW() + INTERVAL '1 day 1 hour', 'Weekly team sync', '123e4567-e89b-12d3-a456-426614174000', INTERVAL '00:00:05'),
    (gen_random_uuid(), 'Project Deadline', NOW() + INTERVAL '00:00:25', NOW() + INTERVAL '7 days', 'Complete project documentation', '123e4567-e89b-12d3-a456-426614174000', INTERVAL '00:00:06'),
    (gen_random_uuid(), 'Lunch with Client', NOW() + INTERVAL '00:00:30', NOW() + INTERVAL '3 days 2 hours', 'Discuss new contract', '987e6543-e21b-12d3-a456-426614174000', INTERVAL '00:00:07');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
