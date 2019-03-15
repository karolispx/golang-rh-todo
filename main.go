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

	// Set up router and routes
	router := mux.NewRouter().StrictSlash(true)

	// User authentication
	router.HandleFunc(Config.RestAPIPath+"/auth/register", api.Register).Methods("POST") // User registration
	router.HandleFunc(Config.RestAPIPath+"/auth/login", api.Login).Methods("POST")       // User login

	// User tasks. Require authentication.
	router.HandleFunc(Config.RestAPIPath+"/tasks", api.GetTasks).Methods("GET")       // Get tasks
	router.HandleFunc(Config.RestAPIPath+"/tasks", api.CreateTask).Methods("POST")    // Create a task
	router.HandleFunc(Config.RestAPIPath+"/tasks", api.DeleteTasks).Methods("DELETE") // Delete all tasks

	router.HandleFunc(Config.RestAPIPath+"/tasks/{taskid}", api.GetTask).Methods("GET")       // Get a specific task
	router.HandleFunc(Config.RestAPIPath+"/tasks/{taskid}", api.UpdateTask).Methods("PUT")    // Update a specific task
	router.HandleFunc(Config.RestAPIPath+"/tasks/{taskid}", api.DeleteTask).Methods("DELETE") // Delete a specific task

	router.HandleFunc(Config.RestAPIPath+"/tasks/watch/{taskid}", api.WatchTask).Methods("PUT")     // Watch a specific task
	router.HandleFunc(Config.RestAPIPath+"/tasks/unwatch/{taskid}", api.UnwatchTask).Methods("PUT") // Unwatch a specific task

	// Print out the URL of the API
	fmt.Println("Server is running on: " + Config.RestAPIURL + ":" + Config.Port)

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Start server
	log.Fatal(http.ListenAndServe(":"+Config.Port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
