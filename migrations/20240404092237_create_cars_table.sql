-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars
(
    id      SERIAL PRIMARY KEY,
    reg_num TEXT NOT NULL,
    mark    TEXT NULL,
    model TEXT NULL ,
    ownerID INT,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (ownerID) REFERENCES people(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cars;
-- +goose StatementEnd
