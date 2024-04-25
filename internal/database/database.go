package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		date     CHAR(8)      NOT NULL DEFAULT "",
		title    VARCHAR(256) NOT NULL DEFAULT "",
		comment  TEXT,
		repeat   VARCHAR(16)  NOT NULL DEFAULT ""
	);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE INDEX Date ON scheduler (date);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
