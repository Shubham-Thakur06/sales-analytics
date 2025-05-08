package container

import (
	"fmt"

	"sales-analytics/internal/api"
	"sales-analytics/internal/config"
	"sales-analytics/internal/models"
	"sales-analytics/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Container holds all the dependencies for the application
type Container struct {
	Config         *config.Config
	Logger         *logrus.Logger
	DB             *gorm.DB
	Cron           *cron.Cron
	LoaderService  *services.LoaderService
	RevenueService *services.RevenueService
	Router         *api.Router
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

	// Store database connection
	container.DB = database

	// Auto-migrate the database schemas
	if err := database.AutoMigrate(&models.Customer{}, &models.Product{}, &models.Order{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %v", err)
	}

	// Store config
	container.Config = config

	// Initialize services
	container.LoaderService = services.NewLoaderService(database, container.Logger, config.BatchSize)
	container.RevenueService = services.NewRevenueService(database)

	// Initialize cron
	container.Cron = cron.New()
	if _, err := container.Cron.AddFunc(config.CronSpec, func() {
		if err := container.LoaderService.LoadData(config.CSVPath); err != nil {
			container.Logger.Errorf("Error in scheduled data refresh: %v", err)
		}
	}); err != nil {
		return nil, err
	}

	// Initialize router
	container.Router = api.NewRouter(
		container.LoaderService,
		container.RevenueService,
		container.Logger,
		config.CSVPath,
	)

	return container, nil
}

// Start starts all the background services
func (c *Container) Start() {
	c.Cron.Start()
}

// Stop gracefully stops all services
func (c *Container) Stop() {
	c.Cron.Stop()
	sqlDB, err := c.DB.DB()
	if err != nil {
		c.Logger.Errorf("Error getting underlying *sql.DB: %v", err)
		return
	}
	sqlDB.Close()
}

// SetupHTTPServer sets up the HTTP server with all routes
func (c *Container) SetupHTTPServer() *gin.Engine {
	router := gin.Default()
	c.Router.SetupRoutes(router)
	return router
}
