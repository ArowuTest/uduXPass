package email

import (
	"context"
	"fmt"
	"time"

	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/infrastructure/pdf"
)

// SendTicketPDFEmail sends ticket PDFs to the customer
// Note: tickets must have OrderLine relation preloaded with TicketTier
func (s *SMTPEmailService) SendTicketPDFEmail(ctx context.Context, order *entities.Order, tickets []*entities.Ticket, orderLines []*entities.OrderLine, event *entities.Event) error {
	subject := fmt.Sprintf("Your Tickets for %s - Order %s", event.Name, order.Code)

	// Generate PDF for each ticket
	pdfGenerator := pdf.NewTicketPDFGenerator()
	attachments := make([]EmailAttachment, 0, len(tickets))

	for i, ticket := range tickets {
		// Find the corresponding order line for this ticket
		var orderLine *entities.OrderLine
		for _, ol := range orderLines {
			if ol.ID == ticket.OrderLineID {
				orderLine = ol
				break
			}
		}
		
		if orderLine == nil || orderLine.TicketTier == nil {
			return fmt.Errorf("order line or ticket tier not found for ticket %s", ticket.ID.String())
		}
		
		// Prepare ticket data
		ticketData := pdf.TicketData{
			TicketID:      ticket.ID.String(),
			QRCode:        ticket.SerialNumber,
			EventName:     event.Name,
			EventDate:     event.EventDate,
			VenueName:     event.VenueName,
			VenueAddress:  event.VenueAddress,
			TierName:      orderLine.TicketTier.Name,
			Price:         orderLine.TicketTier.Price,
			CustomerName:  order.CustomerFirstName + " " + order.CustomerLastName,
			CustomerEmail: order.CustomerEmail,
			OrderID:       order.Code,
			TicketNumber:  i + 1,
			TotalTickets:  len(tickets),
		}

		// Generate PDF
		pdfBytes, err := pdfGenerator.GenerateTicketPDF(ticketData)
		if err != nil {
			return fmt.Errorf("failed to generate PDF for ticket %s: %w", ticket.ID.String(), err)
		}

		// Add to attachments
		filename := fmt.Sprintf("ticket_%d_%s.pdf", i+1, ticket.SerialNumber[:8])
		attachments = append(attachments, EmailAttachment{
			Filename:    filename,
			ContentType: "application/pdf",
			Data:        pdfBytes,
		})
	}

	// Prepare email body
	customerName := order.CustomerFirstName + " " + order.CustomerLastName
	data := map[string]interface{}{
		"OrderCode":    order.Code,
		"CustomerName": customerName,
		"EventName":    event.Name,
		"EventDate":    event.EventDate.Format("Monday, January 2, 2006 at 3:04 PM"),
		"VenueName":    event.VenueName,
		"VenueAddress": event.VenueAddress,
		"TicketCount":  len(tickets),
		"Total":        order.TotalAmount,
		"Year":         time.Now().Year(),
	}

	body, err := s.renderTemplate("ticket_pdf_email.html", data)
	if err != nil {
		// Use fallback template if custom template doesn't exist
		body = s.getTicketPDFEmailTemplate(data)
	}

	// Send email with PDF attachments
	return s.sendEmailWithAttachments(order.CustomerEmail, subject, body, attachments)
}

// getTicketPDFEmailTemplate returns a default HTML template for ticket PDF emails
func (s *SMTPEmailService) getTicketPDFEmailTemplate(data map[string]interface{}) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #0066cc 0%%, #0052a3 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .ticket-info { background: white; padding: 20px; margin: 20px 0; border-radius: 8px; border-left: 4px solid #0066cc; }
        .info-row { margin: 10px 0; }
        .label { font-weight: bold; color: #666; }
        .value { color: #333; }
        .footer { text-align: center; margin-top: 30px; padding: 20px; color: #999; font-size: 12px; }
        .button { display: inline-block; padding: 12px 30px; background: #0066cc; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎫 Your Tickets Are Ready!</h1>
        <p>Order Confirmation: %s</p>
    </div>
    
    <div class="content">
        <p>Hi %s,</p>
        
        <p>Thank you for your purchase! Your tickets for <strong>%s</strong> are attached to this email as PDF files.</p>
        
        <div class="ticket-info">
            <h3>Event Details</h3>
            <div class="info-row">
                <span class="label">Event:</span>
                <span class="value">%s</span>
            </div>
            <div class="info-row">
                <span class="label">Date & Time:</span>
                <span class="value">%s</span>
            </div>
            <div class="info-row">
                <span class="label">Venue:</span>
                <span class="value">%s</span>
            </div>
            <div class="info-row">
                <span class="label">Address:</span>
                <span class="value">%s</span>
            </div>
            <div class="info-row">
                <span class="label">Number of Tickets:</span>
                <span class="value">%d</span>
            </div>
            <div class="info-row">
                <span class="label">Total Paid:</span>
                <span class="value">₦%.2f</span>
            </div>
        </div>
        
        <h3>📎 Attached Files</h3>
        <p>You will find %d PDF ticket(s) attached to this email. Each ticket contains:</p>
        <ul>
            <li>✓ QR code for entry validation</li>
            <li>✓ Event details and venue information</li>
            <li>✓ Ticket holder information</li>
            <li>✓ Important entry instructions</li>
        </ul>
        
        <h3>📱 How to Use Your Tickets</h3>
        <ol>
            <li>Download and save the PDF tickets to your device</li>
            <li>You can print them or show them on your mobile device</li>
            <li>Present the QR code at the venue entrance for scanning</li>
            <li>Each ticket is valid for one entry only</li>
        </ol>
        
        <p><strong>Important:</strong> Please arrive early to avoid queues. Doors open 1 hour before the event starts.</p>
        
        <p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>
        
        <p>Enjoy the event! 🎉</p>
    </div>
    
    <div class="footer">
        <p>Powered by uduXPass - Enterprise Ticketing Platform</p>
        <p>&copy; %d uduXPass. All rights reserved.</p>
    </div>
</body>
</html>
`,
		data["OrderCode"],
		data["CustomerName"],
		data["EventName"],
		data["EventName"],
		data["EventDate"],
		data["VenueName"],
		data["VenueAddress"],
		data["TicketCount"],
		data["Total"],
		data["TicketCount"],
		data["Year"],
	)
}
