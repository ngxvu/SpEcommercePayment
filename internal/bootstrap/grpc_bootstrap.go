package bootstrap

import (
	"context"
	"fmt"
	"payment/internal/grpc/handlers"
	"payment/internal/grpc/server"
	"payment/internal/repositories"
	"payment/internal/services"
	"strconv"
)

func StartGRPC(app *App) (*server.GRPCServer, error) {

	grpcPort, err := strconv.Atoi(app.Config.GRPCPort)
	if err != nil || grpcPort == 0 {
		grpcPort = 50052
	}
	httpPort, err := strconv.Atoi(app.Config.HTTPPort)
	if err != nil || httpPort == 0 {
		httpPort = 8081
	}

	grpcAddr := fmt.Sprintf(":%d", grpcPort)
	httpAddr := fmt.Sprintf(":%d", httpPort)

	newPgRepo := app.PGRepo

	paymentRepo := repositories.NewPaymentRepository(newPgRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	handler := handlers.NewPaymentHandler(paymentService)

	grpcServer := server.NewGRPCServer(handler, grpcAddr, httpAddr)

	ctx := context.Background()

	go func() {
		if err := grpcServer.Run(ctx); err != nil {
			panic(err)
		}
	}()

	return grpcServer, nil
}
