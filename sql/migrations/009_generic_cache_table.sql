-- +goose Up

drop table order_price_guarantees;
DROP INDEX idx_order_price_expires_at;


CREATE UNLOGGED TABLE order_price_guarantees (
    user_id INTEGER NOT NULL,
    pickup_addr TEXT NOT NULL,
    delivery_addr TEXT NOT NULL,
    distance_meters INTEGER NOT NULL,
    driving_seconds INTEGER NOT NULL,
    price_options JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    PRIMARY KEY (delivery_addr, pickup_addr, user_id)
);

CREATE INDEX idx_order_price_expires_at ON order_price_guarantees(expires_at);



-- +goose Down
DROP TABLE order_price_guarantees;


