package database

import (
	"database/sql"
	"log"
)

func InitDatabase(dbFile string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CreateTableUsers)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized and table created if not exists")
	return db, nil
}
