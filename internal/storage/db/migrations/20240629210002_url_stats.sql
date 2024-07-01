-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url_stats (
    id BIGSERIAL PRIMARY KEY NOT NULL ,
    short_url TEXT NOT NULL,
    request_number INTEGER DEFAULT 0 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url_stats;
-- +goose StatementEnd
