-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS owners (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT NULL,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE owners;
-- +goose StatementEnd
