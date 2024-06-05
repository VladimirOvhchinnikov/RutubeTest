package database

const (
	CreateTableUsers = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER UNIQUE,
		first_name TEXT,
		last_name TEXT,
		birth_date TEXT
	);`

	CreateTableSubscriptions = `
	CREATE TABLE IF NOT EXISTS subscriptions (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		subscriber_id INTEGER,
		subscribed_to_id INTEGER,
		FOREIGN KEY(subscriber_id) REFERENCES users(telegram_id),
		FOREIGN KEY(subscribed_to_id) REFERENCES users(telegram_id)
	);`
)
