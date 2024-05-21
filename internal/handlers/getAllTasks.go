package handlers

import (
	"encoding/json"
	"net/http"
	"start/internal/models"
	"time"
)

type TaskResponse struct {
	Tasks []models.Task `json:"tasks"`
}

// GetTasks возвращает ближайшие задачи по дате
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	var tasks []models.Task
	var err error

	// Получаем значение параметра search из URL-запроса и присваиваем его переменной
	findWord := r.URL.Query().Get("search")

	// Когда параметр search пустой, то возвращаем все задачи
	// иначе возвращаем задачи указанной по дате или по слову
	if findWord == "" {
		tasks, err = h.Db.GetAllTasks()
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, errGetId)
			return
		}
	} else {
		searchDate, err := time.Parse("02.01.2006", findWord)
		if err != nil {
			tasks, err = h.Db.SearchWordDB(findWord)
			if err != nil {
				ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			tasks, err = h.Db.SearchDateDB(searchDate.Format("20060102"))
			if err != nil {
				ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// Конвертируем ответ в JSON
	response := TaskResponse{
		Tasks: tasks,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}
