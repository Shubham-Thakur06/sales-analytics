package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"sales-analytics/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoadStatus struct {
	IsLoading    bool      `json:"is_loading"`
	StartTime    time.Time `json:"start_time,omitempty"`
	RecordsRead  int       `json:"records_read"`
	LastError    string    `json:"last_error,omitempty"`
	LastComplete time.Time `json:"last_complete,omitempty"`
}

type LoaderService struct {
	db          *gorm.DB
	logger      *logrus.Logger
	loadingLock sync.Mutex
	status      LoadStatus
	batchSize   int
}

func NewLoaderService(db *gorm.DB, logger *logrus.Logger, batchSize int) *LoaderService {
	return &LoaderService{
		db:        db,
		logger:    logger,
		batchSize: batchSize,
	}
}

// GetStatus returns the current loading status
func (s *LoaderService) GetStatus() LoadStatus {
	s.loadingLock.Lock()
	defer s.loadingLock.Unlock()
	return s.status
}

// IsLoading returns whether a data load is currently in progress
func (s *LoaderService) IsLoading() bool {
	s.loadingLock.Lock()
	defer s.loadingLock.Unlock()
	return s.status.IsLoading
}

// LoadData initiates the data loading process in the background
func (s *LoaderService) LoadData(csvPath string) error {
	s.loadingLock.Lock()
	if s.status.IsLoading {
		s.loadingLock.Unlock()
		return fmt.Errorf("data refresh is already in progress, please try again later")
	}
	s.status = LoadStatus{
		IsLoading:   true,
		StartTime:   time.Now(),
		RecordsRead: 0,
	}
	s.loadingLock.Unlock()

	// Start the loading process in a goroutine
	go func() {
		if err := s.processCSV(csvPath); err != nil {
			s.loadingLock.Lock()
			s.status.LastError = err.Error()
			s.status.IsLoading = false
			s.loadingLock.Unlock()
			s.logger.Errorf("Error loading data: %v", err)
			return
		}

		s.loadingLock.Lock()
		s.status.IsLoading = false
		s.status.LastComplete = time.Now()
		s.status.LastError = ""
		s.loadingLock.Unlock()
		s.logger.Info("Data loaded successfully")
	}()

	return nil
}

// processCSV handles the actual CSV processing
func (s *LoaderService) processCSV(csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %v", err)
	}

	var (
		customerMap = make(map[string]models.Customer)
		productMap  = make(map[string]models.Product)
		orders      []models.Order
		record      []string
		recordCount int
	)

	for {
		record, err = reader.Read()
		if err != nil {
			break // End of file
		}

		recordCount++
		s.loadingLock.Lock()
		s.status.RecordsRead = recordCount
		s.loadingLock.Unlock()

		// Parse numeric values
		quantity, _ := strconv.Atoi(record[7])                // Quantity Sold
		unitPrice, _ := strconv.ParseFloat(record[8], 64)     // Unit Price
		discount, _ := strconv.ParseFloat(record[9], 64)      // Discount
		shippingCost, _ := strconv.ParseFloat(record[10], 64) // Shipping Cost

		// Parse date
		date, err := time.Parse("2006-01-02", record[6]) // Date of Sale
		if err != nil {
			return fmt.Errorf("error parsing date: %v", err)
		}

		// Update customer map (last record wins for duplicates)
		customerMap[record[2]] = models.Customer{
			CustomerID: record[2],
			Name:       record[12],
			Email:      record[13],
			Address:    record[14],
			Region:     record[5],
		}

		// Update product map (last record wins for duplicates)
		productMap[record[1]] = models.Product{
			ProductID: record[1],
			Name:      record[3],
			Category:  record[4],
			UnitPrice: unitPrice,
		}

		orders = append(orders, models.Order{
			OrderID:       record[0],
			ProductID:     record[1],
			CustomerID:    record[2],
			DateOfSale:    date,
			Quantity:      quantity,
			Discount:      discount,
			ShippingCost:  shippingCost,
			PaymentMethod: record[11],
		})

		// Process in batches
		if len(orders) >= s.batchSize {
			if err := s.processBatch(mapToSlice(customerMap), mapToSlice(productMap), orders); err != nil {
				return err
			}
			orders = orders[:0]
			customerMap = make(map[string]models.Customer)
			productMap = make(map[string]models.Product)
		}
	}

	// Process remaining records
	if len(orders) > 0 {
		if err := s.processBatch(mapToSlice(customerMap), mapToSlice(productMap), orders); err != nil {
			return err
		}
	}

	return nil
}

// mapToSlice converts a map to a slice
func mapToSlice[T any](m map[string]T) []T {
	result := make([]T, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func (s *LoaderService) processBatch(customers []models.Customer, products []models.Product, orders []models.Order) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Batch upsert customers
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "customer_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "email", "address", "region"}),
		}).CreateInBatches(customers, s.batchSize).Error; err != nil {
			return fmt.Errorf("error upserting customers: %v", err)
		}

		// Batch upsert products
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "category", "unit_price"}),
		}).CreateInBatches(products, s.batchSize).Error; err != nil {
			return fmt.Errorf("error upserting products: %v", err)
		}

		// Batch insert orders (skip if exists)
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "order_id"}},
			DoNothing: true,
		}).CreateInBatches(orders, s.batchSize).Error; err != nil {
			return fmt.Errorf("error creating orders: %v", err)
		}

		return nil
	})
}
