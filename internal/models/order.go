package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a sales order in the system
type Order struct {
	gorm.Model
	OrderID       string    `gorm:"column:order_id;uniqueIndex;type:varchar(50)" json:"order_id"`
	CustomerID    string    `gorm:"column:customer_id;not null;type:varchar(50);index" json:"customer_id"`
	ProductID     string    `gorm:"column:product_id;not null;type:varchar(50);index" json:"product_id"`
	DateOfSale    time.Time `gorm:"column:date_of_sale;not null;index" json:"date_of_sale"`
	Quantity      int       `gorm:"column:quantity;not null" json:"quantity"`
	Discount      float64   `gorm:"column:discount;not null;type:decimal(10,2)" json:"discount"`
	ShippingCost  float64   `gorm:"column:shipping_cost;not null;type:decimal(10,2)" json:"shipping_cost"`
	PaymentMethod string    `gorm:"column:payment_method;not null;type:varchar(50)" json:"payment_method"`
}

func (Order) TableName() string {
	return "orders"
}
