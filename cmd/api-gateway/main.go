package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"car-rentals/internal/pb"
	"car-rentals/internal/platform"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type app struct {
	users    pb.UserServiceClient
	cars     pb.CarServiceClient
	rentals  pb.RentalServiceClient
	payments pb.PaymentServiceClient
	notices  pb.NotificationServiceClient
}

var httpRequestsTotal uint64

func main() {
	a := app{
		users:    pb.NewUserServiceClient(platform.Dial(platform.Env("USER_SERVICE_ADDR", "localhost:50051"))),
		cars:     pb.NewCarServiceClient(platform.Dial(platform.Env("CAR_SERVICE_ADDR", "localhost:50052"))),
		rentals:  pb.NewRentalServiceClient(platform.Dial(platform.Env("RENTAL_SERVICE_ADDR", "localhost:50053"))),
		payments: pb.NewPaymentServiceClient(platform.Dial(platform.Env("PAYMENT_SERVICE_ADDR", "localhost:50054"))),
		notices:  pb.NewNotificationServiceClient(platform.Dial(platform.Env("NOTIFICATION_SERVICE_ADDR", "localhost:50055"))),
	}
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"}}))
	r.Use(metricsMiddleware)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { write(w, map[string]string{"status": "ok"}) })
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintln(w, "# HELP car_rentals_http_requests_total Total HTTP requests handled by API Gateway.")
		fmt.Fprintln(w, "# TYPE car_rentals_http_requests_total counter")
		fmt.Fprintf(w, "car_rentals_http_requests_total %d\n", atomic.LoadUint64(&httpRequestsTotal))
	})

	r.Post("/api/users/register", handle(func(r *http.Request) (any, error) {
		var in pb.RegisterUserRequest
		decode(r, &in)
		return a.users.RegisterUser(r.Context(), &in)
	}))
	r.Post("/api/users/login", handle(func(r *http.Request) (any, error) {
		var in pb.LoginUserRequest
		decode(r, &in)
		return a.users.LoginUser(r.Context(), &in)
	}))
	r.Get("/api/users/{id}", handle(func(r *http.Request) (any, error) { return a.users.GetUserProfile(r.Context(), id(r)) }))
	r.Put("/api/users/{id}", handle(func(r *http.Request) (any, error) {
		var in pb.UpdateUserProfileRequest
		decode(r, &in)
		in.Id = chi.URLParam(r, "id")
		return a.users.UpdateUserProfile(r.Context(), &in)
	}))

	r.Post("/api/cars", handle(func(r *http.Request) (any, error) {
		var in pb.Car
		decode(r, &in)
		return a.cars.CreateCar(r.Context(), &in)
	}))
	r.Get("/api/cars", handle(func(r *http.Request) (any, error) {
		q := r.URL.Query()
		if len(q) > 0 {
			return a.cars.SearchCars(r.Context(), &pb.SearchCarsRequest{Brand: q.Get("brand"), Model: q.Get("model"), Type: q.Get("type"), Location: q.Get("location"), Availability: q.Get("availability")})
		}
		return a.cars.ListCars(r.Context(), &pb.Empty{})
	}))
	r.Get("/api/cars/available", handle(func(r *http.Request) (any, error) { return a.cars.ListAvailableCars(r.Context(), &pb.Empty{}) }))
	r.Get("/api/cars/{id}", handle(func(r *http.Request) (any, error) { return a.cars.GetCarById(r.Context(), id(r)) }))
	r.Put("/api/cars/{id}/status", handle(func(r *http.Request) (any, error) {
		var in pb.UpdateCarStatusRequest
		decode(r, &in)
		in.Id = chi.URLParam(r, "id")
		return a.cars.UpdateCarStatus(r.Context(), &in)
	}))
	r.Post("/api/cars/{id}/favorite", handle(func(r *http.Request) (any, error) {
		var in pb.FavoriteRequest
		decode(r, &in)
		in.CarId = chi.URLParam(r, "id")
		return a.cars.AddFavorite(r.Context(), &in)
	}))
	r.Get("/api/users/{id}/favorites", handle(func(r *http.Request) (any, error) { return a.cars.ListFavorites(r.Context(), id(r)) }))
	r.Post("/api/cars/{id}/reviews", handle(func(r *http.Request) (any, error) {
		var in pb.Review
		decode(r, &in)
		in.CarId = chi.URLParam(r, "id")
		return a.cars.AddReview(r.Context(), &in)
	}))

	r.Post("/api/rentals", handle(func(r *http.Request) (any, error) {
		var in pb.CreateRentalRequest
		decode(r, &in)
		return a.rentals.CreateRental(r.Context(), &in)
	}))
	r.Get("/api/rentals/{id}", handle(func(r *http.Request) (any, error) { return a.rentals.GetRentalById(r.Context(), id(r)) }))
	r.Get("/api/users/{id}/rentals", handle(func(r *http.Request) (any, error) { return a.rentals.ListUserRentals(r.Context(), id(r)) }))
	r.Post("/api/rentals/{id}/cancel", handle(func(r *http.Request) (any, error) { return a.rentals.CancelRental(r.Context(), id(r)) }))
	r.Post("/api/rentals/{id}/complete", handle(func(r *http.Request) (any, error) {
		var in pb.CompleteRentalRequest
		decode(r, &in)
		in.Id = chi.URLParam(r, "id")
		return a.rentals.CompleteRental(r.Context(), &in)
	}))
	r.Post("/api/rentals/price", handle(func(r *http.Request) (any, error) {
		var in pb.PriceRequest
		decode(r, &in)
		return a.rentals.CalculateRentalPrice(r.Context(), &in)
	}))

	r.Post("/api/payments", handle(func(r *http.Request) (any, error) {
		var in pb.ProcessPaymentRequest
		decode(r, &in)
		return a.payments.ProcessPayment(r.Context(), &in)
	}))
	r.Get("/api/payments/{id}", handle(func(r *http.Request) (any, error) { return a.payments.GetPaymentStatus(r.Context(), id(r)) }))
	r.Post("/api/payments/{id}/refund", handle(func(r *http.Request) (any, error) { return a.payments.RefundPayment(r.Context(), id(r)) }))
	r.Post("/api/notifications/email", handle(func(r *http.Request) (any, error) {
		var in pb.NotificationRequest
		decode(r, &in)
		return a.notices.SendNotificationEmail(r.Context(), &in)
	}))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func id(r *http.Request) *pb.GetByIdRequest { return &pb.GetByIdRequest{Id: chi.URLParam(r, "id")} }
func decode(r *http.Request, v any)         { _ = json.NewDecoder(r.Body).Decode(v) }
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		atomic.AddUint64(&httpRequestsTotal, 1)
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
func handle(fn func(*http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := fn(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		write(w, out)
	}
}
