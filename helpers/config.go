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

	configExists := true

	if err != nil {
		configExists = false
	}

	var Config Configuration

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)

	if err != nil {
		configExists = false
	}

	// If config.json file doesn't exist, get config variables from environmental variables
	if configExists == false {
		if os.Getenv("DBHOST") != "" {
			Config.DBHost = os.Getenv("DBHOST")
		}

		if os.Getenv("DBPORT") != "" {
			Config.DBPort = os.Getenv("DBPORT")
		}

		if os.Getenv("DBUSER") != "" {
			Config.DBUser = os.Getenv("DBUSER")
		}

		if os.Getenv("DBPASSWORD") != "" {
			Config.DBPassword = os.Getenv("DBPASSWORD")
		}

		if os.Getenv("DBNAME") != "" {
			Config.DBName = os.Getenv("DBNAME")
		}

		if os.Getenv("PORT") != "" {
			Config.Port = os.Getenv("PORT")
		}

		if os.Getenv("JWTSIGNINGKEY") != "" {
			Config.JWTSigningKey = os.Getenv("JWTSIGNINGKEY")
		}

		if os.Getenv("RESTAPIPATH") != "" {
			Config.RestAPIPath = os.Getenv("RESTAPIPATH")
		}

		if os.Getenv("RESTAPIURL") != "" {
			Config.RestAPIURL = os.Getenv("RESTAPIURL")
		}
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
