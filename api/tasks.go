package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/karolispx/golang-rh-todo/helpers"
	"github.com/karolispx/golang-rh-todo/models"
)

// GetTasks - get user tasks.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Authenticate user to make sure user a valid user is attempting to view tasks.
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
