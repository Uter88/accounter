CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    object_id INTEGER
    object_type TEXT,
    type TEXT,
    path TEXT,
    date INTEGER,
    updates TEXT
);