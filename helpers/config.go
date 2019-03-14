package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

// InitDB - initialize DB
func InitDB() *sql.DB {
	// Get config
	var Config = GetConfig()

	// Initialize error variable
	var err error

	// Connection string
	DBInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config.DBHost, Config.DBPort, Config.DBUser, Config.DBPassword, Config.DBName,
	)

	// Open connection
	DB, err := sql.Open("postgres", DBInfo)

	// Panic if connection breaks or can not be opened
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}

	return DB
}
