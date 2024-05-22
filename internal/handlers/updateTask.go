package handlers

import (
	"encoding/json"
	"net/http"
	"start/internal/date"
	"start/internal/models"
	"strconv"
	"time"
)

// UpdateTask редактирует задачу
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	if task.Title == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongTitleFormat)
		return
	}

	if task.Date != "" {
		_, err := time.Parse("20060102", task.Date)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongDateFormat)
			return
		}
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

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

	if task.ID == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errGetTasks)
		return
	}

	numTaskID, err := strconv.Atoi(task.ID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongTaskIDFormat)
		return
	}

	_, err = h.Db.TaskByIdDB(numTaskID)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	// Обновляем задачу в базе данных
	err = h.Db.UpdateTaskDB(task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errUpdateDb)
		return
	}

	// Возвращаем JSON-ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]int{})
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
}
