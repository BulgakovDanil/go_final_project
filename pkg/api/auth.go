package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key")

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedPass := password

	// Если пароль не установлен - авторизация не нужна
	if expectedPass == "" {
		writeJson(w, map[string]string{"token": ""}, http.StatusOK)
		return
	}

	// Читаем пароль из запроса
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Проверяем пароль
	if req.Password != expectedPass {
		writeJsonError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": req.Password,                         // сохраняем пароль в токен
		"exp":           time.Now().Add(8 * time.Hour).Unix(), // 8 часов
	})

	// Подписываем токен
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writeJsonError(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	writeJson(w, map[string]string{
		"token": tokenString,
	}, http.StatusOK)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			// Получаем токен из куки
			cookie, err := r.Cookie("token")
			if err != nil {
				writeJsonError(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			jwtToken := cookie.Value
			var valid bool

			// Парсим и проверяем JWT токен
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				// Проверяем метод подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtKey, nil
			})

			if err == nil && token.Valid {
				// Получаем claims
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					// Проверяем что пароль в токене совпадает с текущим
					if passHash, ok := claims["password_hash"].(string); ok {
						valid = (passHash == pass)
					}
				}
			}

			if !valid {
				writeJsonError(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
