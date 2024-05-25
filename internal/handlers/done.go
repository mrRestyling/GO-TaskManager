package handlers

import (
	"net/http"
	"start/internal/date"
	"start/internal/models"
	"strconv"
	"time"
)

// DoneTask помечает задачу как выполненную
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	task.ID = r.URL.Query().Get("id")

	// Проверяем, что запрос содержит ID
	if task.ID == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errGetTasks)
		return
	}

	// Парсим ID
	numTaskID, err := strconv.Atoi(task.ID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// Пытаемся получить задачу по ID
	task, err = h.Db.TaskById(numTaskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	// Помечаем задачу как выполненную.
	// Если указан повтор для задачи, то обновляем дату
	// с помощью функции NextDate
	if task.Repeat == "" {
		err = h.Db.DoneTasks(numTaskID)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		nextDate, err := date.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		// Обновляем дату
		task.Date = nextDate
		err = h.Db.UpdateTask(task)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
