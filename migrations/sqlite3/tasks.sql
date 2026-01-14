CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    task_id TEXT NULL,
    status TEXT,
    description TEXT,
    work_begin INTEGER,
    work_end INTEGER,
    date INTEGER
);