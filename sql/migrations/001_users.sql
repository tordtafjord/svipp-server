-- +goose Up

CREATE TABLE users (
                       id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                       name TEXT,
                       phone TEXT NOT NULL UNIQUE,
                       email TEXT UNIQUE,
                       password TEXT,
                       device_token TEXT,
                       temporary BOOLEAN,
                       rate_total INTEGER NOT NULL DEFAULT 0,
                       rates INTEGER NOT NULL DEFAULT 0,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE driver (
                       id INTEGER PRIMARY KEY,
                       status TEXT NOT NULL DEFAULT 'UNAVAILABLE',
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMPTZ DEFAULT NULL,
                       FOREIGN KEY (id) REFERENCES users(id)
);

-- Create a spatial index on the current_location column

-- +goose Down
DROP TABLE driver;
DROP TABLE users;
