package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	if taskID == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errGetTasks)
		return
	}

	numTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
	}

	task, err := h.Db.TaskByID(numTaskID)

	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, errGetId)
		return
	}
	jsonResponse, err := json.Marshal(task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}
