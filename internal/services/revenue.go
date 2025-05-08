package services

import (
	"fmt"
	"time"

	"sales-analytics/internal/models"

	"gorm.io/gorm"
)

type RevenueService struct {
	db *gorm.DB
}

func NewRevenueService(db *gorm.DB) *RevenueService {
	return &RevenueService{db: db}
}

func (s *RevenueService) GetTotalRevenue(startDate, endDate time.Time) (*models.RevenueResponse, error) {
	var totalRevenue float64

	err := s.db.Model(&models.Order{}).
		Joins("JOIN products ON products.product_id = orders.product_id").
		Where("orders.date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM((products.unit_price * orders.quantity) - orders.discount + orders.shipping_cost), 0) as total_revenue").
		Scan(&totalRevenue).Error

	if err != nil {
		return nil, fmt.Errorf("error calculating total revenue: %v", err)
	}

	return &models.RevenueResponse{TotalRevenue: totalRevenue}, nil
}

func (s *RevenueService) GetRevenueByProduct(startDate, endDate time.Time) ([]models.ProductRevenue, error) {
	var results []models.ProductRevenue

	err := s.db.Model(&models.Product{}).
		Select("products.product_id, products.name as product_name, COALESCE(SUM((products.unit_price * orders.quantity) - orders.discount + orders.shipping_cost), 0) as revenue").
		Joins("LEFT JOIN orders ON products.product_id = orders.product_id AND orders.date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Group("products.product_id, products.name").
		Order("revenue DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error querying revenue by product: %v", err)
	}

	return results, nil
}

func (s *RevenueService) GetRevenueByCategory(startDate, endDate time.Time) ([]models.CategoryRevenue, error) {
	var results []models.CategoryRevenue

	err := s.db.Model(&models.Product{}).
		Select("products.category, COALESCE(SUM((products.unit_price * orders.quantity) - orders.discount + orders.shipping_cost), 0) as revenue").
		Joins("LEFT JOIN orders ON products.product_id = orders.product_id AND orders.date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Group("products.category").
		Order("revenue DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error querying revenue by category: %v", err)
	}

	return results, nil
}

func (s *RevenueService) GetRevenueByRegion(startDate, endDate time.Time) ([]models.RegionRevenue, error) {
	var results []models.RegionRevenue

	err := s.db.Model(&models.Customer{}).
		Select("customers.region, COALESCE(SUM((products.unit_price * orders.quantity) - orders.discount + orders.shipping_cost), 0) as revenue").
		Joins("LEFT JOIN orders ON customers.customer_id = orders.customer_id AND orders.date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Joins("LEFT JOIN products ON products.product_id = orders.product_id").
		Group("customers.region").
		Order("revenue DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error querying revenue by region: %v", err)
	}

	return results, nil
}
