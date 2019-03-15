package api

import (
	"encoding/json"
	"net/http"

	"github.com/karolispx/golang-rh-todo/helpers"
	"github.com/karolispx/golang-rh-todo/models"
)

// User information
type User struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Password2  string `json:"password2"`
	Token      string `json:"token"`
	LastAction string `json:"last_action"`
}

// Register - process user registration
func Register(w http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)

	// Unable to parse payload
	if err != nil {
		helpers.RestAPIRespond(w, r, "Please provide all information.", "", "error", 422)

		return
	}

	if user.Email == "" || user.Password == "" || user.Password2 == "" {
		// Check if user email and both passwords have been provided
		helpers.RestAPIRespond(w, r, "Please provide all information.", "", "error", 422)

		return
	} else if helpers.ValidateEmailAddress(user.Email) != true {
		// Check if user email provided is valid
		helpers.RestAPIRespond(w, r, "Email address is not valid!", "", "error", 422)

		return
	} else if user.Password != user.Password2 {
		// Check if password provided match
		helpers.RestAPIRespond(w, r, "Passwords do not match!", "", "error", 422)

		return
	}

	DB := helpers.InitDB()

	countUsers := models.CountUsersWithEmailAddress(DB, user.Email)

	if countUsers > 0 {
		helpers.RestAPIRespond(w, r, "This email address is taken already!", "", "error", 422)

		return
	}

	lastInsertID := models.CreateUser(DB, user.Email, user.Password)

	defer DB.Close()

	if lastInsertID < 1 {
		helpers.DefaultErrorRestAPIRespond(w, r)

		return
	}

	generateJWT := helpers.GenerateJWT(lastInsertID)

	if generateJWT != "" {
		helpers.RestAPIRespond(w, r, "You have registered successfully!", &User{Token: generateJWT}, "success", 201)

		return
	}

	helpers.DefaultErrorRestAPIRespond(w, r)

	return
}

// Login - process user login
func Login(w http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		helpers.RestAPIRespond(w, r, "Please provide all information.", "", "error", 422)

		return
	}

	if user.Email == "" || user.Password == "" {
		helpers.RestAPIRespond(w, r, "Please provide all information.", "", "error", 422)

		return
	}

	DB := helpers.InitDB()

	userID := models.UserValidLogin(DB, user.Email, user.Password)

	defer DB.Close()

	if userID > 0 {
		generateJWT := helpers.GenerateJWT(userID)

		if generateJWT != "" {
			helpers.RestAPIRespond(w, r, "You have logged in successfully!", &User{Token: generateJWT}, "success", 200)

			return
		}

		helpers.DefaultErrorRestAPIRespond(w, r)

		return
	}

	helpers.RestAPIRespond(w, r, "We could not log you in.", "", "error", 403)

	return
}
