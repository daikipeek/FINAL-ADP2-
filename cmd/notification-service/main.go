package main

import (
	"car-rentals/internal/notification"
	"car-rentals/internal/pb"
	"car-rentals/internal/platform"

	"google.golang.org/grpc"
)

func main() {
	nc := platform.NATS()
	platform.ServeGRPC(":50055", func(s *grpc.Server) {
		pb.RegisterNotificationServiceServer(s, notification.New(nc))
	})
}
