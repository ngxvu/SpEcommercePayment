package services

import (
	"context"
	model "payment/internal/models"
	"payment/internal/repositories"
	"payment/pkg/core/logger"
	"payment/pkg/http/utils/app_errors"
	pb "payment/pkg/proto/paymentpb"
	"time"
)

type PaymentProcessor interface {
	Process(ctx context.Context, req *pb.PayRequest) (*pb.PayResponse, error)
}

type PaymentService struct {
	repo repositories.PaymentRepoInterface
}

func NewPaymentService(repo repositories.PaymentRepoInterface) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Process(ctx context.Context, req *pb.PayRequest) (*pb.PayResponse, error) {

	log := logger.WithTag("PaymentService|Process")

	if req == nil {
		err := app_errors.AppError(app_errors.StatusInternalServerError, app_errors.StatusInternalServerError)
		logger.LogError(log, err, "nil request")
		return nil, err
	}

	if req.Amount <= 0 {
		err := app_errors.AppError(app_errors.StatusInternalServerError, app_errors.StatusInternalServerError)
		logger.LogError(log, err, "invalid amount")
		return nil, err
	}
	if req.EventId == "" {
		err := app_errors.AppError(app_errors.StatusInternalServerError, app_errors.StatusInternalServerError)
		logger.LogError(log, err, "missing event ID for idempotency")
		return nil, err
	}

	createReq := &model.CreatePaymentRequest{
		OrderID:        req.OrderId,
		IdempotencyKey: req.EventId,
		Amount:         int64(req.Amount),
	}

	payment, created, err := s.repo.CreateOrGetPayment(ctx, createReq)
	if err != nil {
		logger.LogError(log, err, "failed to create or get payment")
		return nil, app_errors.AppError(app_errors.StatusInternalServerError, app_errors.StatusInternalServerError)

	}

	// Nếu đã tồn tại và đã có trạng thái cuối thì trả ngay
	if !created && (payment.Status == model.PaymentAuthorized || payment.Status == model.PaymentDeclined) {
		return &pb.PayResponse{
			Message:   string(payment.Status),
			PaymentId: payment.ID.String(),
			Status:    string(payment.Status),
		}, nil
	}

	// Chỉ thực hiện gateway nếu vừa tạo hoặc vẫn ở trạng thái PENDING
	if payment.Status == model.PaymentPending {
		select {
		case <-ctx.Done():
			_ = s.repo.UpdateStatus(ctx, payment.ID.String(), model.PaymentDeclined, "context canceled")
			return &pb.PayResponse{
				Message:   string(payment.Status),
				PaymentId: payment.ID.String(),
				Status:    string(model.PaymentDeclined),
			}, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}

		// Fake gateway success logic (có thể thay bằng gọi thật)
		gatewaySuccess := true

		if gatewaySuccess {
			_ = s.repo.UpdateStatus(ctx, payment.ID.String(), model.PaymentAuthorized, "")
			payment.Status = model.PaymentAuthorized
		} else {
			_ = s.repo.UpdateStatus(ctx, payment.ID.String(), model.PaymentDeclined, "gateway fail")
			payment.Status = model.PaymentDeclined
		}
	}

	return &pb.PayResponse{
		Message:   string(model.PaymentAuthorized),
		PaymentId: payment.ID.String(),
		Status:    string(payment.Status),
	}, nil
}

// Helpers (cần tự cài đặt parseUUID)
func parseUUID(s string) (u [16]byte) {
	return
}
