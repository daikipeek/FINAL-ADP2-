package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"car-rentals/internal/pb"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Service struct {
	db   *sql.DB
	nats *nats.Conn
}

func New(db *sql.DB, nc *nats.Conn) *Service { return &Service{db: db, nats: nc} }

func (s *Service) ProcessPayment(ctx context.Context, r *pb.ProcessPaymentRequest) (*pb.Payment, error) {
	status := "success"
	if r.Amount <= 0 || r.Method == "fail" {
		status = "failed"
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	p := &pb.Payment{Id: uuid.NewString(), RentalId: r.RentalId, UserId: r.UserId, Amount: r.Amount, Method: r.Method, Status: status}
	_, err = tx.ExecContext(ctx, `insert into payments(id,rental_id,user_id,amount,status,method) values($1,$2,$3,$4,$5,$6)`,
		p.Id, p.RentalId, p.UserId, p.Amount, p.Status, p.Method)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.publish("payment."+status, p)
	if status == "failed" {
		return p, errors.New("payment failed")
	}
	return p, nil
}

func (s *Service) GetPaymentStatus(ctx context.Context, r *pb.GetByIdRequest) (*pb.Payment, error) {
	return s.byID(ctx, r.Id)
}

func (s *Service) RefundPayment(ctx context.Context, r *pb.GetByIdRequest) (*pb.Payment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var p pb.Payment
	err = tx.QueryRowContext(ctx, `update payments set status='refunded',updated_at=now() where id=$1 returning id,rental_id,user_id,amount,status,method`, r.Id).
		Scan(&p.Id, &p.RentalId, &p.UserId, &p.Amount, &p.Status, &p.Method)
	if err != nil {
		return nil, err
	}
	return &p, tx.Commit()
}

func (s *Service) byID(ctx context.Context, id string) (*pb.Payment, error) {
	var p pb.Payment
	err := s.db.QueryRowContext(ctx, `select id,rental_id,user_id,amount,status,method from payments where id=$1`, id).
		Scan(&p.Id, &p.RentalId, &p.UserId, &p.Amount, &p.Status, &p.Method)
	return &p, err
}

func (s *Service) publish(subject string, v any) {
	raw, _ := json.Marshal(v)
	_ = s.nats.Publish(subject, raw)
}
