package models

import (
	"database/sql"

	"github.com/karolispx/golang-rh-todo/helpers"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// CountUsersWithEmailAddress to see if user with this email address exists already
func CountUsersWithEmailAddress(DB *sql.DB, emailAddress string) int {
	countUsers := 0

	// See if user with this email address already exists
	row := DB.QueryRow("SELECT COUNT(*) FROM users where email_address = $1", emailAddress)

	err := row.Scan(&countUsers)

	if err != nil {
		panic(err)
	}

	return countUsers
}

// CreateUser in the DB
func CreateUser(DB *sql.DB, emailAddress string, password string) int {
	var lastInsertID int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	userPassword := string(hashedPassword)

	// Insert account into db
	err = DB.QueryRow("INSERT INTO users(email_address, password, date_created ) VALUES($1, $2, $3) returning userid;", emailAddress, userPassword, helpers.GetCurrentDateTime()).Scan(&lastInsertID)

	if err != nil {
		panic(err)
	}

	return lastInsertID
}

// UserValidLogin by email address and password
func UserValidLogin(DB *sql.DB, emailAddress string, password string) int {
	// See if user with this email address and password exists
	rows, err := DB.Query("SELECT * FROM users where email_address = $1", emailAddress)

	if err != nil {
		panic(err)
	}

	userExists := false

	var passwordToCheck string
	var emailToCheck string
	var userID int

	// Foreach user in db
	for rows.Next() {
		var idFromDB int
		var emailFromDB string
		var passwordFromDB string
		var DateCreated string

		err = rows.Scan(&idFromDB, &emailFromDB, &passwordFromDB, &DateCreated)

		if err != nil {
			panic(err)
		}

		userID = idFromDB
		emailToCheck = emailFromDB
		passwordToCheck = passwordFromDB
		userExists = true
	}

	userValidLogin := false

	if userExists && emailToCheck != "" && passwordToCheck != "" {
		if err = bcrypt.CompareHashAndPassword([]byte(passwordToCheck), []byte(password)); err != nil {
			userValidLogin = false
		} else {
			userValidLogin = true
		}
	}

	if userValidLogin == true {
		return userID
	}

	return 0
}
