package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	// TODO

	return nil
}
