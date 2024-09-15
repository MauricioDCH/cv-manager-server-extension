package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser                 string
	DBPassword             string
	DBName                 string
	InstanceConnectionName string
	PrivateIP              string
	MasterKey              string
	GeminiAPIKey           string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		InstanceConnectionName: os.Getenv("INSTANCE_CONNECTION_NAME"),
		PrivateIP:              os.Getenv("PRIVATE_IP"),
		MasterKey:              os.Getenv("MASTER_KEY"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
	}

	// Validar si las variables esenciales están presentes
	if config.DBUser == "" || config.DBPassword == "" || config.DBName == "" || config.InstanceConnectionName == "" {
		return nil, fmt.Errorf("one or more required environment variables are missing or empty")
	}

	return config, nil
}
