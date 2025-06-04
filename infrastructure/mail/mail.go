package mail

import (
	"bytes"
	"fmt"
	"go-api/config"
	"html/template"
	"path/filepath"
	"time"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	config        *config.Config
	templateCache map[string]*template.Template
}

// MailData represents the data structure for mail templates
type MailData map[string]any

func NewMailClient(cfg *config.Config) *MailService {
	service := &MailService{
		config:        cfg,
		templateCache: make(map[string]*template.Template),
	}

	// Load templates on initialization
	service.loadTemplates()

	return service
}

// loadTemplates loads all mail templates into memory
func (s *MailService) loadTemplates() {
	templateDir := "infrastructure/mail/templates"
	// Define available templates
	templates := []string{"welcome", "password_reset", "mail_verification"}

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

// SendMail sends a generic mail (base method)
func (s *MailService) SendMail(to, subject, htmlBody, textBody string) error {
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", m.FormatAddress(s.config.FromMail, s.config.FromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set body - prioritize HTML body
	if htmlBody != "" {
		m.SetBody("text/html", htmlBody)
	} else if textBody != "" {
		m.SetBody("text/plain", textBody)
	}

	// Create dialer
	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.MailUsername, s.config.MailPassword)

	// Send mail
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}

// SendTemplateMail sends a mail using a specified template with dynamic data
func (s *MailService) SendTemplateMail(to, subject, templateName string, data MailData) error {
	// Set default values
	if data["Year"] == 0 {
		data["Year"] = time.Now().Year()
	}
	if data["AppName"] == "" {
		data["AppName"] = "Go API App"
	}
	if data["SupportMail"] == "" {
		data["SupportMail"] = s.config.FromMail
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

	// Send mail with only HTML body
	return s.SendMail(to, subject, htmlBody, "")
}

// SendWelcomeMail sends a welcome mail using the welcome template
func (s *MailService) SendWelcomeMail(userMail, userName string) error {
	data := MailData{
		"UserName": userName,
		"LoginURL": "", // Add your login URL here if needed
	}

	return s.SendTemplateMail(userMail, "Welcome to Go API App!", "welcome", data)
}

// SendPasswordResetMail sends a password reset mail using the password_reset template
func (s *MailService) SendPasswordResetMail(userMail, userName, resetToken string, expirationMinutes int) error {
	data := MailData{
		"UserName":       userName,
		"ResetToken":     resetToken,
		"ResetURL":       "", // Add your reset URL here if needed
		"ExpirationTime": expirationMinutes,
	}

	return s.SendTemplateMail(userMail, "Password Reset Request", "password_reset", data)
}

// SendMailVerificationMail sends a mail verification using the mail_verification template
func (s *MailService) SendMailVerificationMail(userMail, userName, verificationCode string, expirationMinutes int) error {
	data := MailData{
		"UserName":         userName,
		"VerificationCode": verificationCode,
		"VerificationURL":  "", // Add your verification URL here if needed
		"ExpirationTime":   expirationMinutes,
	}

	return s.SendTemplateMail(userMail, "Please Verify Your Mail", "mail_verification", data)
}
