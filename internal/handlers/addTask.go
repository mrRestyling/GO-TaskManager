package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"start/internal/date"
	"start/internal/storage"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		// Декодируем данные из тела запроса в переменную таск
		var task storage.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
			return
		}
		// Поле title обязательное
		if task.Title == "" {
			http.Error(w, "Заголовок не может быть пустым", http.StatusBadRequest)
			return
		}
		// Проверяем, что дата указана в формате 20060102 и что функция time.Parse() корректно её распознаёт.
		if task.Date != "" {
			_, err := time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
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
					http.Error(w, "Неправильно задано повторение", http.StatusBadRequest)
					return
				}

			}
		}

	}

}
