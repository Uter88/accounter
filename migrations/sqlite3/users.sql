CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	login TEXT,
	password TEXT,
	name TEXT,
	surname TEXT,
	patronymic TEXT,
	price_per_hour REAL
);