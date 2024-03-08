CREATE DATABASE recipes_db;

\c recipes_db

CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    ingredients TEXT,
    instructions TEXT
);
