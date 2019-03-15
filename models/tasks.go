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
	Watching    string `json:"watching"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

// TasksQueryParameters - parameters for tasks' query
type TasksQueryParameters struct {
	Limit    int
	Offset   int
	OrderBy  string
	Order    string
	Watching string
}

// GetUserTasks - get tasks.
func GetUserTasks(DB *sql.DB, userID int, tasksQueryParameters TasksQueryParameters) ([]TaskDetails, int) {
	countTasksReturned := 0

	var watchingQuery string

	if tasksQueryParameters.Watching != "" {
		watchingQuery = " AND watching = '" + tasksQueryParameters.Watching + "'"
	}

	// Get user tasks based on filters provided
	rows, err := DB.Query(
		"SELECT * FROM tasks WHERE userid = $1 "+watchingQuery+" ORDER BY "+tasksQueryParameters.OrderBy+" "+tasksQueryParameters.Order+" OFFSET $2 LIMIT $3",
		userID,
		tasksQueryParameters.Offset,
		tasksQueryParameters.Limit,
	)

	if err != nil {
		panic(err)
	}

	var tasks []TaskDetails

	// Foreach task
	for rows.Next() {
		var TaskID int
		var UserID int
		var Task string
		var Watching string
		var DateCreated string
		var DateUpdated string

		err = rows.Scan(&TaskID, &UserID, &Task, &Watching, &DateCreated, &DateUpdated)

		if err != nil {
			panic(err)
		}

		tasks = append(tasks, TaskDetails{
			TaskID:      TaskID,
			Task:        Task,
			Watching:    Watching,
			DateCreated: DateCreated,
			DateUpdated: DateUpdated,
		})

		countTasksReturned++
	}

	return tasks, countTasksReturned
}

// GetUserTask - get a task.
func GetUserTask(DB *sql.DB, taskid string, userID int) (TaskDetails, bool) {
	// Get user task
	rows, err := DB.Query(
		"SELECT * FROM tasks WHERE userid = $1 AND taskid = $2 ORDER BY taskid DESC LIMIT 1",
		userID,
		taskid,
	)

	if err != nil {
		panic(err)
	}

	taskExists := false
	var task TaskDetails

	// Foreach task - should only be max 1 in the list
	for rows.Next() {
		var TaskID int
		var UserID int
		var Task string
		var Watching string
		var DateCreated string
		var DateUpdated string

		err = rows.Scan(&TaskID, &UserID, &Task, &Watching, &DateCreated, &DateUpdated)

		if err != nil {
			panic(err)
		}

		task = TaskDetails{
			TaskID:      TaskID,
			Task:        Task,
			Watching:    Watching,
			DateCreated: DateCreated,
			DateUpdated: DateUpdated,
		}

		taskExists = true
	}

	return task, taskExists
}

// CreateUserTask - create a task.
func CreateUserTask(DB *sql.DB, userID int, task string) (TaskDetails, bool) {
	var lastInsertID int

	// Insert user task db
	getCurrentDateTime := helpers.GetCurrentDateTime()

	err := DB.QueryRow("INSERT INTO tasks(userid, task, watching, date_created, date_updated ) VALUES($1, $2, 'no', $3, $4) returning taskid;", userID, task, getCurrentDateTime, getCurrentDateTime).Scan(&lastInsertID)

	if err != nil {
		panic(err)
	}

	if lastInsertID < 1 {
		return TaskDetails{}, true
	}

	createdTask := TaskDetails{
		TaskID:      lastInsertID,
		Task:        task,
		Watching:    "no",
		DateCreated: getCurrentDateTime,
		DateUpdated: getCurrentDateTime,
	}

	return createdTask, false
}

// UpdateUserTask - update a task.
func UpdateUserTask(DB *sql.DB, taskid string, userID int, task string) (TaskDetails, bool) {
	lastUpdateID := 0
	watching := "no"

	getCurrentDateTime := helpers.GetCurrentDateTime()

	err := DB.QueryRow("UPDATE tasks SET task = $1, date_updated = $2 WHERE taskid = $3 AND userid = $4 returning taskid, watching;",
		task, getCurrentDateTime, taskid, userID).Scan(&lastUpdateID, &watching)

	if err != nil {
		panic(err)
	}

	if lastUpdateID < 1 {
		return TaskDetails{}, true
	}

	updatedTask := TaskDetails{
		TaskID:      lastUpdateID,
		Task:        task,
		Watching:    watching,
		DateCreated: getCurrentDateTime,
		DateUpdated: getCurrentDateTime,
	}

	return updatedTask, false
}

// DeleteUserTask - delete a task.
func DeleteUserTask(DB *sql.DB, taskid string, userID int) bool {
	_, err := DB.Exec("DELETE FROM tasks WHERE taskid = $1 AND userid = $2", taskid, userID)

	if err != nil {
		return false
	}

	return true
}

// DeleteUserTasks - delete all tasks.
func DeleteUserTasks(DB *sql.DB, userID int) bool {
	_, err := DB.Exec("DELETE FROM tasks WHERE userid = $1", userID)

	if err != nil {
		return false
	}

	return true
}

// WatchingUserTask - watch/unwatch a task.
func WatchingUserTask(DB *sql.DB, taskid string, userID int, watchingAction string) int {
	lastUpdatedID := 0

	if watchingAction != "no" && watchingAction != "yes" {
		return 0
	}

	err := DB.QueryRow("UPDATE tasks SET watching = $1, date_updated = $2 WHERE taskid = $3 AND userid = $4 returning taskid;",
		watchingAction, helpers.GetCurrentDateTime(), taskid, userID).Scan(&lastUpdatedID)

	if err != nil {
		panic(err)
	}

	return lastUpdatedID
}
