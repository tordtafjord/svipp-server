-- +goose Up
CREATE TABLE orders
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL,
    sender_id           INTEGER   NOT NULL,
    recipient_id        INTEGER   NOT NULL,
    driver_id           INTEGER   DEFAULT NULL,
    pickup_address   TEXT   NOT NULL,
    delivery_address TEXT   NOT NULL,
    status              TEXT      NOT NULL DEFAULT 'PENDING',
    distance INTEGER NOT NULL,
    driving_minutes DOUBLE PRECISION NOT NULL,
    price               DOUBLE PRECISION NOT NULL,
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

CREATE TABLE temp_order
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL,
    pickup_address   TEXT   NOT NULL,
    delivery_address TEXT   NOT NULL,
    distance INTEGER,
    driving_minutes DOUBLE PRECISION,
    price DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);



-- Ratings table to store ratings between users and drivers
CREATE TABLE rating (
                         id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                         order_id INTEGER NOT NULL,
                         rater_id INTEGER NOT NULL,
                         ratee_id INTEGER NOT NULL,
                         rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (order_id) REFERENCES orders(id),
                         FOREIGN KEY (rater_id) REFERENCES users(id),
                         FOREIGN KEY (ratee_id) REFERENCES users(id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_user_ratings()
    RETURNS TRIGGER AS $$
BEGIN
    -- Update the user's rate_total and rates
    UPDATE users
    SET rate_total = rate_total + NEW.rating,
        rates = rates + 1
    WHERE id = NEW.ratee_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Create the trigger
CREATE TRIGGER after_rating_insert
    AFTER INSERT ON rating
    FOR EACH ROW
EXECUTE FUNCTION update_user_ratings();



-- +goose Down
-- Drop the trigger
DROP TRIGGER IF EXISTS after_rating_insert ON rating;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_user_ratings();

DROP TABLE rating;
DROP TABLE orders;
DROP TABLE temp_order;
