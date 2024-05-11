package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"start/internal/database"
	"start/internal/date"
	"start/internal/models"
)

type Handler struct {
	Db *database.Database
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {

	// Декодируем данные из тела запроса в переменную таск
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		// http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}
	// Поле title обязательное
	// Возврат json ошибки (пересмотреть)
	if task.Title == "" {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongTitleFormat)
		// response := map[string]string{
		// 	"error": "Заголовок не может быть пустым",
		// }
		// jsonResponse, err := json.Marshal(response)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// w.Header().Set("Content-Type", "application/json; charset=UTF-8") // Указываем, что возвращаем JSON
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write(jsonResponse)
		return
	}
	// Проверяем, что дата указана в формате 20060102 и что функция time.Parse() корректно её распознаёт.
	if task.Date != "" {
		_, err := time.Parse("20060102", task.Date)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongDateFormat)

			// errorResponse := map[string]string{
			// 	"error": "Неправильный формат даты",
			// }
			// jsonResponse, _ := json.Marshal(errorResponse)

			// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			// w.WriteHeader(http.StatusBadRequest)
			// w.Write(jsonResponse)
			return
		}
	}

	// поле date не указано или содержит пустую строку, берётся сегодняшнее число：
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// Если дата меньше сегодняшнего числа, есть два варианта:
	// 1. если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число;
	// 2. при указанном правиле повторения вам нужно вычислить и записать в таблицу дату выполнения, которая будет больше сегодняшнего числа.
	//    Для этого используйте функцию NextDate(), которую вы уже написали раньше.

	if task.Date < time.Now().Format("20060102") {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else if task.Repeat != "" {
			task.Date, err = date.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				ResponseWithErrorJSON(w, http.StatusBadRequest, errWrongRepeatFormat)
				// errorResponse := map[string]string{
				// 	"error": "Неправильно задано повторение",
				// }
				// jsonResponse, _ := json.Marshal(errorResponse)

				// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				// w.WriteHeader(http.StatusBadRequest)
				// w.Write(jsonResponse)
				return
			}
		}
	}

	id, err := h.Db.AddTask(task)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, errPostId)
		// http.Error(w, "Не удалось добавить задачу", http.StatusBadRequest)
		return
	}
	// w.Write([]byte(fmt.Sprintf("%d", id))) // пишем ответ в тело запроса,
	w.Header().Set("Content-Type", "application/json; charset=UTF-8") // устанавливаем заголовок, чтобы показать, что это JSON.
	// почитать про методы (.)
	w.WriteHeader(http.StatusOK)

	// тут не надо
	err = json.NewEncoder(w).Encode(map[string]int{"id": id})
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		// http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
