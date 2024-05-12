package handlers

import (
	"net/http"
	"start/internal/models"
	"strconv"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	task.ID = r.URL.Query().Get("id")

	if task.ID == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errGetTasks)
		return
	}

	numTaskID, err := strconv.Atoi(task.ID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	task, err = h.Db.TaskByID(numTaskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	err = h.Db.DoneTasksDB(numTaskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))

}
