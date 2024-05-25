package handlers

import (
	"net/http"
	"start/internal/models"
	"strconv"
)

// DeleteTask удаляет задачу по ID
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	task.ID = r.URL.Query().Get("id")

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

	// Удаляем задачу
	err = h.Db.DoneTasks(numTaskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	// Отправляем ответ с пустым телом
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))

}
