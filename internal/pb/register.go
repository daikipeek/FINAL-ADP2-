package pb

import (
	"context"

	"google.golang.org/grpc"
)

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "carrentals.UserService", HandlerType: (*UserServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "RegisterUser", Handler: handler(srv.RegisterUser)},
		{MethodName: "LoginUser", Handler: handler(srv.LoginUser)},
		{MethodName: "GetUserProfile", Handler: handler(srv.GetUserProfile)},
		{MethodName: "UpdateUserProfile", Handler: handler(srv.UpdateUserProfile)},
	}}, srv)
}

func RegisterCarServiceServer(s grpc.ServiceRegistrar, srv CarServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "carrentals.CarService", HandlerType: (*CarServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "CreateCar", Handler: handler(srv.CreateCar)},
		{MethodName: "GetCarById", Handler: handler(srv.GetCarById)},
		{MethodName: "ListCars", Handler: handler(srv.ListCars)},
		{MethodName: "ListAvailableCars", Handler: handler(srv.ListAvailableCars)},
		{MethodName: "UpdateCarStatus", Handler: handler(srv.UpdateCarStatus)},
		{MethodName: "SearchCars", Handler: handler(srv.SearchCars)},
		{MethodName: "AddFavorite", Handler: handler(srv.AddFavorite)},
		{MethodName: "ListFavorites", Handler: handler(srv.ListFavorites)},
		{MethodName: "AddReview", Handler: handler(srv.AddReview)},
	}}, srv)
}

func RegisterRentalServiceServer(s grpc.ServiceRegistrar, srv RentalServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "carrentals.RentalService", HandlerType: (*RentalServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "CreateRental", Handler: handler(srv.CreateRental)},
		{MethodName: "GetRentalById", Handler: handler(srv.GetRentalById)},
		{MethodName: "ListUserRentals", Handler: handler(srv.ListUserRentals)},
		{MethodName: "CancelRental", Handler: handler(srv.CancelRental)},
		{MethodName: "CompleteRental", Handler: handler(srv.CompleteRental)},
		{MethodName: "CalculateRentalPrice", Handler: handler(srv.CalculateRentalPrice)},
	}}, srv)
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "carrentals.PaymentService", HandlerType: (*PaymentServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "ProcessPayment", Handler: handler(srv.ProcessPayment)},
		{MethodName: "GetPaymentStatus", Handler: handler(srv.GetPaymentStatus)},
		{MethodName: "RefundPayment", Handler: handler(srv.RefundPayment)},
	}}, srv)
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "carrentals.NotificationService", HandlerType: (*NotificationServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "SendNotificationEmail", Handler: handler(srv.SendNotificationEmail)},
	}}, srv)
}

func handler[TReq any, TResp any](fn func(context.Context, *TReq) (*TResp, error)) grpc.MethodHandler {
	return func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(TReq)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return fn(ctx, in)
		}
		info := &grpc.UnaryServerInfo{Server: srv}
		return interceptor(ctx, in, info, func(ctx context.Context, req any) (any, error) {
			return fn(ctx, req.(*TReq))
		})
	}
}
