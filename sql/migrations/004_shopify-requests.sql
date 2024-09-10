-- +goose Up

CREATE TABLE shopify_requests (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    raw_json JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE shopify_requests;