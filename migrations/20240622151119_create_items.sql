-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       word        TEXT NOT NULL,
       translation TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
