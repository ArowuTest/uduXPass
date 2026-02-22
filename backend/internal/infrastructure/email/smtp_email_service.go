package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/services"
)

type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService() services.EmailService {
	return &SMTPEmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

// SendTicketEmail sends ticket information to the customer
func (s *SMTPEmailService) SendTicketEmail(ctx context.Context, order *entities.Order, tickets []*entities.Ticket) error {
	subject := fmt.Sprintf("Your Tickets for Order %s", order.Code)
	
	// Prepare ticket data for template
	type TicketData struct {
		Code       string
		TierName   string
		QRCodeURL  string
		EventName  string
		EventDate  time.Time
		VenueName  string
		VenueCity  string
	}
	
	ticketDataList := make([]TicketData, len(tickets))
	for i, ticket := range tickets {
		qrURL := ""
		if ticket.QRCodeImageURL != nil {
			qrURL = *ticket.QRCodeImageURL
		}
		ticketDataList[i] = TicketData{
			Code:      ticket.SerialNumber,
			TierName:  "Ticket", // Will be populated from order line
			QRCodeURL: qrURL,
			// Event details would come from joined data
		}
	}
	
	customerName := order.CustomerFirstName + " " + order.CustomerLastName
	data := map[string]interface{}{
		"OrderCode":    order.Code,
		"CustomerName": customerName,
		"Tickets":      ticketDataList,
		"Total":        order.TotalAmount,
		"Year":         time.Now().Year(),
	}
	
	body, err := s.renderTemplate("ticket_email.html", data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}
	
	return s.sendEmail(order.CustomerEmail, subject, body)
}

// SendOrderConfirmation sends order confirmation email
func (s *SMTPEmailService) SendOrderConfirmation(ctx context.Context, order *entities.Order) error {
	subject := fmt.Sprintf("Order Confirmation - %s", order.Code)
	
	customerName := order.CustomerFirstName + " " + order.CustomerLastName
	data := map[string]interface{}{
		"OrderCode":    order.Code,
		"CustomerName": customerName,
		"Total":        order.TotalAmount,
		"Status":       order.Status,
		"CreatedAt":    order.CreatedAt,
		"Year":         time.Now().Year(),
	}
	
	body, err := s.renderTemplate("order_confirmation.html", data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}
	
	return s.sendEmail(order.CustomerEmail, subject, body)
}

// SendWelcomeEmail sends welcome email to new users
func (s *SMTPEmailService) SendWelcomeEmail(ctx context.Context, user *entities.User) error {
	subject := "Welcome to uduXPass!"
	
	data := map[string]interface{}{
		"FirstName": user.FirstName,
		"Email":     user.Email,
		"Year":      time.Now().Year(),
	}
	
	body, err := s.renderTemplate("welcome_email.html", data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}
	
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	return s.sendEmail(email, subject, body)
}

// SendPasswordResetEmail sends password reset link
func (s *SMTPEmailService) SendPasswordResetEmail(ctx context.Context, email, resetToken string) error {
	subject := "Password Reset Request"
	
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), resetToken)
	
	data := map[string]interface{}{
		"ResetURL": resetURL,
		"Year":     time.Now().Year(),
	}
	
	body, err := s.renderTemplate("password_reset.html", data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}
	
	return s.sendEmail(email, subject, body)
}

// sendEmail sends an email using SMTP
func (s *SMTPEmailService) sendEmail(to, subject, body string) error {
	// If SMTP is not configured, log and return (dev mode)
	if s.host == "" || s.port == "" {
		fmt.Printf("[Email Service] Would send email to %s: %s\n", to, subject)
		fmt.Printf("[Email Service] Body preview: %s\n", body[:min(len(body), 200)])
		return nil
	}
	
	// Build email message
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body))
	
	// Connect to SMTP server
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	fmt.Printf("[Email Service] Successfully sent email to %s: %s\n", to, subject)
	return nil
}

// renderTemplate renders an HTML email template
func (s *SMTPEmailService) renderTemplate(templateName string, data interface{}) (string, error) {
	templatePath := filepath.Join("templates", "emails", templateName)
	
	// Check if template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// Use default template if file doesn't exist
		return s.getDefaultTemplate(templateName, data)
	}
	
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// getDefaultTemplate returns a simple default template
func (s *SMTPEmailService) getDefaultTemplate(templateName string, data interface{}) (string, error) {
	switch templateName {
	case "ticket_email.html":
		return s.renderTicketTemplate(data)
	case "order_confirmation.html":
		return s.renderOrderConfirmationTemplate(data)
	case "welcome_email.html":
		return s.renderWelcomeTemplate(data)
	case "password_reset.html":
		return s.renderPasswordResetTemplate(data)
	default:
		return "<html><body><p>Email content</p></body></html>", nil
	}
}

func (s *SMTPEmailService) renderTicketTemplate(data interface{}) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .ticket { border: 2px solid #667eea; border-radius: 8px; padding: 20px; margin: 20px 0; }
        .qr-code { text-align: center; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎟️ Your Tickets are Ready!</h1>
        </div>
        <div style="padding: 20px;">
            <p>Hi {{.CustomerName}},</p>
            <p>Your tickets for order <strong>{{.OrderCode}}</strong> are ready!</p>
            
            {{range .Tickets}}
            <div class="ticket">
                <h3>{{.TierName}}</h3>
                <p><strong>Ticket Code:</strong> {{.Code}}</p>
                <div class="qr-code">
                    <img src="{{.QRCodeURL}}" alt="QR Code" style="max-width: 200px;">
                </div>
                <p style="font-size: 12px; color: #666;">Present this QR code at the venue for entry</p>
            </div>
            {{end}}
            
            <p><strong>Total Paid:</strong> ₦{{.Total}}</p>
            <p>See you at the event! 🎉</p>
        </div>
        <div class="footer">
            <p>&copy; {{.Year}} uduXPass. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
	t, err := template.New("ticket").Parse(tmpl)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

func (s *SMTPEmailService) renderOrderConfirmationTemplate(data interface{}) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .content { padding: 20px; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Order Confirmed!</h1>
        </div>
        <div class="content">
            <p>Hi {{.CustomerName}},</p>
            <p>Your order <strong>{{.OrderCode}}</strong> has been confirmed!</p>
            <p><strong>Total:</strong> ₦{{.Total}}</p>
            <p><strong>Status:</strong> {{.Status}}</p>
            <p>You will receive your tickets shortly after payment is processed.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.Year}} uduXPass. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
	t, err := template.New("order").Parse(tmpl)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

func (s *SMTPEmailService) renderWelcomeTemplate(data interface{}) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .content { padding: 20px; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to uduXPass!</h1>
        </div>
        <div class="content">
            <p>Hi {{.FirstName}},</p>
            <p>Welcome to uduXPass - your premium event ticketing platform!</p>
            <p>We're excited to have you join our community. Start exploring amazing events and book your tickets today!</p>
            <p><strong>Your registered email:</strong> {{.Email}}</p>
        </div>
        <div class="footer">
            <p>&copy; {{.Year}} uduXPass. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

func (s *SMTPEmailService) renderPasswordResetTemplate(data interface{}) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .content { padding: 20px; }
        .button { background-color: #667eea; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Password Reset Request</h1>
        </div>
        <div class="content">
            <p>You requested to reset your password.</p>
            <p>Click the button below to reset your password:</p>
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
            <p style="font-size: 12px; color: #666;">If you didn't request this, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.Year}} uduXPass. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
	t, err := template.New("reset").Parse(tmpl)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
