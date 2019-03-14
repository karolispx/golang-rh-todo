package api

import (
	"encoding/json"
	"net/http"

	"github.com/karolispx/golang-rh-todo/helpers"
	"github.com/karolispx/golang-rh-todo/models"
)

// User information
type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
	Token     string `json:"token"`
}

// Register - process user registration
func Register(w http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		response := "Please provide all information."

		helpers.RestAPIRespond(w, r, response, "error", 422)

		return
	}

	if user.Email == "" || user.Password == "" || user.Password2 == "" {
		response := "Please provide all information."

		helpers.RestAPIRespond(w, r, response, "error", 422)

		return
	} else if helpers.ValidateEmailAddress(user.Email) != true {
		response := "Email address is not valid!"

		helpers.RestAPIRespond(w, r, response, "error", 422)

		return
	} else if user.Password != user.Password2 {
		response := "Passwords do not match!"

		helpers.RestAPIRespond(w, r, response, "error", 422)

		return
	}

	DB := helpers.InitDB()

	countUsers := models.CountUsersWithEmailAddress(DB, user.Email)

	if countUsers > 0 {
		response := "This email address is taken already!"

		helpers.RestAPIRespond(w, r, response, "error", 422)
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
		helpers.RestAPIRespond(w, r, generateJWT, "success", 200)

		return
	}

	helpers.DefaultErrorRestAPIRespond(w, r)

	return
}
