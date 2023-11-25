CREATE INDEX IF NOT EXISTS gifts_title_idx ON gifts USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS gifts_status_idx ON gifts USING GIN (status);
