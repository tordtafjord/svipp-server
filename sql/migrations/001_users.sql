-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

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

CREATE TABLE business (
                        id BIGINT PRIMARY KEY,
                        name TEXT NOT NULL,
                        org_id BIGINT UNIQUE NOT NULL,
                        address TEXT NOT NULL,
                        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMPTZ DEFAULT NULL,
                        FOREIGN KEY (id) REFERENCES users(id)
);

CREATE TABLE shopify_api_config (
    api_key TEXT PRIMARY KEY,
    quote_key TEXT NOT NULL,
    business_id BIGINT NOT NULL,
    pickup_address TEXT,
    pickup_coords GEOGRAPHY(POINT, 4326),
    pickup_instructions TEXT,
    pickup_window_start TIMESTAMPTZ,
    pickup_window_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (business_id) REFERENCES business(id)
);

CREATE UNIQUE INDEX idx_shop_quote_key ON shopify_api_config(quote_key);


-- +goose Down
DROP TABLE driver;
DROP TABLE users;
DROP EXTENSION IF EXISTS postgis;

