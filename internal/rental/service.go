package rental

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"car-rentals/internal/pb"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const dateLayout = "2006-01-02"

type Service struct {
	db   *sql.DB
	cars pb.CarServiceClient
	nats *nats.Conn
}

func New(db *sql.DB, cars pb.CarServiceClient, nc *nats.Conn) *Service {
	return &Service{db: db, cars: cars, nats: nc}
}

func (s *Service) CreateRental(ctx context.Context, r *pb.CreateRentalRequest) (*pb.Rental, error) {
	car, err := s.cars.GetCarById(ctx, &pb.GetByIdRequest{Id: r.CarId})
	if err != nil {
		return nil, err
	}
	if car.Status != "available" {
		return nil, errors.New("car is not available")
	}
	price, err := calculate(car.DailyRate, r.StartDate, r.EndDate, r.EndDate, r.Insurance)
	if err != nil {
		return nil, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	out := &pb.Rental{Id: uuid.NewString(), UserId: r.UserId, CarId: r.CarId, StartDate: r.StartDate, EndDate: r.EndDate, Insurance: r.Insurance, TotalPrice: price.Total, Status: "confirmed"}
	_, err = tx.ExecContext(ctx, `insert into rentals(id,user_id,car_id,start_date,end_date,insurance,total_price,status) values($1,$2,$3,$4,$5,$6,$7,$8)`,
		out.Id, out.UserId, out.CarId, out.StartDate, out.EndDate, out.Insurance, out.TotalPrice, out.Status)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	_, _ = s.cars.UpdateCarStatus(ctx, &pb.UpdateCarStatusRequest{Id: r.CarId, Status: "rented"})
	s.publish("rental.created", out)
	return out, nil
}

func (s *Service) GetRentalById(ctx context.Context, r *pb.GetByIdRequest) (*pb.Rental, error) {
	row := s.db.QueryRowContext(ctx, rentalSelect()+` where id=$1`, r.Id)
	return scanRental(row)
}

func (s *Service) ListUserRentals(ctx context.Context, r *pb.GetByIdRequest) (*pb.RentalList, error) {
	rows, err := s.db.QueryContext(ctx, rentalSelect()+` where user_id=$1 order by created_at desc`, r.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := &pb.RentalList{}
	for rows.Next() {
		item, err := scanRental(rows)
		if err != nil {
			return nil, err
		}
		list.Rentals = append(list.Rentals, item)
	}
	return list, rows.Err()
}

func (s *Service) CancelRental(ctx context.Context, r *pb.GetByIdRequest) (*pb.Rental, error) {
	item, err := s.setStatus(ctx, r.Id, "cancelled", 0)
	if err == nil {
		_, _ = s.cars.UpdateCarStatus(ctx, &pb.UpdateCarStatusRequest{Id: item.CarId, Status: "available"})
		s.publish("rental.cancelled", item)
	}
	return item, err
}

func (s *Service) CompleteRental(ctx context.Context, r *pb.CompleteRentalRequest) (*pb.Rental, error) {
	item, err := s.GetRentalById(ctx, &pb.GetByIdRequest{Id: r.Id})
	if err != nil {
		return nil, err
	}
	car, err := s.cars.GetCarById(ctx, &pb.GetByIdRequest{Id: item.CarId})
	if err != nil {
		return nil, err
	}
	price, err := calculate(car.DailyRate, item.StartDate, item.EndDate, r.ReturnDate, item.Insurance)
	if err != nil {
		return nil, err
	}
	item, err = s.setStatus(ctx, r.Id, "completed", price.LateFee)
	if err == nil {
		_, _ = s.cars.UpdateCarStatus(ctx, &pb.UpdateCarStatusRequest{Id: item.CarId, Status: "available"})
	}
	return item, err
}

func (s *Service) CalculateRentalPrice(ctx context.Context, r *pb.PriceRequest) (*pb.PriceResponse, error) {
	car, err := s.cars.GetCarById(ctx, &pb.GetByIdRequest{Id: r.CarId})
	if err != nil {
		return nil, err
	}
	return calculate(car.DailyRate, r.StartDate, r.EndDate, r.ReturnDate, r.Insurance)
}

func (s *Service) setStatus(ctx context.Context, id, status string, lateFee float64) (*pb.Rental, error) {
	row := s.db.QueryRowContext(ctx, `update rentals set status=$2,late_fee=$3,updated_at=now() where id=$1 returning id,user_id,car_id,start_date::text,end_date::text,insurance,total_price,status,late_fee`, id, status, lateFee)
	return scanRental(row)
}

func calculate(rate float64, start, end, returned string, insurance bool) (*pb.PriceResponse, error) {
	s, err := time.Parse(dateLayout, start)
	if err != nil {
		return nil, err
	}
	e, err := time.Parse(dateLayout, end)
	if err != nil {
		return nil, err
	}
	if e.Before(s) {
		return nil, errors.New("end date must be after start date")
	}
	days := int(e.Sub(s).Hours()/24) + 1
	total := float64(days) * rate
	if insurance {
		total += float64(days) * 15
	}
	lateFee := 0.0
	if returned != "" {
		rt, err := time.Parse(dateLayout, returned)
		if err != nil {
			return nil, err
		}
		if rt.After(e) {
			lateFee = float64(int(rt.Sub(e).Hours()/24)) * rate * 0.5
			total += lateFee
		}
	}
	return &pb.PriceResponse{Total: total, LateFee: lateFee, Days: int32(days)}, nil
}

func (s *Service) publish(subject string, v any) {
	raw, _ := json.Marshal(v)
	_ = s.nats.Publish(subject, raw)
}

func rentalSelect() string {
	return `select id,user_id,car_id,start_date::text,end_date::text,insurance,total_price,status,late_fee from rentals`
}

type scanner interface{ Scan(...any) error }

func scanRental(row scanner) (*pb.Rental, error) {
	var r pb.Rental
	err := row.Scan(&r.Id, &r.UserId, &r.CarId, &r.StartDate, &r.EndDate, &r.Insurance, &r.TotalPrice, &r.Status, &r.LateFee)
	return &r, err
}
