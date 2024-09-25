-- +goose Up
-- Enable PostGIS extension if not already enabled
CREATE EXTENSION IF NOT EXISTS postgis;

-- Add new columns
ALTER TABLE orders
    ADD COLUMN pickup_coords GEOGRAPHY(POINT, 4326),
    ADD COLUMN delivery_coords GEOGRAPHY(POINT, 4326);


-- +goose Down
-- Remove the new columns
ALTER TABLE orders
DROP COLUMN pickup_coords,
DROP COLUMN delivery_coords;

-- Optionally, remove PostGIS extension if it's no longer needed
DROP EXTENSION IF EXISTS postgis;
