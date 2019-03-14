package api

import (
	"net/http"

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

		// TODO: Pagination/filtering/sorting/search
		limit := 10
		offset := 10
		sort := "asc"
		keyword := ""

		tasksQueryParameters := models.TasksQueryParameters{
			Limit:   limit,
			Offset:  offset,
			Sort:    sort,
			Keyword: keyword,
		}

		// Get user tasks
		getUserTasks := models.GetUserTasks(userID, DB, tasksQueryParameters)

		helpers.RestAPIRespond(w, r, getUserTasks, "success", 200)

		return
	}
}
