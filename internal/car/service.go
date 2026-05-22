package car

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"car-rentals/internal/pb"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Service struct {
	db    *sql.DB
	cache *redis.Client
}

func New(db *sql.DB, cache *redis.Client) *Service { return &Service{db: db, cache: cache} }

func (s *Service) CreateCar(ctx context.Context, c *pb.Car) (*pb.Car, error) {
	if c.Id == "" {
		c.Id = uuid.NewString()
	}
	if c.Status == "" {
		c.Status = "available"
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `insert into cars(id,brand,model,type,location,status,daily_rate,image_url) values($1,$2,$3,$4,$5,$6,$7,$8)`,
		c.Id, c.Brand, c.Model, c.Type, c.Location, c.Status, c.DailyRate, c.ImageUrl); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.cache.Del(ctx, "cars:available")
	return c, nil
}

func (s *Service) GetCarById(ctx context.Context, r *pb.GetByIdRequest) (*pb.Car, error) {
	row := s.db.QueryRowContext(ctx, carSelect()+` where c.id=$1 group by c.id`, r.Id)
	return scanCar(row)
}

func (s *Service) ListCars(ctx context.Context, _ *pb.Empty) (*pb.CarList, error) {
	rows, err := s.db.QueryContext(ctx, carSelect()+` group by c.id order by c.created_at desc`)
	return scanCars(rows, err)
}

func (s *Service) ListAvailableCars(ctx context.Context, _ *pb.Empty) (*pb.CarList, error) {
	if raw, err := s.cache.Get(ctx, "cars:available").Bytes(); err == nil {
		var cached pb.CarList
		if json.Unmarshal(raw, &cached) == nil {
			return &cached, nil
		}
	}
	list, err := s.SearchCars(ctx, &pb.SearchCarsRequest{Availability: "available"})
	if err == nil {
		if raw, marshalErr := json.Marshal(list); marshalErr == nil {
			s.cache.Set(ctx, "cars:available", raw, 30*time.Second)
		}
	}
	return list, err
}

func (s *Service) UpdateCarStatus(ctx context.Context, r *pb.UpdateCarStatusRequest) (*pb.Car, error) {
	_, err := s.db.ExecContext(ctx, `update cars set status=$2,updated_at=now() where id=$1`, r.Id, r.Status)
	s.cache.Del(ctx, "cars:available")
	if err != nil {
		return nil, err
	}
	return s.GetCarById(ctx, &pb.GetByIdRequest{Id: r.Id})
}

func (s *Service) SearchCars(ctx context.Context, r *pb.SearchCarsRequest) (*pb.CarList, error) {
	var args []any
	var where []string
	add := func(field, val string) {
		if val != "" {
			args = append(args, "%"+strings.ToLower(val)+"%")
			where = append(where, "lower("+field+") like $"+itoa(len(args)))
		}
	}
	add("brand", r.Brand)
	add("model", r.Model)
	add("type", r.Type)
	add("location", r.Location)
	if r.Availability != "" {
		args = append(args, r.Availability)
		where = append(where, "status=$"+itoa(len(args)))
	}
	query := carSelect()
	if len(where) > 0 {
		query += " where " + strings.Join(where, " and ")
	}
	query += " group by c.id order by c.created_at desc"
	rows, err := s.db.QueryContext(ctx, query, args...)
	return scanCars(rows, err)
}

func (s *Service) AddFavorite(ctx context.Context, r *pb.FavoriteRequest) (*pb.Empty, error) {
	_, err := s.db.ExecContext(ctx, `insert into favorites(user_id,car_id) values($1,$2) on conflict do nothing`, r.UserId, r.CarId)
	return &pb.Empty{}, err
}

func (s *Service) ListFavorites(ctx context.Context, r *pb.GetByIdRequest) (*pb.CarList, error) {
	rows, err := s.db.QueryContext(ctx, carSelect()+` join favorites f on f.car_id=c.id where f.user_id=$1 group by c.id`, r.Id)
	return scanCars(rows, err)
}

func (s *Service) AddReview(ctx context.Context, r *pb.Review) (*pb.Review, error) {
	if r.Id == "" {
		r.Id = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `insert into reviews(id,user_id,car_id,rating,comment) values($1,$2,$3,$4,$5)`, r.Id, r.UserId, r.CarId, r.Rating, r.Comment)
	return r, err
}

func carSelect() string {
	return `select c.id,c.brand,c.model,c.type,c.location,c.status,c.daily_rate,c.image_url,coalesce(avg(r.rating),0) from cars c left join reviews r on r.car_id=c.id`
}

type scanner interface{ Scan(...any) error }

func scanCar(row scanner) (*pb.Car, error) {
	var c pb.Car
	err := row.Scan(&c.Id, &c.Brand, &c.Model, &c.Type, &c.Location, &c.Status, &c.DailyRate, &c.ImageUrl, &c.Rating)
	return &c, err
}
func scanCars(rows *sql.Rows, err error) (*pb.CarList, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := &pb.CarList{}
	for rows.Next() {
		c, err := scanCar(rows)
		if err != nil {
			return nil, err
		}
		list.Cars = append(list.Cars, c)
	}
	return list, rows.Err()
}
func itoa(n int) string { return string(rune('0' + n)) }
