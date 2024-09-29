-- +goose Up

-- Add new columns
ALTER TABLE orders
    ADD COLUMN delivery_window_start TIMESTAMPTZ,
    ADD COLUMN delivery_window_end TIMESTAMPTZ;


-- +goose Down
-- Remove the new columns
ALTER TABLE orders
DROP COLUMN delivery_window_start,
DROP COLUMN delivery_window_end;

