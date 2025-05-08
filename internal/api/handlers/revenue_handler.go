package handlers

import (
	"fmt"
	"net/http"
	"time"

	"sales-analytics/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RevenueHandler struct {
	revenueService *services.RevenueService
	logger         *logrus.Logger
}

func NewRevenueHandler(revenueService *services.RevenueService, logger *logrus.Logger) *RevenueHandler {
	return &RevenueHandler{
		revenueService: revenueService,
		logger:         logger,
	}
}

// GetTotalRevenue handles the total revenue calculation
func (h *RevenueHandler) GetTotalRevenue(c *gin.Context) {
	startDate, endDate, err := h.getDateRange(c)
	if err != nil {
		return // Error response already handled in getDateRange
	}

	revenue, err := h.revenueService.GetTotalRevenue(startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get total revenue")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate total revenue",
		})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// GetRevenueByProduct handles revenue calculation by product
func (h *RevenueHandler) GetRevenueByProduct(c *gin.Context) {
	startDate, endDate, err := h.getDateRange(c)
	if err != nil {
		return // Error response already handled in getDateRange
	}

	revenue, err := h.revenueService.GetRevenueByProduct(startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get revenue by product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate revenue by product",
		})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// GetRevenueByCategory handles revenue calculation by category
func (h *RevenueHandler) GetRevenueByCategory(c *gin.Context) {
	startDate, endDate, err := h.getDateRange(c)
	if err != nil {
		return // Error response already handled in getDateRange
	}

	revenue, err := h.revenueService.GetRevenueByCategory(startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get revenue by category")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate revenue by category",
		})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// GetRevenueByRegion handles revenue calculation by region
func (h *RevenueHandler) GetRevenueByRegion(c *gin.Context) {
	startDate, endDate, err := h.getDateRange(c)
	if err != nil {
		return // Error response already handled in getDateRange
	}

	revenue, err := h.revenueService.GetRevenueByRegion(startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get revenue by region")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate revenue by region",
		})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// getDateRange extracts and validates date range from request
func (h *RevenueHandler) getDateRange(c *gin.Context) (time.Time, time.Time, error) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Check if dates are provided
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both start_date and end_date are required in format YYYY-MM-DD",
		})
		return time.Time{}, time.Time{}, fmt.Errorf("missing date parameters")
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		h.logger.WithError(err).WithField("start_date", startDateStr).Error("Invalid start date format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid start date '%s'. Date must be in format YYYY-MM-DD", startDateStr),
		})
		return time.Time{}, time.Time{}, err
	}

	// Parse end date
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		h.logger.WithError(err).WithField("end_date", endDateStr).Error("Invalid end date format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid end date '%s'. Date must be in format YYYY-MM-DD", endDateStr),
		})
		return time.Time{}, time.Time{}, err
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "End date cannot be before start date",
		})
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range")
	}

	return startDate, endDate, nil
}
