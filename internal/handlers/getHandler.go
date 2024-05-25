package handlers

import (
	"net/http"
	"time"

	"start/internal/date"
	"start/internal/models"
)

// NextDateHandler - обработчик GET-запроса для получения следующей даты
// *для проверки функции NextDate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	// Извлекаем параметры GET-запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Парсим дату в формате "20060102"
	now, err := time.Parse(models.FormatDate, nowStr)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongDateFormat)
		return
	}

	// Вызываем функцию NextDate для получения следующей даты
	// на основе текущей даты, заданной даты и повтора
	nextDate, err := date.NextDate(now, dateStr, repeatStr)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// Записываем содержимое переменной nextDate в ответный поток "w" в виде байтового массива
	w.Write([]byte(nextDate))
}
