package models

import (
	"database/sql"
	"time"

	"github.com/karolispx/golang-rh-todo/helpers"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// CountUsersWithEmailAddress - check to see if user with this email address exists already
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

// CreateUser - create user.
func CreateUser(DB *sql.DB, emailAddress string, password string) int {
	var lastInsertID int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	userPassword := string(hashedPassword)

	// Insert account into db
	getCurrentDateTime := helpers.GetCurrentDateTime()

	err = DB.QueryRow("INSERT INTO users(email_address, password, last_action, date_created ) VALUES($1, $2, $3, $4) returning userid;", emailAddress, userPassword, getCurrentDateTime, getCurrentDateTime).Scan(&lastInsertID)

	if err != nil {
		panic(err)
	}

	return lastInsertID
}

// UserValidLogin - check if user login is valid.
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
		var LastAction string
		var DateCreated string

		err = rows.Scan(&idFromDB, &emailFromDB, &passwordFromDB, &LastAction, &DateCreated)

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

// UpdateUserLastAction - update user last action.
func UpdateUserLastAction(DB *sql.DB, userID int) int {
	lastUpdatedID := 0

	err := DB.QueryRow("UPDATE users SET last_action = $1 WHERE userid = $2 returning userid;",
		helpers.GetCurrentDateTime(), userID).Scan(&lastUpdatedID)

	if err != nil {
		panic(err)
	}

	return lastUpdatedID
}

// UserNeedsCooldown - check to see if user is not trying to spam the system.
func UserNeedsCooldown(DB *sql.DB, userID int) bool {
	var lastAction string

	row := DB.QueryRow("SELECT last_action FROM users where userid = $1", userID)

	err := row.Scan(&lastAction)

	if err != nil {
		panic(err)
	}

	// Compare the date/time of last action and the current date/time to see if enough time has passed from last action
	if lastAction != "" {
		lastActionTime, err := time.Parse("2006.01.02 15:04:05", lastAction)

		if err == nil {
			// 3 seconds antispam protection
			cooldownTime := lastActionTime.Add(3 * time.Second)
			currentTime := time.Now()

			// timeLeft := currentTime.Sub(lastActionTime)
			timeLeft := cooldownTime.Sub(currentTime)

			if timeLeft > 0 {
				return true
			}
		}
	}

	return false
}
