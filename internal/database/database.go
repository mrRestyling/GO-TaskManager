package database

import (
	"database/sql"
	"errors"
	"fmt"
	"start/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

const count = 50

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

func (d *Database) GetAllTasks() ([]models.Task, error) {
	tasks := []models.Task{}

	rows, err := d.Db.Query(`SELECT * FROM scheduler ORDER BY date ASC LIMIT ?`, count)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() { // ?? метод для перебора всех строк
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// TaskByID	возвращает задачу по id
// либо возвращает ошибку
func (d *Database) TaskByID(id int) (models.Task, error) {

	task := models.Task{}

	err := d.Db.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

// UpdateTaskDB обновляет задачу в базе данных
func (d *Database) UpdateTaskDB(task models.Task) error {
	_, err := d.Db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return errors.New("задача не найдена")
	}
	return nil
}

func (d *Database) DoneTasksDB(id int) error {

	_, err := d.Db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return errors.New("ошибка при удалении задачи")
	}
	return nil
}
