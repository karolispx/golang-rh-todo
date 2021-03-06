package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// RestAPIResponse message
type RestAPIResponse struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Token information
type Token struct {
	UserID int
	jwt.StandardClaims
}

// RestAPIRespond - process rest api response
func RestAPIRespond(w http.ResponseWriter, r *http.Request, message string, data interface{}, responseType string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(status)

	returnResponse := RestAPIResponse{
		Type:    responseType,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(returnResponse)
}

// DefaultErrorRestAPIRespond - respond with default server 500 error
func DefaultErrorRestAPIRespond(w http.ResponseWriter, r *http.Request) {
	RestAPIRespond(w, r, "An error occured. Please try again later.", "", "error", 500)
}

// GetCurrentDateTime in string format
func GetCurrentDateTime() string {
	currentTime := time.Now()

	return currentTime.Format("2006.01.02 15:04:05")
}

// GenerateJWT - generating JWT
func GenerateJWT(userID int) string {
	Config := GetConfig()

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["UserID"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 10080).Unix() // 7 days - 10080

	tokenString, err := token.SignedString([]byte(Config.JWTSigningKey))

	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

// ValidateJWT - validate JWT
func ValidateJWT(w http.ResponseWriter, r *http.Request) int {
	tokenHeader := r.Header.Get("Authorization")

	if tokenHeader == "" {
		RestAPIRespond(w, r, "You need to be logged in to access this page.", "", "error", 403)

		return 0
	}

	splitted := strings.Split(tokenHeader, " ")

	if len(splitted) != 2 {
		RestAPIRespond(w, r, "You need to be logged in to access this page.", "", "error", 403)

		return 0
	}

	Config := GetConfig()

	tokenPart := splitted[1]

	headerToken := &Token{}

	token, err := jwt.ParseWithClaims(tokenPart, headerToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSigningKey), nil
	})

	// Malformed token, returns with http code 403 as usual
	if err != nil {
		RestAPIRespond(w, r, "You need to be logged in to access this page.", "", "error", 403)

		return 0
	}

	if !token.Valid {
		RestAPIRespond(w, r, "You need to be logged in to access this page.", "", "error", 403)

		return 0
	}

	return headerToken.UserID
}

// ValidateEmailAddress - validate email address
func ValidateEmailAddress(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	return Re.MatchString(email)
}
