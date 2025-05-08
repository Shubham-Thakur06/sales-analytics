package container

import (
	"fmt"

	// "sales-analytics/internal/api"
	"sales-analytics/internal/config"
	"sales-analytics/internal/models"

	// "sales-analytics/internal/services"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Container holds all the dependencies for the application
type Container struct {
	Config *config.Config
	Logger *logrus.Logger
	DB     *gorm.DB
	Cron   *cron.Cron
}

// NewContainer initializes a new dependency container
func NewContainer(config *config.Config) (*Container, error) {
	container := &Container{}

	// Initialize logger
	container.Logger = logrus.New()
	container.Logger.SetFormatter(&logrus.JSONFormatter{})

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBSSLMode,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the database schemas
	if err := database.AutoMigrate(&models.Customer{}, &models.Product{}, &models.Order{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %v", err)
	}

	// Store config
	container.Config = config

	return container, nil
}
