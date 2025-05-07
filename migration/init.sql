package migration

CREATE TABLE IF NOT EXISTS(
                              id serial PRIMARY KEY,
                              name TEXT,
                              surname TEXT,
                              patronymic TEXT,
                              age INT,
                              gender TEXT,
                              nationality TEXT
);
