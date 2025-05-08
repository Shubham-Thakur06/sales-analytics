package handlers

import (
	"net/http"

	"sales-analytics/internal/services"

	"github.com/gin-gonic/gin"
)

type RefreshHandler struct {
	loaderService *services.LoaderService
	csvFilePath   string
}

func NewRefreshHandler(loaderService *services.LoaderService, csvFilePath string) *RefreshHandler {
	return &RefreshHandler{
		loaderService: loaderService,
		csvFilePath:   csvFilePath,
	}
}

// RefreshData handles the data refresh operation
func (h *RefreshHandler) RefreshData(c *gin.Context) {
	if h.loaderService.IsLoading() {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Data refresh is already in progress, please try again later",
		})
		return
	}

	if err := h.loaderService.LoadData(h.csvFilePath); err != nil {
		if err.Error() == "data refresh is already in progress, please try again later" {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data refresh completed successfully",
	})
}
