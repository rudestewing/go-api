# Email Service Documentation

## Overview

The Email Service has been redesigned to be dynamic and template-based. You can now easily send emails by specifying a template name and providing the required variables.

## Folder Structure

```
emails/
└── templates/
    ├── welcome.html
    ├── welcome.txt
    ├── password_reset.html
    ├── password_reset.txt
    ├── email_verification.html
    └── email_verification.txt
```

## How to Use

### 1. Basic Usage with Template

```go
// Create email data
data := service.EmailData{
    UserName: "John Doe",
    LoginURL: "https://yourapp.com/login",
}

// Send email using template
err := emailService.SendTemplateEmail(
    "user@example.com",
    "Welcome to Our App!",
    "welcome",
    data
)
```

### 2. Pre-built Methods

The service provides convenient methods for common email types:

#### Welcome Email

```go
err := emailService.SendWelcomeEmail("user@example.com", "John Doe")
```

#### Password Reset Email

```go
err := emailService.SendPasswordResetEmail(
    "user@example.com",
    "John Doe",
    "reset_token_123",
    30 // expiration in minutes
)
```

#### Email Verification

```go
err := emailService.SendEmailVerificationEmail(
    "user@example.com",
    "John Doe",
    "verification_code_123",
    15 // expiration in minutes
)
```

## Available Template Variables

### EmailData Structure

```go
type EmailData struct {
    UserName         string // User's name
    Year             int    // Current year (auto-filled if not provided)
    LoginURL         string // Login page URL
    ResetToken       string // Password reset token
    ResetURL         string // Password reset URL
    ExpirationTime   int    // Token expiration time in minutes
    VerificationCode string // Email verification code
    VerificationURL  string // Email verification URL
    AppName          string // Application name (defaults to "Go API App")
    SupportEmail     string // Support email (defaults to config.FromEmail)
}
```

## Creating New Templates

### 1. Create HTML Template

Create a new file in `emails/templates/` with `.html` extension:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>{{.Subject}}</title>
    <!-- Add your CSS styles here -->
  </head>
  <body>
    <h1>Hello {{.UserName}}!</h1>
    <p>Your custom content here...</p>

    {{if .CustomVariable}}
    <p>Custom variable: {{.CustomVariable}}</p>
    {{end}}

    <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
  </body>
</html>
```

### 2. Create Text Template

Create a corresponding `.txt` file:

```txt
Hello {{.UserName}}!

Your custom content here...

{{if .CustomVariable}}
Custom variable: {{.CustomVariable}}
{{end}}

© {{.Year}} {{.AppName}}. All rights reserved.
```

### 3. Add Template to Service

Update the `loadTemplates()` method in `email_service.go`:

```go
templates := []string{"welcome", "password_reset", "email_verification", "your_new_template"}
```

### 4. Create a Convenience Method (Optional)

```go
func (s *EmailService) SendYourCustomEmail(userEmail, userName, customData string) error {
    data := EmailData{
        UserName:       userName,
        CustomVariable: customData,
    }

    return s.SendTemplateEmail(userEmail, "Your Subject", "your_new_template", data)
}
```

## Template Syntax

The templates use Go's `html/template` and `text/template` packages. Available syntax:

- `{{.VariableName}}` - Display variable
- `{{if .Variable}}...{{end}}` - Conditional blocks
- `{{if .Variable}}...{{else}}...{{end}}` - If-else blocks
- `{{range .Items}}...{{end}}` - Loop through slices
- `{{with .Variable}}...{{end}}` - Set context

## Example: Adding a Custom Template

Let's say you want to add an "account_suspended" template:

1. **Create `emails/templates/account_suspended.html`:**

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Account Suspended</title>
    <style>
      /* Your styles here */
    </style>
  </head>
  <body>
    <h1>Account Suspended</h1>
    <p>Hello {{.UserName}},</p>
    <p>Your account has been suspended for: {{.Reason}}</p>
    {{if .ContactURL}}
    <p><a href="{{.ContactURL}}">Contact Support</a></p>
    {{end}}
  </body>
</html>
```

2. **Create `emails/templates/account_suspended.txt`:**

```txt
Account Suspended

Hello {{.UserName}},

Your account has been suspended for: {{.Reason}}

{{if .ContactURL}}
Contact Support: {{.ContactURL}}
{{end}}
```

3. **Add to service and create method:**

```go
func (s *EmailService) SendAccountSuspendedEmail(userEmail, userName, reason string) error {
    data := EmailData{
        UserName: userName,
        Reason:   reason, // You might need to add this field to EmailData
    }

    return s.SendTemplateEmail(userEmail, "Account Suspended", "account_suspended", data)
}
```

4. **Usage:**

```go
err := emailService.SendAccountSuspendedEmail(
    "user@example.com",
    "John Doe",
    "Violation of terms of service"
)
```

## Benefits

1. **Seamless Usage** - Just specify template name and variables
2. **Maintainable** - Templates separated from code
3. **Flexible** - Easy to add new templates
4. **Consistent** - Automatic handling of common variables like year, app name
5. **Dual Format** - Automatic HTML and text versions
6. **Template Caching** - Templates loaded once at startup for performance
