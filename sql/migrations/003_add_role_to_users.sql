-- +goose Up
ALTER TABLE users
ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

-- Update existing drivers
UPDATE users
SET role = 'driver'
WHERE id IN (SELECT id FROM driver);

-- +goose Down
ALTER TABLE users
DROP COLUMN role;