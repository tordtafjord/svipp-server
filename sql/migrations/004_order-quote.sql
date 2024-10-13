-- +goose Up

CREATE UNLOGGED TABLE order_quote (
    user_id bigint NOT NULL,
    pickup_address TEXT NOT NULL,
    delivery_address TEXT NOT NULL,
    delivery_window_start TIMESTAMPTZ,
    delivery_window_end TIMESTAMPTZ,
    distance_meters INTEGER NOT NULL,
    driving_seconds INTEGER NOT NULL,
    price_options JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (delivery_address, pickup_address, user_id)
);

-- +goose Down
DROP TABLE order_quote;


