package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// GetTaskByID возвращает задачу по ID
func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	// Проверяем, что запрос содержит ID
	if taskID == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errGetTasks)
		return
	}

	// Парсим ID
	numTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
	}

	// Пытаемся получить задачу по ID
	task, err := h.Db.TaskByID(numTaskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, errGetId)
		return
	}

	// Конвертируем ответ в JSON
	jsonResponse, err := json.Marshal(task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}
