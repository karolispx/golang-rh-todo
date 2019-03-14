package api

import (
	"net/http"

	"github.com/karolispx/golang-rh-todo/helpers"
)

// TestRoute will return 'Hello world!'
func TestRoute(w http.ResponseWriter, r *http.Request) {
	helpers.RestAPIRespond(w, r, "Hello world!", "success", 200)
}
