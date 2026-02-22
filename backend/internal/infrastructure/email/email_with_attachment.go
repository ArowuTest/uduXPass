package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
)

// EmailAttachment represents a file attachment
type EmailAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// sendEmailWithAttachments sends an email with PDF attachments
func (s *SMTPEmailService) sendEmailWithAttachments(to, subject, htmlBody string, attachments []EmailAttachment) error {
	// If SMTP is not configured, log and return (dev mode)
	if s.host == "" || s.port == "" {
		fmt.Printf("[Email Service] Would send email with %d attachments to %s: %s\n", len(attachments), to, subject)
		fmt.Printf("[Email Service] Body preview: %s\n", htmlBody[:min(len(htmlBody), 200)])
		return nil
	}

	// Create multipart message
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Set boundary for multipart message
	boundary := writer.Boundary()

	// Write email headers
	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

	// Build header string
	var headerStr strings.Builder
	for k, v := range headers {
		headerStr.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	headerStr.WriteString("\r\n")

	// Write HTML body part
	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"text/html; charset=UTF-8"},
		"Content-Transfer-Encoding": []string{"quoted-printable"},
	})
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %w", err)
	}
	_, err = htmlPart.Write([]byte(htmlBody))
	if err != nil {
		return fmt.Errorf("failed to write HTML body: %w", err)
	}

	// Add attachments
	for _, attachment := range attachments {
		attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":              []string{attachment.ContentType},
			"Content-Transfer-Encoding": []string{"base64"},
			"Content-Disposition":       []string{fmt.Sprintf("attachment; filename=\"%s\"", attachment.Filename)},
		})
		if err != nil {
			return fmt.Errorf("failed to create attachment part: %w", err)
		}

		// Encode attachment data to base64
		encoded := base64.StdEncoding.EncodeToString(attachment.Data)
		
		// Write in chunks of 76 characters (RFC 2045)
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			_, err = attachmentPart.Write([]byte(encoded[i:end] + "\r\n"))
			if err != nil {
				return fmt.Errorf("failed to write attachment data: %w", err)
			}
		}
	}

	// Close multipart writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Combine headers and body
	message := []byte(headerStr.String() + buf.String())

	// Send email via SMTP
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	err = smtp.SendMail(addr, auth, s.from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("[Email Service] Successfully sent email with %d attachments to %s: %s\n", len(attachments), to, subject)
	return nil
}
