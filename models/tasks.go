package models

import (
	"database/sql"

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
	Sort    string
	Keyword string
}

// GetUserTasks - get user tasks
func GetUserTasks(userID int, DB *sql.DB, tasksQueryParameters TasksQueryParameters) []TaskDetails {
	// Get tasks for this user
	rows, err := DB.Query("SELECT * FROM tasks WHERE userid = $1", userID)

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
