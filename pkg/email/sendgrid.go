package email

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridProvider struct {
	client    *sendgrid.Client
	fromEmail string
}

func NewSendGridProvider(apiKey, fromEmail string) *SendGridProvider {
	return &SendGridProvider{
		client:    sendgrid.NewSendClient(apiKey),
		fromEmail: fromEmail,
	}
}

func (p *SendGridProvider) Send(ctx context.Context, to, subject, htmlContent string) error {
	from := mail.NewEmail("Risefit", p.fromEmail)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, "", htmlContent)

	response, err := p.client.SendWithContext(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid returned non-success status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
