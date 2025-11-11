package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	pb "payment/pkg/proto/paymentpb"
)

type GRPCServer struct {
	server     *grpc.Server
	grpcAddr   string
	httpAddr   string
	httpServer *http.Server
	lis        net.Listener
}

func NewGRPCServer(handler pb.PaymentServiceServer, grpcAddr, httpAddr string) *GRPCServer {
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, handler)
	return &GRPCServer{
		server:   s,
		grpcAddr: grpcAddr,
		httpAddr: httpAddr,
	}
}

// Run starts the gRPC server and the grpc-gateway HTTP server.
// Cancel the provided ctx to trigger graceful shutdown.
func (s *GRPCServer) Run(ctx context.Context) error {
	// start gRPC listener
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}
	s.lis = lis

	grpcErrCh := make(chan error, 1)
	go func() {
		log.Printf("gRPC server running on %s", s.grpcAddr)
		if err := s.server.Serve(lis); err != nil {
			grpcErrCh <- err
		}
	}()

	// wait for cancel or server error
	select {
	case <-ctx.Done():
		log.Println("shutting down servers")
		// Graceful shutdown of HTTP and gRPC
		_ = s.httpServer.Shutdown(context.Background())
		s.server.GracefulStop()
		return ctx.Err()
	case err := <-grpcErrCh:
		return fmt.Errorf("gRPC server error: %w", err)
	}
}

// Stop triggers an immediate graceful shutdown.
func (s *GRPCServer) Stop() {
	if s.httpServer != nil {
		_ = s.httpServer.Shutdown(context.Background())
	}
	if s.server != nil {
		s.server.GracefulStop()
	}
}
