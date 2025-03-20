CREATE TABLE users (
  id TEXT PRIMARY KEY CHECK (length(id) = 32),
  name TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL
);
