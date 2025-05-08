package api

import (
	"sales-analytics/internal/api/handlers"

	"github.com/gin-gonic/gin"

	"sales-analytics/internal/services"

	"github.com/sirupsen/logrus"
)

type Router struct {
	refreshHandler *handlers.RefreshHandler
	revenueHandler *handlers.RevenueHandler
}

func NewRouter(loaderService *services.LoaderService, revenueService *services.RevenueService, logger *logrus.Logger, csvFilePath string) *Router {
	return &Router{
		refreshHandler: handlers.NewRefreshHandler(loaderService, csvFilePath),
		revenueHandler: handlers.NewRevenueHandler(revenueService, logger),
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Data refresh endpoint
		api.POST("/refresh", r.refreshHandler.RefreshData)

		// Revenue endpoints
		api.GET("/revenue", r.revenueHandler.GetTotalRevenue)
		api.GET("/revenue/product", r.revenueHandler.GetRevenueByProduct)
		api.GET("/revenue/category", r.revenueHandler.GetRevenueByCategory)
		api.GET("/revenue/region", r.revenueHandler.GetRevenueByRegion)
	}
}
