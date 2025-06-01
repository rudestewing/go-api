package service

import (
	"fmt"
	"go-api/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendWelcomeEmail sends a welcome email to newly registered users
func (s *EmailService) SendWelcomeEmail(userEmail, userName string) error {
	// Create message
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", m.FormatAddress(s.config.FromEmail, s.config.FromName))
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Welcome to Go API App!")

	// Create HTML body
	htmlBody := s.createWelcomeEmailHTML(userName)
	m.SetBody("text/html", htmlBody)

	// Create plain text alternative
	textBody := s.createWelcomeEmailText(userName)
	m.AddAlternative("text/plain", textBody)

	// Create dialer
	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.EmailUsername, s.config.EmailPassword)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

// createWelcomeEmailHTML creates the HTML version of the welcome email
func (s *EmailService) createWelcomeEmailHTML(userName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Go API App</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            background-color: #007bff;
            color: white;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            margin: -30px -30px 30px -30px;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .content {
            text-align: left;
        }
        .highlight {
            background-color: #e7f3ff;
            padding: 15px;
            border-left: 4px solid #007bff;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
            font-size: 14px;
        }
        .btn {
            display: inline-block;
            background-color: #007bff;
            color: white;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 5px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to Go API App!</h1>
        </div>

        <div class="content">
            <h2>Hello %s!</h2>

            <p>Thank you for registering with <strong>Go API App</strong>! We're excited to have you on board.</p>

            <div class="highlight">
                <h3>🚀 Your account has been successfully created!</h3>
                <p>You can now start using our API services with your registered email address.</p>
            </div>

            <h3>What's Next?</h3>
            <ul>
                <li>🔐 Log in to your account using your credentials</li>
                <li>📚 Explore our API documentation</li>
                <li>💻 Start building amazing applications</li>
                <li>🆘 Contact our support team if you need assistance</li>
            </ul>

            <p>If you have any questions or need help getting started, don't hesitate to reach out to our support team.</p>

            <p>Happy coding!</p>

            <p><strong>The Go API Team</strong></p>
        </div>

        <div class="footer">
            <p>This email was sent automatically. Please do not reply to this email.</p>
            <p>&copy; 2025 Go API App. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, userName)
}

// createWelcomeEmailText creates the plain text version of the welcome email
func (s *EmailService) createWelcomeEmailText(userName string) string {
	return fmt.Sprintf(`
Welcome to Go API App!

Hello %s!

Thank you for registering with Go API App! We're excited to have you on board.

Your account has been successfully created!
You can now start using our API services with your registered email address.

What's Next?
- Log in to your account using your credentials
- Explore our API documentation
- Start building amazing applications
- Contact our support team if you need assistance

If you have any questions or need help getting started, don't hesitate to reach out to our support team.

Happy coding!

The Go API Team

---
This email was sent automatically. Please do not reply to this email.
© 2025 Go API App. All rights reserved.
`, userName)
}

// SendEmail sends a generic email (can be used for other email types in the future)
func (s *EmailService) SendEmail(to, subject, htmlBody, textBody string) error {
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", m.FormatAddress(s.config.FromEmail, s.config.FromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set body
	if htmlBody != "" {
		m.SetBody("text/html", htmlBody)
		if textBody != "" {
			m.AddAlternative("text/plain", textBody)
		}
	} else {
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
