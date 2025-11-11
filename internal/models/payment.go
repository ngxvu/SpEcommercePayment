package models

import (
	"github.com/google/uuid"
	"payment/pkg/http/utils"
)

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "PENDING"
	PaymentAuthorized PaymentStatus = "AUTHORIZED"
	PaymentDeclined   PaymentStatus = "DECLINED"
)

type Payment struct {
	BaseModel
	OrderID        string        `gorm:"index"`
	IdempotencyKey string        `gorm:"uniqueIndex:uniq_idem_key"`
	Amount         int64         `gorm:"not null"`
	Status         PaymentStatus `gorm:"type:text;index;not null"`
	Attempts       int           `gorm:"default:0"`
	LastError      string        `gorm:"type:text"`
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
