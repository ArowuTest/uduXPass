package services

import (
	"context"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	// SendTicketEmail sends ticket information to the customer
	SendTicketEmail(ctx context.Context, order *entities.Order, tickets []*entities.Ticket) error
	
	// SendTicketPDFEmail sends ticket PDFs to the customer
	SendTicketPDFEmail(ctx context.Context, order *entities.Order, tickets []*entities.Ticket, event *entities.Event) error
	
	// SendOrderConfirmation sends order confirmation email
	SendOrderConfirmation(ctx context.Context, order *entities.Order) error
	
	// SendWelcomeEmail sends welcome email to new users
	SendWelcomeEmail(ctx context.Context, user *entities.User) error
	
	// SendPasswordResetEmail sends password reset link
	SendPasswordResetEmail(ctx context.Context, email, resetToken string) error
}
