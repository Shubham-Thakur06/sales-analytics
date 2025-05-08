package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	AppPort    int
	CSVPath    string
	CronSpec   string
	BatchSize  int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432 // default PostgreSQL port
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		appPort = 8080 // default application port
	}

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		batchSize = 1000 // default batch size
	}

	// Get database credentials from OS environment variables
	dbUser := os.Getenv("PG_DB_USER")
	if dbUser == "" {
		dbUser = "postgres" // default user if not set
	}

	dbPassword := os.Getenv("PG_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres" // default password if not set
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSL_MODE"),
		AppPort:    appPort,
		CSVPath:    os.Getenv("CSV_FILE_PATH"),
		CronSpec:   os.Getenv("REFRESH_CRON"),
		BatchSize:  batchSize,
	}, nil
}
