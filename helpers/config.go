package helpers

import (
	"encoding/json"
	"os"
)

// Configuration struct containing application's variables
type Configuration struct {
	DBHost        string `json:"DBHost"`
	DBPort        string `json:"DBPort"`
	DBUser        string `json:"DBUser"`
	DBPassword    string `json:"DBPassword"`
	DBName        string `json:"DBName"`
	Port          string `json:"Port"`
	JWTSigningKey string `json:"JWTSigningKey"`
	RestAPIPath   string `json:"RestAPIPath"`
	RestAPIURL    string `json:"RestAPIURL"`
}

// GetConfig - get application's variables
func GetConfig() Configuration {
	// Extract variable from config.json file
	file, err := os.Open("config.json")

	if err != nil {
		panic(err)
	}

	var Config Configuration

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)

	if err != nil {
		panic(err)
	}

	return Config
}
