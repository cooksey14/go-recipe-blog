-- +goose Up

CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    ingredients TEXT,
    instructions TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS recipes;
