# Car Rentals Microservices

Full-stack car rental system with Go microservices, gRPC service-to-service communication, an HTTP API Gateway, PostgreSQL, Redis caching, NATS events, Gmail SMTP-compatible email notifications, and a React frontend.

## Run

```bash
docker-compose up --build
```

Open:

- Frontend: http://localhost:3000
- API Gateway health: http://localhost:8080/health
- API Gateway metrics: http://localhost:8080/metrics
- NATS monitor: http://localhost:8222
- Grafana: http://localhost:3001 (`admin` / `admin`)
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100
- Tempo: http://localhost:3200

PostgreSQL is initialized from `migrations/` and includes sample cars. Redis caches available cars for the Car Service. Rental and Payment services publish NATS events consumed by Notification Service.

## Gmail SMTP

Set these environment variables in `docker-compose.yml` for real emails:

```text
SMTP_USER=your@gmail.com
SMTP_PASS=your-gmail-app-password
SMTP_FROM=your@gmail.com
DEFAULT_EMAIL_TO=customer@example.com
```

When SMTP credentials are blank, Notification Service logs a dry-run email.

## Tests

```bash
go test ./...
cd frontend && npm test
```

## HTTP API

- `POST /api/users/register`
- `POST /api/users/login`
- `GET /api/users/{id}`
- `PUT /api/users/{id}`
- `GET /api/cars`
- `GET /api/cars/available`
- `GET /api/cars/{id}`
- `POST /api/cars`
- `PUT /api/cars/{id}/status`
- `POST /api/cars/{id}/favorite`
- `GET /api/users/{id}/favorites`
- `POST /api/cars/{id}/reviews`
- `POST /api/rentals`
- `GET /api/rentals/{id}`
- `GET /api/users/{id}/rentals`
- `POST /api/rentals/{id}/cancel`
- `POST /api/rentals/{id}/complete`
- `POST /api/rentals/price`
- `POST /api/payments`
- `GET /api/payments/{id}`
- `POST /api/payments/{id}/refund`
- `POST /api/notifications/email`

## Notes

The `proto/car_rentals.proto` file is the service contract. The repo uses a small JSON gRPC codec and manually registered descriptors in `internal/pb` so the Docker build does not require `protoc`.

## Requirement Checklist

- Clean architecture: each service is separated into `cmd/<service>` entrypoints plus `internal/<domain>` business logic and `internal/platform` infrastructure.
- Microservices and gateway: API Gateway plus User, Car, Rental, Payment, and Notification services.
- gRPC endpoints: 24 RPC methods are defined in `proto/car_rentals.proto`.
- Message queue: NATS subjects include `rental.created`, `rental.cancelled`, `payment.success`, and `payment.failed`.
- Databases, cache, migrations, transactions: PostgreSQL migrations are in `migrations/`, Redis caches available cars, and write flows use SQL transactions.
- Email: Notification Service supports Gmail SMTP and dry-run logging.
- Tests: Go unit tests plus a bufconn gRPC integration test are included.
- Frontend bonus: React web app is in `frontend/`.
- Grafana bonus: Grafana, Prometheus metrics, Loki logs, and Tempo tracing backend are wired in `docker-compose.yml`.
