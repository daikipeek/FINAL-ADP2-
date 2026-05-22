package platform

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func OpenDB() *sql.DB {
	db, err := sql.Open("postgres", Env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	for ctx.Err() == nil {
		if err := db.PingContext(ctx); err == nil {
			return db
		}
		time.Sleep(time.Second)
	}
	log.Fatal(ctx.Err())
	return nil
}

func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: Env("REDIS_ADDR", "localhost:6379")})
}

func NATS() *nats.Conn {
	nc, err := nats.Connect(Env("NATS_URL", nats.DefaultURL), nats.RetryOnFailedConnect(true), nats.MaxReconnects(20), nats.ReconnectWait(time.Second))
	if err != nil {
		log.Fatal(err)
	}
	return nc
}

func Dial(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

type jsonCodec struct{}

func (jsonCodec) Name() string                       { return "json" }
func (jsonCodec) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (jsonCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

func ServeGRPC(addr string, register func(*grpc.Server)) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	register(s)
	log.Printf("gRPC listening on %s", addr)
	log.Fatal(s.Serve(lis))
}
