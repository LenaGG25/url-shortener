-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY NOT NULL ,
    short_url TEXT NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
