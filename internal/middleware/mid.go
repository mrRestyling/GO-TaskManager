package middleware

import (
	"errors"
	"net/http"
	"os"
	"start/internal/handlers"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if password is set
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			// Get JWT token from cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				handlers.ResponseWithErrorJSON(w, http.StatusUnauthorized, errors.New("ошибка авторизации"))
				return
			}

			jwtTokenString := cookie.Value

			// Validate JWT token
			token, err := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
				// Verify secret key
				secret := []byte("mysecretkey")
				return secret, nil
			})
			// log.Println(token.Valid)
			// log.Println(token)
			if err != nil || !token.Valid {
				handlers.ResponseWithErrorJSON(w, http.StatusUnauthorized, errors.New("ошибка авторизации"))
				return
			}

		}
		next(w, r)
	})
}
