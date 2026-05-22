package main

import (
	"car-rentals/internal/car"
	"car-rentals/internal/pb"
	"car-rentals/internal/platform"

	"google.golang.org/grpc"
)

func main() {
	db := platform.OpenDB()
	cache := platform.Redis()
	platform.ServeGRPC(":50052", func(s *grpc.Server) {
		pb.RegisterCarServiceServer(s, car.New(db, cache))
	})
}
