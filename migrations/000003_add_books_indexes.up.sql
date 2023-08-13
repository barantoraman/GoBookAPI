CREATE INDEX IF NOT EXISTS books_isbn_idx ON books (isbn);
CREATE INDEX IF NOT EXISTS books_title_idx ON books USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS books_author_idx ON books USING GIN (to_tsvector('simple', author));
CREATE INDEX IF NOT EXISTS books_genres_idx ON books USING GIN (genres);