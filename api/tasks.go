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

// GetTasks - get user tasks.
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
			Limit:   10,
			Offset:  0,
			OrderBy: "taskid",
			Order:   "desc",
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

		// Get user tasks
		getUserTasks := models.GetUserTasks(userID, DB, tasksQueryParameters)

		helpers.RestAPIRespond(w, r, getUserTasks, "success", 200)

		return
	}
}

// CreateTask - create user task.
func CreateTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow creating a task.
	if userID > 0 {
		task := &models.TaskDetails{}

		err := json.NewDecoder(r.Body).Decode(task)

		if err != nil {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "error", 422)

			return
		}

		if task.Task == "" {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "error", 422)

			return
		}

		DB := helpers.InitDB()

		lastInsertID := models.CreateUserTask(userID, DB, task.Task)

		defer DB.Close()

		if lastInsertID < 1 {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		helpers.RestAPIRespond(w, r, lastInsertID, "success", 201)

		return
	}
}

// UpdateTask - update user task.
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow updating a task.
	if userID > 0 {
		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "error", 422)

			return
		}

		task := &models.TaskDetails{}

		err := json.NewDecoder(r.Body).Decode(task)

		if err != nil {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "error", 422)

			return
		}

		if task.Task == "" {
			helpers.RestAPIRespond(w, r, "Please provide a task.", "error", 422)

			return
		}

		DB := helpers.InitDB()

		// Check if this task belongs to the user
		checkTaskBelongsToUser := models.CheckTaskBelongsToUser(DB, taskid, userID)

		if checkTaskBelongsToUser < 1 {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "error", 422)

			return
		}

		updateTask := models.UpdateUserTask(DB, taskid, userID, task.Task)

		defer DB.Close()

		if updateTask < 1 {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		helpers.RestAPIRespond(w, r, "Task has been updated successfully!", "success", 200)

		return
	}
}

// DeleteTask - delete user task.
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user.
	userID := helpers.ValidateJWT(w, r)

	// If user is authenticated, allow deleting a task.
	if userID > 0 {
		vars := mux.Vars(r)
		taskid := vars["taskid"]

		if taskid == "" {
			helpers.RestAPIRespond(w, r, "Please provide task ID.", "error", 422)

			return
		}

		DB := helpers.InitDB()

		// Check if this task belongs to the user
		checkTaskBelongsToUser := models.CheckTaskBelongsToUser(DB, taskid, userID)

		if checkTaskBelongsToUser < 1 {
			helpers.RestAPIRespond(w, r, "This task does not belong to you!", "error", 422)

			return
		}

		deleteTask := models.DeleteUserTask(DB, taskid, userID)

		defer DB.Close()

		if deleteTask == false {
			helpers.DefaultErrorRestAPIRespond(w, r)

			return
		}

		helpers.RestAPIRespond(w, r, "Task has been deleted successfully!", "success", 204)

		return
	}
}
