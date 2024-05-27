package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"start/internal/date"
	"start/internal/models"
	"start/internal/storage"
)

// Handler объявляем структуру в которой передаем ссылку на базу данных
type Handler struct {
	Db       *storage.Database
	Password string
}

// TaskHandler добавляет задачу
func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка: ошибка десериализации JSON
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Task: %+v\n", task)

	// Проверка: валидация задачи вынесена в отдельную функцию
	task.Date, err = validateTask(task) // ПРАВКА !
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// Присвоение идентификатора к добавленной задаче
	id, err := h.Db.AddTask(task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errPostId)
		return
	}

	// Устанавливаем заголовок, чтобы показать, что это JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Устанавливаем статус ответа
	w.WriteHeader(http.StatusOK)

	// Кодируем ответ в формат JSON
	err = json.NewEncoder(w).Encode(map[string]int{"id": id})
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
}

// validateTask - проверка задачи(заголовка и даты)
func validateTask(task models.Task) (string, error) {

	// Проверка: не указан заголовок задачи
	if task.Title == "" {
		return "", errWrongTitleFormat
	}

	// Проверка: дата представлена в формате, отличном от 20060102
	if task.Date != "" {
		_, err := time.Parse(models.FormatDate, task.Date) // ПРАВКА !
		if err != nil {
			// ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongDateFormat)
			return "", errWrongDateFormat
		}
	}

	// Если поле date не указано или содержит пустую строку, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = time.Now().Format(models.FormatDate)
	}

	if len(task.Repeat) > 0 {
		// if task.Repeat[0] != 'd' && task.Repeat[0] != 'w' && task.Repeat[0] != 'm' && task.Repeat[0] != 'y' {
		// 	return "", errors.New("неверное правило повторения")
		// }
		// if task.Repeat[0] == 'd' || task.Repeat[0] == 'w' || task.Repeat[0] == 'm' {
		// 	if len(task.Repeat) < 3 {
		// 		return "", errors.New("неверное правило повторения")
		// 	}
		// }
		_, err := date.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return "", errWrongRepeatFormat
		}
	}

	// 1. Подставляется сегодняшнее число, если правило повторения не указано;
	// 2. Или подставляется следующая дата с помощью функции NextDate()

	if task.Date < time.Now().Format(models.FormatDate) {
		var err error
		if task.Repeat == "" {
			task.Date = time.Now().Format(models.FormatDate)
		} else if task.Repeat != "" {
			task.Date, err = date.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return "", errWrongRepeatFormat
			}
		}
	}
	return task.Date, nil
}
