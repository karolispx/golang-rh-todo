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
