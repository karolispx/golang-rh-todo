package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/karolispx/golang-rh-todo/api"
	"github.com/karolispx/golang-rh-todo/helpers"
)

// Main function of the application
func main() {
	// Initialize routes
	initRoutes()
}

// Initializes routes
func initRoutes() {
	// Get config variables
	Config := helpers.GetConfig()

	// Set up router routes
	router := mux.NewRouter().StrictSlash(true)

	// Test route to ensure router is set up correctly
	router.HandleFunc(Config.RestAPIPath+"/test", api.TestRoute).Methods("GET")

	// Print out the URL of the API
	fmt.Println("Server is running on: " + Config.RestAPIURL + ":" + Config.Port)

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Start server
	log.Fatal(http.ListenAndServe(":"+Config.Port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
