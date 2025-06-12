package email

import (
	"bytes"
	"fmt"
	"go-api/config"
	"html/template"
	"path/filepath"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config        *config.Config
	templateCache map[string]*template.Template
}

// EmailData represents the data structure for email templates
type EmailData map[string]any

func NewEmailClient(cfg *config.Config) *EmailService {
	service := &EmailService{
		config:        cfg,
		templateCache: make(map[string]*template.Template),
	}

	// Load templates on initialization
	service.loadTemplates()

	return service
}

// loadTemplates loads all email templates into memory
func (s *EmailService) loadTemplates() {
	templateDir := "infrastructure/emails/templates"

	// Define available templates
	templates := []string{"welcome", "password_reset", "email_verification"}

	for _, tmplName := range templates {
		htmlPath := filepath.Join(templateDir, tmplName+".html")

		// Parse HTML template
		htmlTmpl, err := template.ParseFiles(htmlPath)
		if err != nil {
			fmt.Printf("Warning: Could not load HTML template %s: %v\n", htmlPath, err)
			continue
		}

		// Store HTML template
		s.templateCache[tmplName] = htmlTmpl
	}
}

// SendEmail sends a generic email (base method)
func (s *EmailService) SendEmail(to, subject, htmlBody, textBody string) error {
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", m.FormatAddress(s.config.FromEmail, s.config.FromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set body - prioritize HTML body
	if htmlBody != "" {
		m.SetBody("text/html", htmlBody)
	} else if textBody != "" {
		m.SetBody("text/plain", textBody)
	}

	// Create dialer
	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.EmailUsername, s.config.EmailPassword)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendTemplateEmail sends an email using a specified template with dynamic data
func (s *EmailService) SendTemplateEmail(to, subject, templateName string, data EmailData) error {
	// Set default values
	if data["Year"] == 0 {
		data["Year"] = time.Now().Year()
	}
	if data["AppName"] == "" {
		data["AppName"] = "Go API App"
	}
	if data["SupportEmail"] == "" {
		data["SupportEmail"] = s.config.FromEmail
	}

	// Get HTML template
	htmlTemplate := s.templateCache[templateName]

	if htmlTemplate == nil {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	var htmlBody string
	var err error

	// Render HTML template
	var htmlBuf bytes.Buffer
	if err = htmlTemplate.Execute(&htmlBuf, data); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}
	htmlBody = htmlBuf.String()

	// Send email with only HTML body
	return s.SendEmail(to, subject, htmlBody, "")
}

// SendWelcomeEmail sends a welcome email using the welcome template
func (s *EmailService) SendWelcomeEmail(userEmail, userName string) error {
	data := EmailData{
		"UserName": userName,
		"LoginURL": "", // Add your login URL here if needed
	}

	return s.SendTemplateEmail(userEmail, "Welcome to Go API App!", "welcome", data)
}

// SendPasswordResetEmail sends a password reset email using the password_reset template
func (s *EmailService) SendPasswordResetEmail(userEmail, userName, resetToken string, expirationMinutes int) error {
	data := EmailData{
		"UserName":       userName,
		"ResetToken":     resetToken,
		"ResetURL":       "", // Add your reset URL here if needed
		"ExpirationTime": expirationMinutes,
	}

	return s.SendTemplateEmail(userEmail, "Password Reset Request", "password_reset", data)
}

// SendEmailVerificationEmail sends an email verification using the email_verification template
func (s *EmailService) SendEmailVerificationEmail(userEmail, userName, verificationCode string, expirationMinutes int) error {
	data := EmailData{
		"UserName":         userName,
		"VerificationCode": verificationCode,
		"VerificationURL":  "", // Add your verification URL here if needed
		"ExpirationTime":   expirationMinutes,
	}

	return s.SendTemplateEmail(userEmail, "Please Verify Your Email", "email_verification", data)
}
