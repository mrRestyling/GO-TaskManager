package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"start/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

const count = 50

type Database struct {
	Db *sql.DB
}

// New создает новый экземпляр базы данных
func New() (*Database, error) {

	// *** Задание со звездочкой, переменная окружения
	fileName := os.Getenv("TODO_DBFILE")
	if fileName == "" {
		fileName = "scheduler.db"
	}
	// ***

	// Открываем соединение к базе данных
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу (если она не создана)
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

	// Создаем индекс по полю date
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS Date ON scheduler (date);`)
	if err != nil {
		return nil, err
	}

	return &Database{Db: db}, nil
}

// Функция для закрытия соединения
func (d *Database) Close() error {
	return d.Db.Close()
}

// AddTask добавляет задачу в базу данных
func (d *Database) AddTask(task models.Task) (int, error) {

	res, err := d.Db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить задачу: %v", err)
	}

	id, err := res.LastInsertId() // метод дб для получения последнего вставленного id
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetAllTasks возвращает все задачи из базы данных
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

// DeleteTaskDB удаляет задачу из базы данных
func (d *Database) DoneTasksDB(id int) error {

	_, err := d.Db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return errors.New("ошибка при удалении задачи")
	}
	return nil
}

// SearchWordDB возвращает задачи по ключевому слову
func (d *Database) SearchWordDB(title string) ([]models.Task, error) {
	tasks := []models.Task{}

	rows, err := d.Db.Query(`SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`, "%"+title+"%", "%"+title+"%", count)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	// Перебор всех строк в результате выполнения запроса к базе данных
	// для каждой строки создается новый объект models.Task,
	// и затем данные из строки сканируются
	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// SearchDateDB возвращает задачи по дате
func (d *Database) SearchDateDB(date string) ([]models.Task, error) {
	tasks := []models.Task{}

	rows, err := d.Db.Query(`SELECT * FROM scheduler WHERE date = ? LIMIT ?`, date, count)
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
