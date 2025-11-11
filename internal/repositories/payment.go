package repositories

import (
	"context"
	"gorm.io/gorm"
	model "payment/internal/models"
	pgGorm "payment/internal/repositories/pg-gorm"
	"payment/pkg/http/utils"
)

type PaymentRepository struct {
	db pgGorm.PGInterface
}

func NewPaymentRepository(newPgRepo pgGorm.PGInterface) *PaymentRepository {
	return &PaymentRepository{db: newPgRepo}
}

type PaymentRepoInterface interface {
	CreatePayment(ctx context.Context, tx *gorm.DB, PaymentRequest *model.CreatePaymentRequest) (*model.CreatePaymentResponse, error)
}

func (a *PaymentRepository) CreatePayment(ctx context.Context, tx *gorm.DB, paymentRequest *model.CreatePaymentRequest) (*model.CreatePaymentResponse, error) {

	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = a.db.DBWithTimeout(ctx)
		defer cancel()
	}

	paymentRecord := &model.Payment{
		OrderID:    paymentRequest.OrderID,
		CustomerID: paymentRequest.CustomerID,
		Amount:     paymentRequest.Amount,
		Status:     paymentRequest.Status,
	}

	if err := tx.Create(paymentRecord).Error; err != nil {
		return nil, err
	}

	response := &model.CreatePaymentResponse{
		Meta: utils.NewMetaData(ctx),
		Data: model.CreatePaymentResponseData{
			PaymentID:  paymentRecord.ID,
			OrderID:    paymentRecord.OrderID,
			CustomerID: paymentRecord.CustomerID,
			Amount:     paymentRecord.Amount,
			Status:     paymentRecord.Status,
		},
	}

	return response, nil
}
