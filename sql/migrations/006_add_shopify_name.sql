-- +goose Up
ALTER TABLE shopify_api_config
    ADD COLUMN location_name TEXT;
ALTER TABLE shopify_api_config
    DROP COLUMN pickup_window_start;
ALTER TABLE shopify_api_config
    DROP COLUMN pickup_window_end;
ALTER TABLE shopify_api_config
    ALTER COLUMN api_key TYPE BYTEA
        USING api_key::bytea;

CREATE TABLE business_hours (
                                api_key BYTEA NOT NULL ,
                                day_of_week INTEGER NOT NULL, -- 0-6 for Monday-Sunday
                                opens_at TIME NOT NULL ,
                                closes_at TIME NOT NULL ,
                                PRIMARY KEY (api_key, day_of_week),
                                FOREIGN KEY (api_key) REFERENCES shopify_api_config(api_key)
);





-- +goose Down
ALTER TABLE shopify_api_config
    DROP COLUMN location_name;

ALTER TABLE shopify_api_config
    ADD COLUMN pickup_window_start timestamptz;
ALTER TABLE shopify_api_config
    ADD COLUMN pickup_window_end timestamptz;

DROP TABLE business_hours;

ALTER TABLE shopify_api_config
    ALTER COLUMN api_key TYPE TEXT
        USING api_key::text;

