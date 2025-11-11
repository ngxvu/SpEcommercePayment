package models

import "github.com/google/uuid"

type Payment struct {
	BaseModel
	OrderID    uuid.UUID `gorm:"type:uuid;index"`
	CustomerID uuid.UUID `gorm:"type:uuid;index"`
	Amount     float64   `gorm:"type:numeric(10,2)"`
	Status     string    `gorm:"type:varchar(50);index"`
	Currency   string    `gorm:"type:varchar(10)"`
}

func (Payment) TableName() string {
	return "payments"
}
