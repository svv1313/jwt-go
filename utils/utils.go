package utils

import (
	"encoding/json"
	"jwt-course/models"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	ResponseJSON(w, error)
}

func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss": "issued by google",
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
