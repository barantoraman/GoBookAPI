CREATE EXTENSION citext;

CREATE TABLE IF NOT EXISTS books(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    isbn text NOT NULL,
    title text NOT NULL,
    author text NOT NULL,
    genres text[] NOT NULL,
    pages integer NOT NULL,
    language text NOT NULL,
    publisher text NOT NULL,
    year integer NOT NULL,
    version integer NOT NULL DEFAULT 1
);