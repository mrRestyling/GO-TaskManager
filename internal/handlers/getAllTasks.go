package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"start/internal/models"
)

type TaskResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Db.GetAllTasks()
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, errors.New("не удалось получить задачи"))
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := TaskResponse{
		Tasks: tasks,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}
