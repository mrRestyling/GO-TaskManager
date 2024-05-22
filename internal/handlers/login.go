package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// LoginSign - обработчик POST-запроса для входа в систему
func LoginSign(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginSign")

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Определяем структуру, которая будет содержать введённый пароль
	// Password сопоставляется с password из JSON
	var request struct {
		Password string `json:"password"`
	}

	// Парсим JSON в структуру request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// Сопоставляем введённый пароль с паролем из переменной окружения
	// если пароли совпадают, то создаем JWT токен
	if request.Password != os.Getenv("TODO_PASSWORD") {
		ResponseWithErrorJSON(w, http.StatusUnauthorized, errWrongPassword)
		return
	} else {
		log.Printf("Password: %s", request.Password)

		// Создаём секретный ключ для подписи
		secret := []byte("mysecretkey")

		// Создаём JWT-токен и указываем алгоритм хеширования
		jwtToken := jwt.New(jwt.SigningMethodHS256)

		// Подписываем токен с помощью секретного ключа
		signedToken, err := jwtToken.SignedString(secret)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		}

		log.Printf("Token: %s", signedToken)

		// Записываем в тело ответа JWT-токен
		response := map[string]string{
			"token": signedToken,
		}

		// Кодируем ответ в JSON и направляем клиенту
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResponse)

		return
	}

}
