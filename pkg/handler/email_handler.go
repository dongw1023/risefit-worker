package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/risefit/email-worker/pkg/email"
)

type EmailHandler struct {
	Provider     email.Provider
	TemplatePath string
}

func NewEmailHandler(provider email.Provider) *EmailHandler {
	return &EmailHandler{
		Provider:     provider,
		TemplatePath: "templates",
	}
}

func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req email.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Logging task as requested
	log.Printf("Received email task: type=%s, email=%s", req.Template, req.Email)

	// Validate email address
	if _, err := mail.ParseAddress(req.Email); err != nil {
		log.Printf("Invalid email address: %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	// Validate template type and required fields
	subject, templateName, requiredFields, err := getTemplateMetadata(req.Template)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, field := range requiredFields {
		if _, ok := req.Data[field]; !ok {
			log.Printf("Missing required field %s for template %s", field, req.Template)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Missing required field: %s", field)})
			return
		}
	}

	// Render template
	tmplPath := filepath.Join(h.TemplatePath, templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmplPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error rendering email"})
		return
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, req.Data); err != nil {
		log.Printf("Error executing template %s: %v", tmplPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error rendering email"})
		return
	}

	// Send email
	err = h.Provider.Send(c.Request.Context(), req.Email, subject, rendered.String())
	if err != nil {
		log.Printf("Error sending email via provider: %v", err)
		// Return 5xx to signal Cloud Tasks to retry
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email"})
		return
	}

	c.Status(http.StatusOK)
}

func getTemplateMetadata(templateType string) (subject string, templateName string, requiredFields []string, err error) {
	switch templateType {
	case "verify_email":
		return "Verify your email", "verify_email.html", []string{"name", "verification_url"}, nil
	case "reset_password":
		return "Reset your password", "reset_password.html", []string{"name", "reset_url"}, nil
	default:
		return "", "", nil, fmt.Errorf("unsupported template: %s", templateType)
	}
}
