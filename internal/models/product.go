package models

import "gorm.io/gorm"

// Product represents a product in the system
type Product struct {
	gorm.Model
	ProductID string  `gorm:"column:product_id;uniqueIndex;type:varchar(50)" json:"product_id"`
	Name      string  `gorm:"column:name;not null;type:varchar(255)" json:"name"`
	Category  string  `gorm:"column:category;not null;type:varchar(100)" json:"category"`
	UnitPrice float64 `gorm:"column:unit_price;not null;type:decimal(10,2)" json:"unit_price"`
}

func (Product) TableName() string {
	return "products"
}
