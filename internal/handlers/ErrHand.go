package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errWrongDateFormat   = errors.New("неправильный формат даты")
	errWrongRepeatFormat = errors.New("неправильный формат повтора")
	errWrongTitleFormat  = errors.New("неправильный формат заголовка")
	errPostId            = errors.New("не удалось добавить задачу по id")
)

func ResponseWithErrorJSON(w http.ResponseWriter, status int, err error) {
	errorResponse := map[string]string{
		"error": err.Error(),
	}
	jsonResponse, _ := json.Marshal(errorResponse)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}
