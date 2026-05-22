package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"car-rentals/internal/pb"
	"car-rentals/internal/platform"
	"github.com/nats-io/nats.go"
)

type Service struct{}

func New(nc *nats.Conn) *Service {
	s := &Service{}
	for _, subject := range []string{"rental.created", "rental.cancelled", "payment.success", "payment.failed"} {
		subject := subject
		_, _ = nc.Subscribe(subject, func(m *nats.Msg) {
			var event map[string]any
			_ = json.Unmarshal(m.Data, &event)
			to := platform.Env("DEFAULT_EMAIL_TO", "customer@example.com")
			body := fmt.Sprintf("Event %s\n\n%s", subject, string(m.Data))
			_, err := s.SendNotificationEmail(context.Background(), &pb.NotificationRequest{To: to, Subject: "Car Rentals: " + subject, Body: body})
			if err != nil {
				log.Printf("email skipped: %v", err)
			}
		})
	}
	return s
}

func (s *Service) SendNotificationEmail(_ context.Context, r *pb.NotificationRequest) (*pb.Empty, error) {
	host := platform.Env("SMTP_HOST", "smtp.gmail.com")
	port := platform.Env("SMTP_PORT", "587")
	user := platform.Env("SMTP_USER", "")
	pass := platform.Env("SMTP_PASS", "")
	from := platform.Env("SMTP_FROM", user)
	if user == "" || pass == "" {
		log.Printf("email dry-run to=%s subject=%q body=%q", r.To, r.Subject, r.Body)
		return &pb.Empty{}, nil
	}
	msg := strings.Join([]string{
		"From: " + from,
		"To: " + r.To,
		"Subject: " + r.Subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		r.Body,
	}, "\r\n")
	auth := smtp.PlainAuth("", user, pass, host)
	return &pb.Empty{}, smtp.SendMail(host+":"+port, auth, from, []string{r.To}, []byte(msg))
}
