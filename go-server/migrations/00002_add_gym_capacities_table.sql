-- +goose Up
-- +goose StatementBegin
CREATE TABLE gym_capacities (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    location_id INTEGER NOT NULL,
    percentage INTEGER NOT NULL CHECK (percentage >= 0 AND percentage <= 100),
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'America/New_York')
);
CREATE INDEX idx_gym_capacities_location_id ON gym_capacities (location_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gym_capacities;
-- +goose StatementEnd
