package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS snippets (
            id SERIAL PRIMARY KEY,
            name TEXT UNIQUE NOT NULL,
            language TEXT,
            tags TEXT,
            content TEXT
        )
    `)
	if err != nil {
		return nil, err
	}
	return db, nil
}
