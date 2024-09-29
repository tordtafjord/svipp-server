-- +goose Up
CREATE UNLOGGED TABLE phone_auth (
    phone TEXT PRIMARY KEY,
    auth_code TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP
);

CREATE INDEX idx_phone_auth_expires_at ON phone_auth(expires_at);

CREATE UNLOGGED TABLE order_price_guarantees (
    hash_key BYTEA PRIMARY KEY,
    price_cents INTEGER NOT NULL,
    distance INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP
);

CREATE INDEX idx_order_price_expires_at ON phone_auth(expires_at);


-- Change to using int and store seconds and cents instead
-- Change driving_minutes to INTEGER and convert to seconds
ALTER TABLE orders
ALTER COLUMN driving_minutes TYPE INTEGER USING (driving_minutes * 60)::INTEGER;

-- Rename driving_minutes to driving_seconds
ALTER TABLE orders
    RENAME COLUMN driving_minutes TO driving_seconds;

-- Change price to INTEGER and convert to cents
ALTER TABLE orders
ALTER COLUMN price TYPE INTEGER USING (price * 100)::INTEGER;

-- Rename price to price_cents
ALTER TABLE orders
    RENAME COLUMN price TO price_cents;



DROP TABLE temp_order;


-- +goose Down
DROP TABLE phone_auth;
DROP TABLE order_price_guarantees;
DROP TABLE temp_order;

