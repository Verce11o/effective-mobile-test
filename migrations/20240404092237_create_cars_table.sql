-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars
(
    id      SERIAL PRIMARY KEY,
    reg_num TEXT NOT NULL,
    mark    TEXT NULL,
    model TEXT NULL ,
    year INT,
    ownerID INT,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (ownerID) REFERENCES owners(id)
);

CREATE INDEX idx_id_created_at ON cars(id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cars;
-- +goose StatementEnd
