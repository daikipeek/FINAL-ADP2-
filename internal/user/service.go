package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"car-rentals/internal/pb"
	"car-rentals/internal/platform"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{ db *sql.DB }

func New(db *sql.DB) *Service { return &Service{db: db} }

func (s *Service) RegisterUser(ctx context.Context, r *pb.RegisterUserRequest) (*pb.UserResponse, error) {
	role := r.Role
	if role == "" {
		role = "customer"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	u := &pb.UserResponse{Id: uuid.NewString(), Name: r.Name, Email: r.Email, Role: role}
	if _, err = tx.ExecContext(ctx, `insert into users(id,name,email,password_hash,role) values($1,$2,$3,$4,$5)`, u.Id, u.Name, u.Email, string(hash), u.Role); err != nil {
		return nil, err
	}
	return u, tx.Commit()
}

func (s *Service) LoginUser(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginResponse, error) {
	var u pb.UserResponse
	var hash string
	err := s.db.QueryRowContext(ctx, `select id,name,email,password_hash,role from users where email=$1`, r.Email).Scan(&u.Id, &u.Name, &u.Email, &hash, &u.Role)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.Password)) != nil {
		return nil, errors.New("invalid email or password")
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.Id, "role": u.Role, "exp": time.Now().Add(24 * time.Hour).Unix(),
	}).SignedString([]byte(platform.Env("JWT_SECRET", "dev-secret")))
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: token, User: &u}, nil
}

func (s *Service) GetUserProfile(ctx context.Context, r *pb.GetByIdRequest) (*pb.UserResponse, error) {
	var u pb.UserResponse
	err := s.db.QueryRowContext(ctx, `select id,name,email,role from users where id=$1`, r.Id).Scan(&u.Id, &u.Name, &u.Email, &u.Role)
	return &u, err
}

func (s *Service) UpdateUserProfile(ctx context.Context, r *pb.UpdateUserProfileRequest) (*pb.UserResponse, error) {
	var u pb.UserResponse
	err := s.db.QueryRowContext(ctx, `update users set name=$2,email=$3,updated_at=now() where id=$1 returning id,name,email,role`, r.Id, r.Name, r.Email).Scan(&u.Id, &u.Name, &u.Email, &u.Role)
	return &u, err
}
