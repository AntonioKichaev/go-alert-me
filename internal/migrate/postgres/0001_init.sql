CREATE TABLE IF NOT EXISTs gauges
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    value DOUBLE PRECISION NOT NULL
    );

CREATE TABLE IF NOT EXISTS counts(
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL UNIQUE,
    value INTEGER NOT NULL
    );