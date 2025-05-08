package models

import "gorm.io/gorm"

// Customer represents a customer in the system
type Customer struct {
	gorm.Model
	CustomerID string `gorm:"column:customer_id;uniqueIndex;type:varchar(50)" json:"customer_id"`
	Name       string `gorm:"column:name;not null;type:varchar(255)" json:"name"`
	Email      string `gorm:"column:email;not null;type:varchar(255)" json:"email"`
	Address    string `gorm:"column:address;not null;type:text" json:"address"`
	Region     string `gorm:"column:region;not null;type:varchar(100)" json:"region"`
}

func (Customer) TableName() string {
	return "customers"
}
