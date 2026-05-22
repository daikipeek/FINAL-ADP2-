package main

import (
	"car-rentals/internal/pb"
	"car-rentals/internal/platform"
	"car-rentals/internal/rental"

	"google.golang.org/grpc"
)

func main() {
	db := platform.OpenDB()
	carConn := platform.Dial(platform.Env("CAR_SERVICE_ADDR", "localhost:50052"))
	nc := platform.NATS()
	platform.ServeGRPC(":50053", func(s *grpc.Server) {
		pb.RegisterRentalServiceServer(s, rental.New(db, pb.NewCarServiceClient(carConn), nc))
	})
}
