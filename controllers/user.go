package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jwt-course/models"
	userRepository "jwt-course/repository/user"
	"jwt-course/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {}

func (c Controller) SignUp(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		json.NewDecoder(r.Body).Decode(&user)
		
		if user.Email == "" {
			error.Message = "Email is required field."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is required field."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)
		
		userRep := userRepository.UserRepository{}
		user = userRep.SignUp(db, user)

		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)
	}
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error
		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is required field."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is required field."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password;

		userRep := userRepository.UserRepository{}
		user, err := userRep.Login(db, user)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "User does not exist"
				utils.RespondWithError(w, http.StatusNotFound, error)
				return
			} else {
				log.Fatal(err)
			}
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
		if err != nil {
			error.Message = "Wrong password."
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err :=  utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		jwt.Token = token
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, jwt)
	}
}

func (c Controller) TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObj models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToke := strings.Split(authHeader, " ")

		if len(bearerToke) == 2 {
			authToken := bearerToke[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There is an error while checking token")
				}
				return []byte(os.Getenv("SECRET_KEY")), nil
			})

			if error != nil {
				errorObj.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObj)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			errorObj.Message = "Token is invalid"
			utils.RespondWithError(w, http.StatusUnauthorized, errorObj)
			return
		}
	})
}