package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/risefit/email-worker/pkg/email"
	"github.com/risefit/email-worker/pkg/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type MockProvider struct {
	Sent bool
}

func (m *MockProvider) Send(ctx context.Context, to, subject, htmlContent string) error {
	m.Sent = true
	return nil
}

func TestSendEmail(t *testing.T) {
	mockProvider := &MockProvider{}
	h := &EmailHandler{
		Provider:     mockProvider,
		TemplatePath: "../../templates",
	}

	r := gin.New()
	r.POST("/send-email", middleware.InternalAuth("testkey"), h.SendEmail)

	payload := email.Request{
		Email:    "test@example.com",
		Template: "verify_email",
		Data: map[string]interface{}{
			"name":             "John Doe",
			"verification_url": "https://example.com/verify",
		},
	}
	body, _ := json.Marshal(payload)

	// Test Success
	req, _ := http.NewRequest(http.MethodPost, "/send-email", bytes.NewBuffer(body))
	req.Header.Set("X-Internal-API-Key", "testkey")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !mockProvider.Sent {
		t.Error("provider should have sent email")
	}

	// Test Unauthorized
	req, _ = http.NewRequest(http.MethodPost, "/send-email", bytes.NewBuffer(body))
	req.Header.Set("X-Internal-API-Key", "wrongkey")
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code for unauthorized: got %v want %v", status, http.StatusUnauthorized)
	}

	// Test Invalid Email
	payload.Email = "invalid-email"
	body, _ = json.Marshal(payload)
	req, _ = http.NewRequest(http.MethodPost, "/send-email", bytes.NewBuffer(body))
	req.Header.Set("X-Internal-API-Key", "testkey")
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid email: got %v want %v", status, http.StatusBadRequest)
	}
}
