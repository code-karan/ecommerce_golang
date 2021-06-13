package main

import (
 	"log"
 	"net/http"
	"github.com/gorilla/handlers"
	"go-ecommerce/store"
)

func main() {

	router := store.NewRouter() // create routes

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	handleCors := handlers.CORS(allowedOrigins, allowedMethods)(router)

	// CORS validations
	log.Fatal(http.ListenAndServe("localhost:8000", handleCors))

}