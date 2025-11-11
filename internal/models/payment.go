package models

import (
	"github.com/google/uuid"
	"payment/pkg/http/utils"
)

type Payment struct {
	BaseModel
	OrderID    uuid.UUID `gorm:"type:uuid;index"`
	CustomerID uuid.UUID `gorm:"type:uuid;index"`
	Amount     float64   `gorm:"type:numeric(10,2)"`
	Status     string    `gorm:"type:varchar(50);index"`
}

func (Payment) TableName() string {
	return "payments"
}

type CreatePaymentRequest struct {
	OrderID    uuid.UUID `json:"order_id" binding:"required"`
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Status     string    `json:"status" binding:"required"`
}

type CreatePaymentResponseData struct {
	PaymentID  uuid.UUID `json:"payment_id"`
	OrderID    uuid.UUID `json:"order_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
}

type CreatePaymentResponse struct {
	Meta *utils.MetaData           `json:"meta"`
	Data CreatePaymentResponseData `json:"data"`
}
