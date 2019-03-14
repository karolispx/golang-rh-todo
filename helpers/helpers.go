package helpers

import (
	"encoding/json"
	"net/http"
)

// RestAPIResponse message
type RestAPIResponse struct {
	Type     string      `json:"type"`
	Response interface{} `json:"response"`
}

// RestAPIRespond - process rest api response
func RestAPIRespond(w http.ResponseWriter, r *http.Request, response interface{}, responseType string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(statusCode)

	returnResponse := RestAPIResponse{
		Type:     responseType,
		Response: response,
	}

	json.NewEncoder(w).Encode(returnResponse)
}
