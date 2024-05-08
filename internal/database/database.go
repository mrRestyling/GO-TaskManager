package database

import (
	"database/sql"
	"fmt"
	"start/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

func New() (*Database, error) {
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

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS Date ON scheduler (date);`)
	if err != nil {
		return nil, err
	}

	return &Database{Db: db}, nil
}

func (d *Database) Close() error {
	return d.Db.Close()
}

func (d *Database) AddTask(task models.Task) (int, error) {

	res, err := d.Db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		// return 0, errors.New("Не удалось добавить задачу")
		// return 0, err
		// return 0, fmt.Errorf("Не удалось добавить задачу: %w", err)
		return 0, fmt.Errorf("не удалось добавить задачу: %v", err)
	}

	id, err := res.LastInsertId() // метод дб для получения последнего вставленного id
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
