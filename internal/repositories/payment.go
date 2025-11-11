package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	CreateOrGetPayment(ctx context.Context, req *model.CreatePaymentRequest) (*model.Payment, bool, error)
	UpdateStatus(ctx context.Context, paymentID string, status model.PaymentStatus, lastErr string) error
}

func (r *PaymentRepository) CreateOrGetPayment(ctx context.Context, req *model.CreatePaymentRequest) (*model.Payment, bool, error) {
	tx, cancel := r.db.DBWithTimeout(ctx)
	defer cancel()

	p := &model.Payment{
		OrderID:        req.OrderID.String(),
		IdempotencyKey: req.IdempotencyKey,
		Amount:         req.Amount,
		Status:         model.PaymentPending,
	}

	// Try insert; do nothing on conflict (idempotency key)
	res := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idempotency_key"}},
		DoNothing: true,
	}).Create(p)

	if res.Error != nil {
		return nil, false, res.Error
	}

	created := res.RowsAffected == 1
	if !created {
		// Fetch existing
		var existing model.Payment
		if err := tx.Where("idempotency_key = ?", req.IdempotencyKey).First(&existing).Error; err != nil {
			return nil, false, err
		}
		// Validate consistency
		if existing.OrderID != req.OrderID.String() || existing.Amount != req.Amount {
			return nil, false, errors.New("idempotency key conflict: payload mismatch")
		}
		return &existing, false, nil
	}

	return p, true, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, paymentID string, status model.PaymentStatus, lastErr string) error {
	tx, cancel := r.db.DBWithTimeout(ctx)
	defer cancel()

	update := map[string]interface{}{
		"status": status,
	}
	if lastErr != "" {
		update["last_error"] = lastErr
		update["attempts"] = gorm.Expr("attempts + 1")
	}

	return tx.Model(&model.Payment{}).Where("id = ?", paymentID).Updates(update).Error
}

// Optional: build response helper
func BuildCreateResponse(ctx context.Context, p *model.Payment) *model.CreatePaymentResponse {
	return model.NewCreatePaymentResponse(p, utils.NewMetaData(ctx))
}
