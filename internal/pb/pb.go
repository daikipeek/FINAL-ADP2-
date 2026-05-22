package pb

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type jsonCodec struct{}

func (jsonCodec) Name() string { return "json" }
func (jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
func (jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type Empty struct{}
type GetByIdRequest struct {
	Id string `json:"id"`
}
type RegisterUserRequest struct{ Name, Email, Password, Role string }
type LoginUserRequest struct{ Email, Password string }
type UpdateUserProfileRequest struct{ Id, Name, Email string }
type UserResponse struct{ Id, Name, Email, Role string }
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type Car struct {
	Id, Brand, Model, Type, Location, Status, ImageUrl string
	DailyRate, Rating                                  float64
}
type CarList struct{ Cars []*Car }
type UpdateCarStatusRequest struct{ Id, Status string }
type SearchCarsRequest struct{ Brand, Model, Type, Location, Availability string }
type FavoriteRequest struct{ UserId, CarId string }
type Review struct {
	Id, UserId, CarId, Comment string
	Rating                     int32
}

type CreateRentalRequest struct {
	UserId, CarId, StartDate, EndDate string
	Insurance                         bool
}
type Rental struct {
	Id, UserId, CarId, StartDate, EndDate, Status string
	Insurance                                     bool
	TotalPrice, LateFee                           float64
}
type RentalList struct{ Rentals []*Rental }
type CompleteRentalRequest struct{ Id, ReturnDate string }
type PriceRequest struct {
	CarId, StartDate, EndDate, ReturnDate string
	Insurance                             bool
}
type PriceResponse struct {
	Total, LateFee float64
	Days           int32
}
type ProcessPaymentRequest struct {
	RentalId, UserId, Method string
	Amount                   float64
}
type Payment struct {
	Id, RentalId, UserId, Status, Method string
	Amount                               float64
}
type NotificationRequest struct{ To, Subject, Body string }

type UserServiceServer interface {
	RegisterUser(context.Context, *RegisterUserRequest) (*UserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginResponse, error)
	GetUserProfile(context.Context, *GetByIdRequest) (*UserResponse, error)
	UpdateUserProfile(context.Context, *UpdateUserProfileRequest) (*UserResponse, error)
}
type CarServiceServer interface {
	CreateCar(context.Context, *Car) (*Car, error)
	GetCarById(context.Context, *GetByIdRequest) (*Car, error)
	ListCars(context.Context, *Empty) (*CarList, error)
	ListAvailableCars(context.Context, *Empty) (*CarList, error)
	UpdateCarStatus(context.Context, *UpdateCarStatusRequest) (*Car, error)
	SearchCars(context.Context, *SearchCarsRequest) (*CarList, error)
	AddFavorite(context.Context, *FavoriteRequest) (*Empty, error)
	ListFavorites(context.Context, *GetByIdRequest) (*CarList, error)
	AddReview(context.Context, *Review) (*Review, error)
}
type RentalServiceServer interface {
	CreateRental(context.Context, *CreateRentalRequest) (*Rental, error)
	GetRentalById(context.Context, *GetByIdRequest) (*Rental, error)
	ListUserRentals(context.Context, *GetByIdRequest) (*RentalList, error)
	CancelRental(context.Context, *GetByIdRequest) (*Rental, error)
	CompleteRental(context.Context, *CompleteRentalRequest) (*Rental, error)
	CalculateRentalPrice(context.Context, *PriceRequest) (*PriceResponse, error)
}
type PaymentServiceServer interface {
	ProcessPayment(context.Context, *ProcessPaymentRequest) (*Payment, error)
	GetPaymentStatus(context.Context, *GetByIdRequest) (*Payment, error)
	RefundPayment(context.Context, *GetByIdRequest) (*Payment, error)
}
type NotificationServiceServer interface {
	SendNotificationEmail(context.Context, *NotificationRequest) (*Empty, error)
}

func unary[TReq any, TResp any](ctx context.Context, cc grpc.ClientConnInterface, method string, req *TReq) (*TResp, error) {
	var out TResp
	err := cc.Invoke(ctx, method, req, &out, grpc.ForceCodec(jsonCodec{}))
	return &out, err
}

type UserServiceClient interface {
	RegisterUser(context.Context, *RegisterUserRequest) (*UserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginResponse, error)
	GetUserProfile(context.Context, *GetByIdRequest) (*UserResponse, error)
	UpdateUserProfile(context.Context, *UpdateUserProfileRequest) (*UserResponse, error)
}
type userClient struct{ cc grpc.ClientConnInterface }

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient { return userClient{cc} }
func (c userClient) RegisterUser(ctx context.Context, r *RegisterUserRequest) (*UserResponse, error) {
	return unary[RegisterUserRequest, UserResponse](ctx, c.cc, "/carrentals.UserService/RegisterUser", r)
}
func (c userClient) LoginUser(ctx context.Context, r *LoginUserRequest) (*LoginResponse, error) {
	return unary[LoginUserRequest, LoginResponse](ctx, c.cc, "/carrentals.UserService/LoginUser", r)
}
func (c userClient) GetUserProfile(ctx context.Context, r *GetByIdRequest) (*UserResponse, error) {
	return unary[GetByIdRequest, UserResponse](ctx, c.cc, "/carrentals.UserService/GetUserProfile", r)
}
func (c userClient) UpdateUserProfile(ctx context.Context, r *UpdateUserProfileRequest) (*UserResponse, error) {
	return unary[UpdateUserProfileRequest, UserResponse](ctx, c.cc, "/carrentals.UserService/UpdateUserProfile", r)
}

type CarServiceClient interface {
	CreateCar(context.Context, *Car) (*Car, error)
	GetCarById(context.Context, *GetByIdRequest) (*Car, error)
	ListCars(context.Context, *Empty) (*CarList, error)
	ListAvailableCars(context.Context, *Empty) (*CarList, error)
	UpdateCarStatus(context.Context, *UpdateCarStatusRequest) (*Car, error)
	SearchCars(context.Context, *SearchCarsRequest) (*CarList, error)
	AddFavorite(context.Context, *FavoriteRequest) (*Empty, error)
	ListFavorites(context.Context, *GetByIdRequest) (*CarList, error)
	AddReview(context.Context, *Review) (*Review, error)
}
type carClient struct{ cc grpc.ClientConnInterface }

func NewCarServiceClient(cc grpc.ClientConnInterface) CarServiceClient { return carClient{cc} }
func (c carClient) CreateCar(ctx context.Context, r *Car) (*Car, error) {
	return unary[Car, Car](ctx, c.cc, "/carrentals.CarService/CreateCar", r)
}
func (c carClient) GetCarById(ctx context.Context, r *GetByIdRequest) (*Car, error) {
	return unary[GetByIdRequest, Car](ctx, c.cc, "/carrentals.CarService/GetCarById", r)
}
func (c carClient) ListCars(ctx context.Context, r *Empty) (*CarList, error) {
	return unary[Empty, CarList](ctx, c.cc, "/carrentals.CarService/ListCars", r)
}
func (c carClient) ListAvailableCars(ctx context.Context, r *Empty) (*CarList, error) {
	return unary[Empty, CarList](ctx, c.cc, "/carrentals.CarService/ListAvailableCars", r)
}
func (c carClient) UpdateCarStatus(ctx context.Context, r *UpdateCarStatusRequest) (*Car, error) {
	return unary[UpdateCarStatusRequest, Car](ctx, c.cc, "/carrentals.CarService/UpdateCarStatus", r)
}
func (c carClient) SearchCars(ctx context.Context, r *SearchCarsRequest) (*CarList, error) {
	return unary[SearchCarsRequest, CarList](ctx, c.cc, "/carrentals.CarService/SearchCars", r)
}
func (c carClient) AddFavorite(ctx context.Context, r *FavoriteRequest) (*Empty, error) {
	return unary[FavoriteRequest, Empty](ctx, c.cc, "/carrentals.CarService/AddFavorite", r)
}
func (c carClient) ListFavorites(ctx context.Context, r *GetByIdRequest) (*CarList, error) {
	return unary[GetByIdRequest, CarList](ctx, c.cc, "/carrentals.CarService/ListFavorites", r)
}
func (c carClient) AddReview(ctx context.Context, r *Review) (*Review, error) {
	return unary[Review, Review](ctx, c.cc, "/carrentals.CarService/AddReview", r)
}

type RentalServiceClient interface {
	CreateRental(context.Context, *CreateRentalRequest) (*Rental, error)
	GetRentalById(context.Context, *GetByIdRequest) (*Rental, error)
	ListUserRentals(context.Context, *GetByIdRequest) (*RentalList, error)
	CancelRental(context.Context, *GetByIdRequest) (*Rental, error)
	CompleteRental(context.Context, *CompleteRentalRequest) (*Rental, error)
	CalculateRentalPrice(context.Context, *PriceRequest) (*PriceResponse, error)
}
type rentalClient struct{ cc grpc.ClientConnInterface }

func NewRentalServiceClient(cc grpc.ClientConnInterface) RentalServiceClient { return rentalClient{cc} }
func (c rentalClient) CreateRental(ctx context.Context, r *CreateRentalRequest) (*Rental, error) {
	return unary[CreateRentalRequest, Rental](ctx, c.cc, "/carrentals.RentalService/CreateRental", r)
}
func (c rentalClient) GetRentalById(ctx context.Context, r *GetByIdRequest) (*Rental, error) {
	return unary[GetByIdRequest, Rental](ctx, c.cc, "/carrentals.RentalService/GetRentalById", r)
}
func (c rentalClient) ListUserRentals(ctx context.Context, r *GetByIdRequest) (*RentalList, error) {
	return unary[GetByIdRequest, RentalList](ctx, c.cc, "/carrentals.RentalService/ListUserRentals", r)
}
func (c rentalClient) CancelRental(ctx context.Context, r *GetByIdRequest) (*Rental, error) {
	return unary[GetByIdRequest, Rental](ctx, c.cc, "/carrentals.RentalService/CancelRental", r)
}
func (c rentalClient) CompleteRental(ctx context.Context, r *CompleteRentalRequest) (*Rental, error) {
	return unary[CompleteRentalRequest, Rental](ctx, c.cc, "/carrentals.RentalService/CompleteRental", r)
}
func (c rentalClient) CalculateRentalPrice(ctx context.Context, r *PriceRequest) (*PriceResponse, error) {
	return unary[PriceRequest, PriceResponse](ctx, c.cc, "/carrentals.RentalService/CalculateRentalPrice", r)
}

type PaymentServiceClient interface {
	ProcessPayment(context.Context, *ProcessPaymentRequest) (*Payment, error)
	GetPaymentStatus(context.Context, *GetByIdRequest) (*Payment, error)
	RefundPayment(context.Context, *GetByIdRequest) (*Payment, error)
}
type paymentClient struct{ cc grpc.ClientConnInterface }

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return paymentClient{cc}
}
func (c paymentClient) ProcessPayment(ctx context.Context, r *ProcessPaymentRequest) (*Payment, error) {
	return unary[ProcessPaymentRequest, Payment](ctx, c.cc, "/carrentals.PaymentService/ProcessPayment", r)
}
func (c paymentClient) GetPaymentStatus(ctx context.Context, r *GetByIdRequest) (*Payment, error) {
	return unary[GetByIdRequest, Payment](ctx, c.cc, "/carrentals.PaymentService/GetPaymentStatus", r)
}
func (c paymentClient) RefundPayment(ctx context.Context, r *GetByIdRequest) (*Payment, error) {
	return unary[GetByIdRequest, Payment](ctx, c.cc, "/carrentals.PaymentService/RefundPayment", r)
}

type NotificationServiceClient interface {
	SendNotificationEmail(context.Context, *NotificationRequest) (*Empty, error)
}
type notificationClient struct{ cc grpc.ClientConnInterface }

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return notificationClient{cc}
}
func (c notificationClient) SendNotificationEmail(ctx context.Context, r *NotificationRequest) (*Empty, error) {
	return unary[NotificationRequest, Empty](ctx, c.cc, "/carrentals.NotificationService/SendNotificationEmail", r)
}
