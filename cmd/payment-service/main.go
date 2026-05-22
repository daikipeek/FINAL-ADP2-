package main

import (
	"car-rentals/internal/payment"
	"car-rentals/internal/pb"
	"car-rentals/internal/platform"

	"google.golang.org/grpc"
)

func main() {
	db := platform.OpenDB()
	nc := platform.NATS()
	platform.ServeGRPC(":50054", func(s *grpc.Server) {
		pb.RegisterPaymentServiceServer(s, payment.New(db, nc))
	})
}
