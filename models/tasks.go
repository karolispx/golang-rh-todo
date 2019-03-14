package models

import (
	"database/sql"

	"github.com/karolispx/golang-rh-todo/helpers"
	_ "github.com/lib/pq"
)

// TaskDetails - tasks information
type TaskDetails struct {
	TaskID      int    `json:"taskid"`
	Task        string `json:"task"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

// TasksQueryParameters - parameters for tasks' query
type TasksQueryParameters struct {
	Limit   int
	Offset  int
	OrderBy string
	Order   string
}

// GetUserTasks - get user tasks
func GetUserTasks(userID int, DB *sql.DB, tasksQueryParameters TasksQueryParameters) []TaskDetails {
	// Get user tasks based on filters provided
	rows, err := DB.Query(
		"SELECT * FROM tasks WHERE userid = $1 ORDER BY "+tasksQueryParameters.OrderBy+" "+tasksQueryParameters.Order+" OFFSET $2 LIMIT $3",
		userID,
		tasksQueryParameters.Offset,
		tasksQueryParameters.Limit,
	)

	if err != nil {
		panic(err)
	}

	var tasks []TaskDetails

	// Foreach coin
	for rows.Next() {
		var TaskID int
		var UserID int
		var Task string
		var DateCreated string
		var DateUpdated string

		err = rows.Scan(&TaskID, &UserID, &Task, &DateCreated, &DateUpdated)

		if err != nil {
			panic(err)
		}

		tasks = append(tasks, TaskDetails{
			TaskID:      TaskID,
			Task:        Task,
			DateCreated: DateCreated,
			DateUpdated: DateUpdated,
		})
	}

	return tasks
}

// CreateUserTask - create user task
func CreateUserTask(userID int, DB *sql.DB, task string) int {
	var lastInsertID int

	// Insert user task db
	getCurrentDateTime := helpers.GetCurrentDateTime()

	err := DB.QueryRow("INSERT INTO tasks(userid, task, date_created, date_updated ) VALUES($1, $2, $3, $4) returning taskid;", userID, task, getCurrentDateTime, getCurrentDateTime).Scan(&lastInsertID)

	if err != nil {
		panic(err)
	}

	return lastInsertID
}