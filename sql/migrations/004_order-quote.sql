-- +goose Up

CREATE UNLOGGED TABLE order_quote (
    user_id bigint NOT NULL,
    pickup_address TEXT NOT NULL,
    delivery_address TEXT NOT NULL,
    distance_meters INTEGER NOT NULL,
    driving_seconds INTEGER NOT NULL,
    price_options JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, pickup_address, delivery_address)
);

-- +goose Down
DROP TABLE order_quote;


