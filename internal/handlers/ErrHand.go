package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Вспомогательные ошибки
var (
	errWrongDateFormat   = errors.New("неправильный формат даты")
	errWrongRepeatFormat = errors.New("неправильный формат повтора")
	errWrongTitleFormat  = errors.New("неправильный формат заголовка")
	errPostId            = errors.New("не удалось добавить задачу по id")
	errGetId             = errors.New("не удалось получить задачу по id")
	errGetTasks          = errors.New("не указан идентификатор задач")
	errWrongTaskIDFormat = errors.New("неверный формат идентификатора задачи")
	errWrongPassword     = errors.New("ошибка авторизации")
)

// ResponseWithErrorJSON возвращает JSON с ошибкой
func ResponseWithErrorJSON(w http.ResponseWriter, status int, err error) {
	errorResponse := map[string]string{
		"error": err.Error(),
	}
	jsonResponse, _ := json.Marshal(errorResponse)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}
