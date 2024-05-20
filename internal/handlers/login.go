package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func LoginSign(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginSign")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// Проверяем введённый пароль
	if request.Password != os.Getenv("TODO_PASSWORD") {
		ResponseWithErrorJSON(w, http.StatusUnauthorized, errWrongPassword)
		return
	} else {
		log.Printf("Password: %s", request.Password)
		// создаем JWT токен его засунуть в поле token и вернуть его клиенту
		secret := []byte("mysecretkey")

		jwtToken := jwt.New(jwt.SigningMethodHS256)
		signedToken, err := jwtToken.SignedString(secret)

		if err != nil {
			ResponseWithErrorJSON(w, http.StatusInternalServerError, err)
		}

		// log.Printf("Token: %s", signedToken)

		response := map[string]string{
			"token": signedToken,
		}
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
