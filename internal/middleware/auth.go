package middleware

import (
	"errors"
	"net/http"
	"os"
	"start/internal/handlers"

	"github.com/golang-jwt/jwt/v5"
)

// Auth функция-посредник механизма middleware
// для проверки аутентификации
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Проверяем наличие пароля в env
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			// Получаем куки с токеном
			cookie, err := r.Cookie("token")
			if err != nil {
				handlers.ResponseWithErrorJSON(w, http.StatusUnauthorized, errors.New("ошибка авторизации"))
				return
			}

			// Записываем значение токена
			jwtTokenString := cookie.Value

			// Валидируем и проверяем токен
			token, err := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
				secret := []byte("mysecretkey")
				return secret, nil
			})
			// log.Println(token.Valid)
			// log.Println(token)

			// Возвращаем ошибку, если токен недействителен
			if err != nil || !token.Valid {
				handlers.ResponseWithErrorJSON(w, http.StatusUnauthorized, errors.New("ошибка авторизации"))
				return
			}

		}
		next(w, r)
	})
}
