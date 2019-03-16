package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/karolispx/golang-rh-todo/helpers"
	"github.com/karolispx/golang-rh-todo/models"
)

// GetTasks - get user tasks by query parameters.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, get user tasks.
	if userID > 0 {
		DB := helpers.InitDB()

		defer DB.Close()

		// Pagination/filtering/sorting/search
		// Default query parameters
		tasksQueryParameters := models.TasksQueryParameters{
			Limit:    10,
			Offset:   0,
			OrderBy:  "taskid",
			Order:    "desc",
			Watching: "",
		}

		// Get query vars from the request
		queryVariables := r.URL.Query()

		// Limit
		limitQueryVar := queryVariables.Get("limit")

		if limitQueryVar != "" {
			limitQueryVarInt, err := strconv.Atoi(limitQueryVar)

			if err == nil && limitQueryVarInt > 0 {
				tasksQueryParameters.Limit = limitQueryVarInt
			}
		}

		// Offset
		offsetQueryVar := queryVariables.Get("offset")

		if offsetQueryVar != "" {
			offsetQueryVarInt, err := strconv.Atoi(offsetQueryVar)

			if err == nil && offsetQueryVarInt > 0 {
				tasksQueryParameters.Offset = offsetQueryVarInt
			}
		}

		// OrderBY
		orderbyQueryVar := queryVariables.Get("orderby")

		if orderbyQueryVar != "" {
			// Convert to lowercase to avoid problems
			orderbyQueryVar = strings.ToLower(orderbyQueryVar)

			if orderbyQueryVar == "task" || orderbyQueryVar == "date_created" || orderbyQueryVar == "date_updated" {
				tasksQueryParameters.OrderBy = orderbyQueryVar
			}
		}

		// Order
		orderQueryVar := queryVariables.Get("order")

		if orderQueryVar != "" {
			// Convert to lowercase to avoid problems
			orderQueryVar = strings.ToLower(orderQueryVar)

			if orderQueryVar == "asc" || orderQueryVar == "desc" {
				tasksQueryParameters.Order = orderQueryVar
			}
		}

		// Watching
		watchingQueryVar := queryVariables.Get("watching")

		if watchingQueryVar != "" {
			// Convert to lowercase to avoid problems
			watchingQueryVar = strings.ToLower(watchingQueryVar)

			if watchingQueryVar == "yes" || watchingQueryVar == "no" {
				tasksQueryParameters.Watching = watchingQueryVar
			}
		}

		// Get user tasks
		userTasks, countTasksReturned := models.GetUserTasks(DB, userID, tasksQueryParameters)

		if countTasksReturned < 1 {
			helpers.RestAPIRespond(w, r, "No tasks available.", userTasks, "error", 422)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		helpers.RestAPIRespond(w, r, "", userTasks, "success", 200)

		return
	}
}

// GetTask - get user task by task ID.
func GetTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, get user task.
	if userID > 0 {
		DB := helpers.InitDB()

		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "", "error", 422)

			return
		}

		// Get user task
		userTask, taskExists := models.GetUserTask(DB, taskid, userID)

		// Check if this task belongs to the user
		if taskExists == false {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "", "error", 422)

			return
		}

		defer DB.Close()

		models.UpdateUserLastAction(DB, userID)

		helpers.RestAPIRespond(w, r, "", userTask, "success", 200)

		return
	}
}

// CreateTask - create user task.
func CreateTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow creating a task.
	if userID > 0 {
		DB := helpers.InitDB()

		// User cooldown period - antispam.
		userNeedsCooldown := models.UserNeedsCooldown(DB, userID)

		if userNeedsCooldown == true {
			helpers.RestAPIRespond(w, r, "Please slow down. You have been making too many requests.", "", "error", 422)

			return
		}

		task := &models.TaskDetails{}

		err := json.NewDecoder(r.Body).Decode(task)

		if err != nil {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "", "error", 422)

			return
		}

		if task.Task == "" {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "", "error", 422)

			return
		}

		createdTask, createError := models.CreateUserTask(DB, userID, task.Task)

		defer DB.Close()

		if createError == true {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		helpers.RestAPIRespond(w, r, "Task has been created successfully!", createdTask, "success", 201)

		return
	}
}

// UpdateTask - update user task by task ID.
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow updating a task.
	if userID > 0 {
		DB := helpers.InitDB()

		// User cooldown period - antispam.
		userNeedsCooldown := models.UserNeedsCooldown(DB, userID)

		if userNeedsCooldown == true {
			helpers.RestAPIRespond(w, r, "Please slow down. You have been making too many requests.", "", "error", 422)

			return
		}

		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "", "error", 422)

			return
		}

		task := &models.TaskDetails{}

		err := json.NewDecoder(r.Body).Decode(task)

		if err != nil {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "", "error", 422)

			return
		}

		if task.Task == "" {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "", "error", 422)

			return
		}

		// Get user task
		_, taskExists := models.GetUserTask(DB, taskid, userID)

		// Check if this task belongs to the user
		if taskExists == false {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "", "error", 422)

			return
		}

		updatedTask, updateError := models.UpdateUserTask(DB, taskid, userID, task.Task)

		defer DB.Close()

		if updateError == true {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		helpers.RestAPIRespond(w, r, "Task has been updated successfully!", updatedTask, "success", 200)

		return
	}
}

// DeleteTask - delete user task by task ID.
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow deleting a task.
	if userID > 0 {
		DB := helpers.InitDB()

		// User cooldown period - antispam.
		userNeedsCooldown := models.UserNeedsCooldown(DB, userID)

		if userNeedsCooldown == true {
			helpers.RestAPIRespond(w, r, "Please slow down. You have been making too many requests.", "", "error", 422)

			return
		}

		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "", "error", 422)

			return
		}

		// Get user task
		_, taskExists := models.GetUserTask(DB, taskid, userID)

		// Check if this task belongs to the user
		if taskExists == false {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "", "error", 422)

			return
		}

		deleteTask := models.DeleteUserTask(DB, taskid, userID)

		defer DB.Close()

		if deleteTask == false {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		helpers.RestAPIRespond(w, r, "Task has been deleted successfully!", "", "success", 200)

		return
	}
}

// DeleteTasks - delete all user tasks by user ID.
func DeleteTasks(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow deleting all tasks.
	if userID > 0 {
		DB := helpers.InitDB()

		// User cooldown period - antispam.
		userNeedsCooldown := models.UserNeedsCooldown(DB, userID)

		if userNeedsCooldown == true {
			helpers.RestAPIRespond(w, r, "Please slow down. You have been making too many requests.", "", "error", 422)

			return
		}

		deleteTasks := models.DeleteUserTasks(DB, userID)

		defer DB.Close()

		if deleteTasks == false {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		helpers.RestAPIRespond(w, r, "All Tasks have been deleted successfully!", "", "success", 200)

		return
	}
}

// WatchTask - watch a task by task ID.
func WatchTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow watching a task.
	if userID > 0 {
		DB := helpers.InitDB()

		// User cooldown period - antispam.
		userNeedsCooldown := models.UserNeedsCooldown(DB, userID)

		if userNeedsCooldown == true {
			helpers.RestAPIRespond(w, r, "Please slow down. You have been making too many requests.", "", "error", 422)

			return
		}

		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "error", "", 422)

			return
		}

		// Get user task
		userTask, taskExists := models.GetUserTask(DB, taskid, userID)

		// Check if this task belongs to the user
		if taskExists == false {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "", "error", 422)

			return
		}

		// Check if user is watching this task
		if userTask.Watching != "" && userTask.Watching == "yes" {
			helpers.RestAPIRespond(w, r, "You can not watch a task that you are already watching!", "", "error", 422)

			return
		}

		lastUpdatedID := models.WatchingUserTask(DB, taskid, userID, "yes")

		defer DB.Close()

		if lastUpdatedID < 1 {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		userTask.Watching = "yes"

		helpers.RestAPIRespond(w, r, "You are now watching this task!", userTask, "success", 201)

		return
	}
}

// UnwatchTask - unwatch a task by task ID.
func UnwatchTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow unwatching a task.
	if userID > 0 {
		DB := helpers.InitDB()

		// User cooldown period - antispam.
		userNeedsCooldown := models.UserNeedsCooldown(DB, userID)

		if userNeedsCooldown == true {
			helpers.RestAPIRespond(w, r, "Please slow down. You have been making too many requests.", "", "error", 422)

			return
		}

		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "", "error", 422)

			return
		}

		// Get user task
		userTask, taskExists := models.GetUserTask(DB, taskid, userID)

		// Check if this task belongs to the user
		if taskExists == false {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "", "error", 422)

			return
		}

		// Check if user is watching this task
		if userTask.Watching != "" && userTask.Watching == "no" {
			helpers.RestAPIRespond(w, r, "You can not unwatch a task that you are not watching!", "", "error", 422)

			return
		}

		lastUpdatedID := models.WatchingUserTask(DB, taskid, userID, "no")

		defer DB.Close()

		if lastUpdatedID < 1 {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		models.UpdateUserLastAction(DB, userID)

		userTask.Watching = "no"

		helpers.RestAPIRespond(w, r, "You are no longer watching this task!", userTask, "success", 201)

		return
	}
}
