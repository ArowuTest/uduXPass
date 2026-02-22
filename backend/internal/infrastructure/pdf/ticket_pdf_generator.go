package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

// TicketPDFGenerator generates PDF tickets with QR codes
type TicketPDFGenerator struct{}

// NewTicketPDFGenerator creates a new PDF generator
func NewTicketPDFGenerator() *TicketPDFGenerator {
	return &TicketPDFGenerator{}
}

// TicketData contains all information needed to generate a ticket PDF
type TicketData struct {
	TicketID      string
	QRCode        string
	EventName     string
	EventDate     time.Time
	VenueName     string
	VenueAddress  string
	TierName      string
	Price         float64
	CustomerName  string
	CustomerEmail string
	OrderID       string
	TicketNumber  int
	TotalTickets  int
}

// GenerateTicketPDF generates a professional PDF ticket
func (g *TicketPDFGenerator) GenerateTicketPDF(ticket TicketData) ([]byte, error) {
	// Create new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins
	pdf.SetMargins(20, 20, 20)

	// Header - Event Name
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(0, 102, 204) // Blue color
	pdf.CellFormat(0, 15, ticket.EventName, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Ticket Type Badge
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(0, 102, 204)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 10, fmt.Sprintf("  %s  ", ticket.TierName), "", 1, "C", true, 0, "")
	pdf.Ln(10)

	// QR Code Section
	qrImage, err := g.generateQRCodeImage(ticket.QRCode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Center QR code
	qrSize := 60.0
	pageWidth := 210.0 // A4 width
	qrX := (pageWidth - qrSize) / 2
	
	// Register QR code image
	imageOpts := gofpdf.ImageOptions{
		ImageType: "PNG",
	}
	pdf.RegisterImageOptionsReader("qrcode", imageOpts, bytes.NewReader(qrImage))
	pdf.ImageOptions("qrcode", qrX, pdf.GetY(), qrSize, qrSize, false, imageOpts, 0, "")
	pdf.Ln(qrSize + 10)

	// Ticket Code
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 5, fmt.Sprintf("Ticket Code: %s", ticket.QRCode), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Event Details Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "Event Details", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Draw line
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	// Event information
	pdf.SetFont("Arial", "", 11)
	g.addInfoRow(pdf, "Date & Time:", ticket.EventDate.Format("Monday, January 2, 2006 at 3:04 PM"))
	g.addInfoRow(pdf, "Venue:", ticket.VenueName)
	g.addInfoRow(pdf, "Address:", ticket.VenueAddress)
	g.addInfoRow(pdf, "Ticket Type:", ticket.TierName)
	g.addInfoRow(pdf, "Price:", fmt.Sprintf("₦%.2f", ticket.Price))
	pdf.Ln(5)

	// Customer Details Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "Ticket Holder", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Draw line
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	// Customer information
	pdf.SetFont("Arial", "", 11)
	g.addInfoRow(pdf, "Name:", ticket.CustomerName)
	g.addInfoRow(pdf, "Email:", ticket.CustomerEmail)
	g.addInfoRow(pdf, "Order ID:", ticket.OrderID)
	g.addInfoRow(pdf, "Ticket:", fmt.Sprintf("%d of %d", ticket.TicketNumber, ticket.TotalTickets))
	pdf.Ln(10)

	// Important Information Box
	pdf.SetFillColor(255, 248, 220) // Light yellow
	pdf.SetDrawColor(255, 193, 7)   // Yellow border
	pdf.Rect(20, pdf.GetY(), 170, 25, "FD")
	
	currentY := pdf.GetY()
	pdf.SetY(currentY + 3)
	pdf.SetX(25)
	
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 5, "Important Information", "", 1, "L", false, 0, "")
	pdf.SetX(25)
	
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(60, 60, 60)
	pdf.MultiCell(160, 4, "• Present this ticket (printed or on mobile) at the venue entrance\n• QR code will be scanned for entry validation\n• Each ticket is valid for one entry only\n• Ticket is non-transferable and non-refundable", "", "L", false)
	
	pdf.SetY(currentY + 25)
	pdf.Ln(10)

	// Footer
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 5, "Powered by uduXPass - Enterprise Ticketing Platform", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Generated on %s", time.Now().Format("January 2, 2006 at 3:04 PM")), "", 1, "C", false, 0, "")

	// Output PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// addInfoRow adds a labeled information row
func (g *TicketPDFGenerator) addInfoRow(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(45, 6, label, "", 0, "L", false, 0, "")
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 6, value, "", 1, "L", false, 0, "")
}

// generateQRCodeImage generates a QR code image as PNG bytes
func (g *TicketPDFGenerator) generateQRCodeImage(content string) ([]byte, error) {
	// Generate QR code with medium recovery level
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// Set size to 256x256 pixels
	qr.DisableBorder = false

	// Generate PNG bytes
	pngBytes, err := qr.PNG(256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code PNG: %w", err)
	}

	return pngBytes, nil
}
