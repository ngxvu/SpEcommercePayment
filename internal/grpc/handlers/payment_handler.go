package handlers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"payment/internal/services"
	pb "payment/pkg/proto/paymentpb"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	svc services.PaymentProcessor
}

func NewPaymentHandler(s services.PaymentProcessor) *PaymentHandler {
	return &PaymentHandler{svc: s}
}

func (h *PaymentHandler) Pay(ctx context.Context, req *pb.PayRequest) (*pb.PayResponse, error) {
	// validate request minimally
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is nil")
	}
	// call business logic
	resp, err := h.svc.Process(ctx, req)
	if err != nil {
		// translate or log as needed; here return response with error text
		if resp == nil {
			return nil, status.Errorf(codes.Internal, "payment processing failed: %v", err)
		}
		return resp, nil
	}
	return resp, nil
}
