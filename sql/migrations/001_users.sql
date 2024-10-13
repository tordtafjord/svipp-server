-- +goose Up

CREATE TABLE users (
                       id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                       first_name TEXT,
                       last_name TEXT,
                       phone TEXT UNIQUE,
                       email TEXT UNIQUE,
                       password TEXT,
                       device_token TEXT,
                       temporary BOOLEAN,
                       role TEXT NOT NULL DEFAULT 'user',
                       rate_total INTEGER NOT NULL DEFAULT 0,
                       rates INTEGER NOT NULL DEFAULT 0,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE driver (
                       id BIGINT PRIMARY KEY,
                       status TEXT NOT NULL DEFAULT 'unavailable',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMPTZ DEFAULT NULL,
                       FOREIGN KEY (id) REFERENCES users(id)
);

-- Create a spatial index on the current_location column

-- +goose Down
DROP TABLE driver;
DROP TABLE users;
