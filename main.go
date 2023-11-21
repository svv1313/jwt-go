package main

import (
	"database/sql"
	"jwt-course/controllers"
	"jwt-course/driver"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	var port string = ":8000"
	db = driver.ConnectDB()
	router := mux.NewRouter()
	controller := controllers.Controller{}

	router.HandleFunc("/signup", controller.SignUp(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleWare(controller.ProtectedEndpoint())).Methods("GET")

	log.Print("Listening on port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}