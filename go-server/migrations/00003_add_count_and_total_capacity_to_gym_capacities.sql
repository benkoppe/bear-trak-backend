-- +goose Up
-- +goose StatementBegin
ALTER TABLE gym_capacities
    ADD COLUMN total_capacity INTEGER,
    ADD COLUMN count INTEGER;

WITH totals(location_id, total_capacity) AS (
    VALUES
        -- current internal ids
        (0, 80),
        (1, 65),
        (2, 50),
        (3, 65),
        (4, 75),
        -- historical external location ids
        (5636, 80),
        (5876, 65),
        (5872, 50),
        (5868, 65),
        (7428, 75),
        -- known court/location ids from external feed
        (7867, 20),
        (7868, 16),
        (8820, 20)
)
UPDATE gym_capacities gc
SET
    total_capacity = t.total_capacity,
    count = ROUND((gc.percentage::numeric * t.total_capacity::numeric) / 100.0)::integer
FROM totals t
WHERE gc.location_id = t.location_id;

-- Fallback for legacy/unexpected location ids: treat percentage as count out of 100.
UPDATE gym_capacities
SET
    total_capacity = 100,
    count = percentage
WHERE total_capacity IS NULL OR count IS NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM gym_capacities
        WHERE total_capacity IS NULL OR count IS NULL
    ) THEN
        RAISE EXCEPTION 'Backfill failed: null values remain after fallback';
    END IF;
END $$;

ALTER TABLE gym_capacities
    ALTER COLUMN total_capacity SET NOT NULL,
    ALTER COLUMN count SET NOT NULL,
    ADD CONSTRAINT gym_capacities_total_capacity_positive CHECK (total_capacity > 0),
    ADD CONSTRAINT gym_capacities_count_nonnegative CHECK (count >= 0),
    ADD CONSTRAINT gym_capacities_count_lte_total CHECK (count <= total_capacity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE gym_capacities
    DROP CONSTRAINT IF EXISTS gym_capacities_count_lte_total,
    DROP CONSTRAINT IF EXISTS gym_capacities_count_nonnegative,
    DROP CONSTRAINT IF EXISTS gym_capacities_total_capacity_positive,
    DROP COLUMN IF EXISTS count,
    DROP COLUMN IF EXISTS total_capacity;
-- +goose StatementEnd
