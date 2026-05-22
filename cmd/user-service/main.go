package main

import (
	"car-rentals/internal/pb"
	"car-rentals/internal/platform"
	"car-rentals/internal/user"

	"google.golang.org/grpc"
)

func main() {
	db := platform.OpenDB()
	platform.ServeGRPC(":50051", func(s *grpc.Server) {
		pb.RegisterUserServiceServer(s, user.New(db))
	})
}
