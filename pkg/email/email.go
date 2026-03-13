package email

import (
	"context"
)

// Request defines the incoming send-email payload.
type Request struct {
	Email    string                 `json:"email"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

// Provider defines the interface for an email service provider.
type Provider interface {
	Send(ctx context.Context, to, subject, htmlContent string) error
}
