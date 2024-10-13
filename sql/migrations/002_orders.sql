-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE orders
(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id bigint NOT NULL,
    sender_id           bigint   NOT NULL,
    recipient_id        bigint   NOT NULL,
    driver_id           bigint   DEFAULT NULL,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    pickup_address   TEXT   NOT NULL,
    delivery_address TEXT   NOT NULL,
    pickup_coords GEOGRAPHY(POINT, 4326),
    delivery_coords GEOGRAPHY(POINT, 4326),
    status              TEXT      NOT NULL DEFAULT 'pending',
    distance_meters INTEGER NOT NULL,
    driving_seconds INTEGER NOT NULL,
    price_cents               INTEGER NOT NULL,
    delivery_window_start TIMESTAMPTZ,
    delivery_window_end TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMPTZ DEFAULT NULL, -- confirmed by recipient
    accepted_at TIMESTAMPTZ DEFAULT NULL,  -- accepted by driver
    picked_up_at TIMESTAMPTZ DEFAULT NULL,
    delivered_at TIMESTAMPTZ DEFAULT NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cancelled_at        TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    FOREIGN KEY (driver_id) REFERENCES users(id)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE UNIQUE INDEX idx_orders_public_id ON orders(public_id);





-- +goose Down
DROP TABLE orders;
DROP EXTENSION IF EXISTS postgis;
