package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// OTPService defines the interface for OTP operations
type OTPService interface {
	GenerateOTP(length int) (string, error)
	GenerateNumericOTP(length int) (string, error)
	GenerateAlphanumericOTP(length int) (string, error)
	IsValidOTP(otp string, length int) bool
}

// DefaultOTPService implements OTP generation and validation
type DefaultOTPService struct{}

// NewOTPService creates a new OTP service
func NewOTPService() *DefaultOTPService {
	return &DefaultOTPService{}
}

// GenerateOTP generates a numeric OTP of specified length
func (s *DefaultOTPService) GenerateOTP(length int) (string, error) {
	return s.GenerateNumericOTP(length)
}

// GenerateNumericOTP generates a numeric OTP
func (s *DefaultOTPService) GenerateNumericOTP(length int) (string, error) {
	if length < 4 {
		length = 4
	}
	if length > 10 {
		length = 10
	}
	
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		otp[i] = byte('0' + num.Int64())
	}
	
	return string(otp), nil
}

// GenerateAlphanumericOTP generates an alphanumeric OTP
func (s *DefaultOTPService) GenerateAlphanumericOTP(length int) (string, error) {
	if length < 4 {
		length = 4
	}
	if length > 20 {
		length = 20
	}
	
	// Use only uppercase letters and digits to avoid confusion
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	otp := make([]byte, length)
	
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		otp[i] = chars[num.Int64()]
	}
	
	return string(otp), nil
}

// IsValidOTP validates an OTP format
func (s *DefaultOTPService) IsValidOTP(otp string, length int) bool {
	if len(otp) != length {
		return false
	}
	
	// Check if all characters are digits
	for _, char := range otp {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	return true
}

// OTPManager manages OTP generation, storage, and verification
type OTPManager struct {
	otpService OTPService
	storage    map[string]*OTPRecord // In production, use Redis or database
}

// OTPRecord represents an OTP record
type OTPRecord struct {
	OTP       string
	Purpose   string
	ExpiresAt time.Time
	Attempts  int
	MaxAttempts int
	CreatedAt time.Time
}

// OTPPurpose defines OTP purposes
type OTPPurpose string

const (
	OTPPurposeLogin          OTPPurpose = "login"
	OTPPurposeRegistration   OTPPurpose = "registration"
	OTPPurposePasswordReset  OTPPurpose = "password_reset"
	OTPPurposePhoneVerification OTPPurpose = "phone_verification"
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
)

// OTPConfig represents OTP configuration
type OTPConfig struct {
	Length      int
	ExpiryTime  time.Duration
	MaxAttempts int
}

// NewOTPManager creates a new OTP manager
func NewOTPManager(otpService OTPService) *OTPManager {
	return &OTPManager{
		otpService: otpService,
		storage:    make(map[string]*OTPRecord),
	}
}

// GenerateAndStore generates an OTP and stores it
func (m *OTPManager) GenerateAndStore(identifier string, purpose OTPPurpose, config OTPConfig) (string, error) {
	// Set default values
	if config.Length == 0 {
		config.Length = 6
	}
	if config.ExpiryTime == 0 {
		config.ExpiryTime = 5 * time.Minute
	}
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}
	
	// Generate OTP
	otp, err := m.otpService.GenerateNumericOTP(config.Length)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}
	
	// Create record
	record := &OTPRecord{
		OTP:         otp,
		Purpose:     string(purpose),
		ExpiresAt:   time.Now().Add(config.ExpiryTime),
		Attempts:    0,
		MaxAttempts: config.MaxAttempts,
		CreatedAt:   time.Now(),
	}
	
	// Store record (key format: identifier:purpose)
	key := fmt.Sprintf("%s:%s", identifier, purpose)
	m.storage[key] = record
	
	return otp, nil
}

// VerifyOTP verifies an OTP
func (m *OTPManager) VerifyOTP(identifier string, purpose OTPPurpose, otp string) (bool, error) {
	key := fmt.Sprintf("%s:%s", identifier, purpose)
	
	record, exists := m.storage[key]
	if !exists {
		return false, fmt.Errorf("OTP not found")
	}
	
	// Check if OTP has expired
	if time.Now().After(record.ExpiresAt) {
		delete(m.storage, key)
		return false, fmt.Errorf("OTP has expired")
	}
	
	// Check if max attempts exceeded
	if record.Attempts >= record.MaxAttempts {
		delete(m.storage, key)
		return false, fmt.Errorf("maximum attempts exceeded")
	}
	
	// Increment attempts
	record.Attempts++
	
	// Verify OTP
	if record.OTP != otp {
		return false, fmt.Errorf("invalid OTP")
	}
	
	// OTP is valid, remove it
	delete(m.storage, key)
	return true, nil
}

// CleanupExpired removes expired OTP records
func (m *OTPManager) CleanupExpired() int {
	count := 0
	now := time.Now()
	
	for key, record := range m.storage {
		if now.After(record.ExpiresAt) {
			delete(m.storage, key)
			count++
		}
	}
	
	return count
}

// GetOTPInfo returns information about an OTP without verifying it
func (m *OTPManager) GetOTPInfo(identifier string, purpose OTPPurpose) (*OTPInfo, error) {
	key := fmt.Sprintf("%s:%s", identifier, purpose)
	
	record, exists := m.storage[key]
	if !exists {
		return nil, fmt.Errorf("OTP not found")
	}
	
	return &OTPInfo{
		Purpose:     purpose,
		ExpiresAt:   record.ExpiresAt,
		Attempts:    record.Attempts,
		MaxAttempts: record.MaxAttempts,
		CreatedAt:   record.CreatedAt,
		IsExpired:   time.Now().After(record.ExpiresAt),
	}, nil
}

// OTPInfo represents OTP information
type OTPInfo struct {
	Purpose     OTPPurpose `json:"purpose"`
	ExpiresAt   time.Time  `json:"expires_at"`
	Attempts    int        `json:"attempts"`
	MaxAttempts int        `json:"max_attempts"`
	CreatedAt   time.Time  `json:"created_at"`
	IsExpired   bool       `json:"is_expired"`
}

// SMSService defines the interface for SMS operations
type SMSService interface {
	SendSMS(phone, message string) error
	SendOTP(phone, otp string, purpose OTPPurpose) error
}

// MockSMSService implements a mock SMS service for testing
type MockSMSService struct {
	shouldFail bool
	sentMessages []SMSMessage
}

// SMSMessage represents a sent SMS message
type SMSMessage struct {
	Phone   string
	Message string
	SentAt  time.Time
}

// NewMockSMSService creates a new mock SMS service
func NewMockSMSService(shouldFail bool) *MockSMSService {
	return &MockSMSService{
		shouldFail:   shouldFail,
		sentMessages: make([]SMSMessage, 0),
	}
}

// SendSMS mock implementation
func (m *MockSMSService) SendSMS(phone, message string) error {
	if m.shouldFail {
		return fmt.Errorf("mock SMS send failed")
	}
	
	m.sentMessages = append(m.sentMessages, SMSMessage{
		Phone:   phone,
		Message: message,
		SentAt:  time.Now(),
	})
	
	return nil
}

// SendOTP mock implementation
func (m *MockSMSService) SendOTP(phone, otp string, purpose OTPPurpose) error {
	message := fmt.Sprintf("Your uduXPass verification code is: %s. Valid for 5 minutes.", otp)
	return m.SendSMS(phone, message)
}

// GetSentMessages returns all sent messages (for testing)
func (m *MockSMSService) GetSentMessages() []SMSMessage {
	return m.sentMessages
}

// ClearSentMessages clears all sent messages (for testing)
func (m *MockSMSService) ClearSentMessages() {
	m.sentMessages = make([]SMSMessage, 0)
}

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendOTP(email, otp string, purpose OTPPurpose) error
}

// MockEmailService implements a mock email service for testing
type MockEmailService struct {
	shouldFail bool
	sentEmails []EmailMessage
}

// EmailMessage represents a sent email message
type EmailMessage struct {
	To      string
	Subject string
	Body    string
	SentAt  time.Time
}

// NewMockEmailService creates a new mock email service
func NewMockEmailService(shouldFail bool) *MockEmailService {
	return &MockEmailService{
		shouldFail: shouldFail,
		sentEmails: make([]EmailMessage, 0),
	}
}

// SendEmail mock implementation
func (m *MockEmailService) SendEmail(to, subject, body string) error {
	if m.shouldFail {
		return fmt.Errorf("mock email send failed")
	}
	
	m.sentEmails = append(m.sentEmails, EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		SentAt:  time.Now(),
	})
	
	return nil
}

// SendOTP mock implementation
func (m *MockEmailService) SendOTP(email, otp string, purpose OTPPurpose) error {
	subject := "uduXPass Verification Code"
	body := fmt.Sprintf(`
		<h2>uduXPass Verification Code</h2>
		<p>Your verification code is: <strong>%s</strong></p>
		<p>This code is valid for 5 minutes.</p>
		<p>If you didn't request this code, please ignore this email.</p>
	`, otp)
	
	return m.SendEmail(email, subject, body)
}

// GetSentEmails returns all sent emails (for testing)
func (m *MockEmailService) GetSentEmails() []EmailMessage {
	return m.sentEmails
}

// ClearSentEmails clears all sent emails (for testing)
func (m *MockEmailService) ClearSentEmails() {
	m.sentEmails = make([]EmailMessage, 0)
}

// MockOTPService implements a mock OTP service for testing
type MockOTPService struct {
	shouldFail bool
	fixedOTP   string
}

// NewMockOTPService creates a new mock OTP service
func NewMockOTPService(shouldFail bool, fixedOTP string) *MockOTPService {
	if fixedOTP == "" {
		fixedOTP = "123456"
	}
	return &MockOTPService{
		shouldFail: shouldFail,
		fixedOTP:   fixedOTP,
	}
}

// GenerateOTP mock implementation
func (m *MockOTPService) GenerateOTP(length int) (string, error) {
	if m.shouldFail {
		return "", fmt.Errorf("mock OTP generation failed")
	}
	return m.fixedOTP[:length], nil
}

// GenerateNumericOTP mock implementation
func (m *MockOTPService) GenerateNumericOTP(length int) (string, error) {
	return m.GenerateOTP(length)
}

// GenerateAlphanumericOTP mock implementation
func (m *MockOTPService) GenerateAlphanumericOTP(length int) (string, error) {
	if m.shouldFail {
		return "", fmt.Errorf("mock OTP generation failed")
	}
	return "ABC123"[:length], nil
}

// IsValidOTP mock implementation
func (m *MockOTPService) IsValidOTP(otp string, length int) bool {
	return len(otp) == length
}

