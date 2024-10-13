-- +goose Up

-- Ratings table to store ratings between users and drivers
CREATE TABLE rating (
                        id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                        order_id bigint NOT NULL,
                        rater_id bigint NOT NULL,
                        rated_id bigint NOT NULL,
                        rating smallint NOT NULL CHECK (rating >= 1 AND rating <= 5),
                        comment TEXT,
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (order_id) REFERENCES orders(id),
                        FOREIGN KEY (rater_id) REFERENCES users(id),
                        FOREIGN KEY (rated_id) REFERENCES users(id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_user_ratings()
    RETURNS TRIGGER AS $$
BEGIN
    -- Update the user's rate_total and rates
UPDATE users
SET rate_total = rate_total + NEW.rating,
    rates = rates + 1
WHERE id = NEW.rated_id;

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