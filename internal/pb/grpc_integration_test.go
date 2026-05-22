package pb

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestUserServiceGRPCIntegration(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	RegisterUserServiceServer(server, fakeUserService{})
	go func() {
		_ = server.Serve(lis)
	}()
	defer server.Stop()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	got, err := NewUserServiceClient(conn).RegisterUser(context.Background(), &RegisterUserRequest{Name: "Aida", Email: "aida@example.com", Role: "customer"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != "aida@example.com" || got.Role != "customer" {
		t.Fatalf("unexpected user response: %#v", got)
	}
}

type fakeUserService struct{}

func (fakeUserService) RegisterUser(_ context.Context, r *RegisterUserRequest) (*UserResponse, error) {
	return &UserResponse{Id: "user-1", Name: r.Name, Email: r.Email, Role: r.Role}, nil
}
func (fakeUserService) LoginUser(context.Context, *LoginUserRequest) (*LoginResponse, error) {
	return &LoginResponse{}, nil
}
func (fakeUserService) GetUserProfile(context.Context, *GetByIdRequest) (*UserResponse, error) {
	return &UserResponse{}, nil
}
func (fakeUserService) UpdateUserProfile(context.Context, *UpdateUserProfileRequest) (*UserResponse, error) {
	return &UserResponse{}, nil
}
