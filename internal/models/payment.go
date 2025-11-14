package models

import (
	"github.com/google/uuid"
	"payment/pkg/http/utils"
)

type PaymentStatus string

const (
	// PaymentPending describes a payment that is created but not yet processed
	PaymentPending PaymentStatus = "PENDING"
	// PaymentAuthorized describes a payment that has been successfully authorized
	PaymentAuthorized PaymentStatus = "AUTHORIZED"
	// PaymentDeclined describes a payment that has been declined
	PaymentDeclined PaymentStatus = "DECLINED"
)

type Payment struct {
	BaseModel
	OrderID        string        `gorm:"index"`
	IdempotencyKey string        `gorm:"uniqueIndex:uniq_idem_key"`
	Amount         float64       `gorm:"not null"`
	Status         PaymentStatus `gorm:"type:text;index;not null"`
	Attempts       int           `gorm:"default:0"`
	LastError      string        `gorm:"type:text"`
}

func (Payment) TableName() string {
	return "payments"
}

type CreatePaymentRequest struct {
	OrderID        string  `json:"order_id" binding:"required"`
	IdempotencyKey string  `json:"idempotency_key" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
}

type CreatePaymentResponseData struct {
	PaymentID      uuid.UUID `json:"payment_id"`
	OrderID        string    `json:"order_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
}

type CreatePaymentResponse struct {
	Meta *utils.MetaData           `json:"meta"`
	Data CreatePaymentResponseData `json:"data"`
}

func NewCreatePaymentResponse(p *Payment, md *utils.MetaData) *CreatePaymentResponse {
	return &CreatePaymentResponse{
		Meta: md,
		Data: CreatePaymentResponseData{
			PaymentID:      p.ID,
			OrderID:        p.OrderID,
			IdempotencyKey: p.IdempotencyKey,
			Amount:         p.Amount,
			Status:         string(p.Status),
		},
	}
}
