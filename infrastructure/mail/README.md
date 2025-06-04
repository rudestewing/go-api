# Mail Service Documentation

## Overview

The Mail Service has been redesigned to be dynamic and template-based. You can now easily send mails by specifying a template name and providing the required variables.

## Folder Structure

```
infrastructure/mail/
├── mail.go
├── README.md
└── templates/
    ├── welcome.html
    ├── password_reset.html
    └── mail_verification.html
```

## How to Use

### 1. Basic Usage with Template

```go
// Create mail data
data := mail.MailData{
    "UserName": "John Doe",
    "LoginURL": "https://yourapp.com/login",
}

// Send mail using template
err := mailService.SendTemplateMail(
    "user@example.com",
    "Welcome to Our App!",
    "welcome",
    data
)
```

### 2. Pre-built Methods

The service provides convenient methods for common mail types:

#### Welcome Mail

```go
err := mailService.SendWelcomeMail("user@example.com", "John Doe")
```

#### Password Reset Mail

```go
err := mailService.SendPasswordResetMail(
    "user@example.com",
    "John Doe",
    "reset-token-123",
    30 // expiration in minutes
)
```

#### Mail Verification

```go
err := mailService.SendMailVerificationMail(
    "user@example.com",
    "John Doe",
    "verification-code-123",
    60 // expiration in minutes
)
```

## Configuration

Add mail configuration to your `config.yaml`:

```yaml
mail:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your-email@example.com"
  password: "your-app-password"
  from_name: "Go API App"
  from_mail: "noreply@example.com"
```

## Template Variables

### Available in All Templates

- `Year` - Current year (automatically set)
- `AppName` - Application name (default: "Go API App")
- `SupportMail` - Support mail address (from config)

### Welcome Template Variables

- `UserName` - The user's name
- `LoginURL` - URL to login page (optional)

### Password Reset Template Variables

- `UserName` - The user's name
- `ResetToken` - Password reset token
- `ResetURL` - Password reset URL (optional)
- `ExpirationTime` - Token expiration time in minutes

### Mail Verification Template Variables

- `UserName` - The user's name
- `VerificationCode` - Verification code
- `VerificationURL` - Verification URL (optional)
- `ExpirationTime` - Code expiration time in minutes

## Creating Custom Templates

1. Create a new HTML template in `infrastructure/mail/templates/`
2. Use Go template syntax with the variables you need
3. Add the template name to the `templates` slice in `loadTemplates()` method
4. Create a new method in `MailService` to send your custom mail

### Example Custom Template

Create `infrastructure/mail/templates/account_suspended.html`:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Account Suspended</title>
  </head>
  <body>
    <h1>Account Suspended</h1>
    <p>Dear {{.UserName}},</p>
    <p>Your account has been suspended due to: {{.Reason}}</p>
    <p>
      If you believe this is an error, please contact support at
      {{.SupportMail}}
    </p>
    <hr />
    <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
  </body>
</html>
```

Add method to `MailService`:

```go
func (s *MailService) SendAccountSuspendedMail(userMail, userName, reason string) error {
    data := MailData{
        "UserName": userName,
        "Reason":   reason,
    }

    return s.SendTemplateMail(userMail, "Account Suspended", "account_suspended", data)
}
```

Use it:

```go
err := mailService.SendAccountSuspendedMail(
    "user@example.com",
    "John Doe",
    "Multiple failed login attempts"
)
```

## Error Handling

The mail service includes comprehensive error handling:

- Template not found errors
- SMTP connection errors
- Template rendering errors

Always check for errors when sending mails:

```go
if err := mailService.SendWelcomeMail("user@example.com", "John Doe"); err != nil {
    log.Printf("Failed to send welcome mail: %v", err)
    // Handle error appropriately
}
```

## Background Mail Sending

For better user experience, send mails in background goroutines:

```go
go func() {
    if err := mailService.SendWelcomeMail(req.Email, req.Name); err != nil {
        log.Printf("Failed to send welcome mail: %v", err)
    }
}()
```

## Testing

When testing, you can mock the mail service or use a test SMTP server like MailHog for development.

## Security Notes

- Store SMTP credentials securely
- Use app passwords for Gmail/Google Workspace
- Consider using environment-specific mail settings
- Validate mail addresses before sending
- Implement rate limiting for mail sending
