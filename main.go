package main

import (
 	"log"
 	"net/http"
 	"os"
	"github.com/gorilla/handlers"
	"./store"
)

func main() {

	router := store.NewRouter() // create routes

	allowedOrigins := handlers.allowedOrigins([]string{"*"})
	allowedMethods := handlers.allowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	handleCors := handlers.CORS(allowedOrigins, allowedMethods)(router)

	// CORS validations
	log.Fatal(http.ListenAndServe(":8000", handleCors))

}