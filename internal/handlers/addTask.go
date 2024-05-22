package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"start/internal/database"
	"start/internal/date"
	"start/internal/models"
)

// Handler объявляем структуру в которой передаем ссылку на базу данных
type Handler struct {
	Db *database.Database
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

	// Проверка: не указан заголовок задачи
	if task.Title == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongTitleFormat)
		return
	}

	// Проверка: дата представлена в формате, отличном от 20060102
	if task.Date != "" {
		_, err := time.Parse("20060102", task.Date)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongDateFormat)
			return
		}
	}

	// Если поле date не указано или содержит пустую строку, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// 1. Подставляется сегодняшнее число, если правило повторения не указано;
	// 2. Или подставляется следующая дата с помощью функции NextDate()

	if task.Date < time.Now().Format("20060102") {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else if task.Repeat != "" {
			task.Date, err = date.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongRepeatFormat)
				return
			}
		}
	}

	// Присвоение идентификатора к добавленной задаче
	id, err := h.Db.AddTaskDB(task)
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
