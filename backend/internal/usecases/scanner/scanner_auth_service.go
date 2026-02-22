package scanner

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/pkg/jwt"
)

// ScannerAuthService handles scanner authentication and session management
type ScannerAuthService struct {
	repoManager repositories.RepositoryManager
	jwtService  jwt.Service
}

// NewScannerAuthService creates a new scanner authentication service
func NewScannerAuthService(repoManager repositories.RepositoryManager, jwtSecret string) *ScannerAuthService {
	jwtService := jwt.NewJWTService(jwtSecret, 15*time.Minute, 7*24*time.Hour, "uduxpass-scanner")
	return &ScannerAuthService{
		repoManager: repoManager,
		jwtService:  jwtService,
	}
}

// Login authenticates a scanner user and returns tokens
func (s *ScannerAuthService) Login(ctx context.Context, username, password, clientIP, userAgent string) (*entities.ScannerLoginResponse, error) {
	// Get scanner by username
	fmt.Printf("DEBUG: Attempting to find scanner with username: %s\n", username)
	scanner, err := s.repoManager.ScannerUsers().GetByUsername(ctx, username)
	if err != nil {
		fmt.Printf("DEBUG: GetByUsername failed with error: %v\n", err)
		// Record failed login attempt
		s.recordLoginAttempt(ctx, nil, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}
	
	fmt.Printf("DEBUG: Found scanner: %s, Status: %s\n", scanner.Username, scanner.Status)
	fmt.Printf("DEBUG: Scanner ID: %s\n", scanner.ID)
	fmt.Printf("DEBUG: Password hash length: %d\n", len(scanner.PasswordHash))

	// Check if account is locked
	if scanner.IsLocked() {
		fmt.Printf("DEBUG: Scanner account is locked\n")
		s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Account is temporarily locked due to multiple failed login attempts",
		}, nil
	}

	// Check if account is active
	if scanner.Status != entities.ScannerStatusActive {
		fmt.Printf("DEBUG: Scanner account is not active, status: %s\n", scanner.Status)
		s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Account is not active",
		}, nil
	}

	// Verify password using stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(scanner.PasswordHash), []byte(password)); err != nil {
		s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}
	
	fmt.Printf("DEBUG: Password verification successful!\n")

	// Generate JWT tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Record successful login
	s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, true)

	// Log login activity
	s.logActivity(ctx, scanner.ID, "login", nil, nil, nil, map[string]interface{}{
		"ip_address": clientIP,
		"user_agent": userAgent,
	})

	return &entities.ScannerLoginResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scanner:      scanner,
		ExpiresIn:    expiresIn,
		Message:      "Login successful",
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *ScannerAuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.ScannerLoginResponse, error) {
	// Validate refresh token and extract scanner ID
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Invalid refresh token",
		}, nil
	}

	scannerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Invalid token format",
		}, nil
	}

	// Get scanner details
	scanner, err := s.repoManager.ScannerUsers().GetByID(ctx, scannerID)
	if err != nil {
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Scanner not found",
		}, nil
	}

	// Check if scanner is still active
	if scanner.Status != entities.ScannerStatusActive {
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Scanner account is not active",
		}, nil
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresIn, err := s.generateTokens(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &entities.ScannerLoginResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Scanner:      scanner,
		ExpiresIn:    expiresIn,
		Message:      "Token refreshed successfully",
	}, nil
}

// Logout handles scanner logout
func (s *ScannerAuthService) Logout(ctx context.Context, scannerID uuid.UUID) error {
	// End any active sessions
	session, err := s.repoManager.ScannerUsers().GetActiveSession(ctx, scannerID)
	if err == nil && session != nil {
		s.EndSession(ctx, session.ID)
	}

	// Log logout activity
	s.logActivity(ctx, scannerID, "logout", nil, nil, nil, nil)

	return nil
}

// StartSession starts a new scanning session
func (s *ScannerAuthService) StartSession(ctx context.Context, scannerID, eventID uuid.UUID) (*entities.ScannerSession, error) {
	// Check if scanner is assigned to this event
	assignedEvents, err := s.repoManager.ScannerUsers().GetAssignedEvents(ctx, scannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check event assignment: %w", err)
	}

	isAssigned := false
	for _, event := range assignedEvents {
		if event.EventID == eventID {
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		return nil, errors.New("scanner is not assigned to this event")
	}

	// End any existing active session
	if activeSession, err := s.repoManager.ScannerUsers().GetActiveSession(ctx, scannerID); err == nil && activeSession != nil {
		s.EndSession(ctx, activeSession.ID)
	}

	// Create new session
	session := &entities.ScannerSession{
		ID:           uuid.New(),
		ScannerID:    scannerID,
		EventID:      eventID,
		StartTime:    time.Now(),
		ScansCount:   0,
		ValidScans:   0,
		InvalidScans: 0,
		TotalRevenue: 0,
		IsActive:     true,
	}

	if err := s.repoManager.ScannerUsers().CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Log session start
	resourceType := "scanner_session"
	s.logActivity(ctx, scannerID, "session_start", &session.ID, &resourceType, &session.ID, map[string]interface{}{
		"event_id": eventID,
	})

	return session, nil
}

// EndSession ends a scanning session
func (s *ScannerAuthService) EndSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.repoManager.ScannerUsers().EndSession(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}

	return nil
}

// ValidateTicket validates a ticket and records the scan
func (s *ScannerAuthService) ValidateTicket(ctx context.Context, scannerID, sessionID, eventID uuid.UUID, ticketCode string, notes *string) (*entities.TicketValidationResponse, error) {
	// Look up the actual ticket
	ticket, err := s.repoManager.Tickets().GetByQRCode(ctx, ticketCode)
	if err != nil {
		return &entities.TicketValidationResponse{
			Success:        false,
			Valid:          false,
			Message:        "Ticket not found",
			ValidationTime: time.Now(),
		}, nil
	}
	
	response := &entities.TicketValidationResponse{
		ValidationTime: time.Now(),
	}

	// Check if ticket is valid
	if ticket.Status == "active" {
		response.Success = true
		response.Valid = true
		response.Message = "Ticket validated successfully"
		standardType := "standard"
		holderName := ""
		response.TicketType = &standardType // Could be enhanced with ticket tier info
		response.HolderName = &holderName // Could be enhanced with attendee info
		response.AlreadyValidated = false

		// Record validation
		validation := &entities.TicketValidation{
			ID:                  uuid.New(),
			TicketID:            ticket.ID, // Use actual ticket ID
			ScannerID:           scannerID,
			SessionID:           sessionID,
			ValidationResult:    "valid",
			ValidationTimestamp: time.Now(),
			Notes:               notes,
		}

		if err := s.repoManager.ScannerUsers().ValidateTicket(ctx, validation); err != nil {
			return nil, fmt.Errorf("failed to record validation: %w", err)
		}

		// Update session stats
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 1, 0, s.getTicketPrice(ticketCode))

		// Log validation activity
		ticketResourceType := "ticket"
		s.logActivity(ctx, scannerID, "ticket_validation", &sessionID, &ticketResourceType, &validation.TicketID, map[string]interface{}{
			"validation_result": "valid",
			"ticket_code":       ticketCode,
			"event_id":          eventID,
		})
	} else {
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket code"
		response.AlreadyValidated = false

		// Record invalid validation
		validation := &entities.TicketValidation{
			ID:                  uuid.New(),
			TicketID:            uuid.New(),
			ScannerID:           scannerID,
			SessionID:           sessionID,
			ValidationResult:    "invalid",
			ValidationTimestamp: time.Now(),
			Notes:               notes,
		}

		s.repoManager.ScannerUsers().ValidateTicket(ctx, validation)
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)

		// Log validation activity
		ticketResourceType := "ticket"
		s.logActivity(ctx, scannerID, "ticket_validation", &sessionID, &ticketResourceType, &validation.TicketID, map[string]interface{}{
			"validation_result": "invalid",
			"ticket_code":       ticketCode,
			"event_id":          eventID,
		})
	}

	return response, nil
}

// Helper methods

func (s *ScannerAuthService) generateTokens(scanner *entities.ScannerUser) (string, string, int64, error) {
	// Generate access token
	accessToken, err := s.jwtService.GenerateAccessToken(scanner.ID, string(scanner.Role))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateRefreshToken(scanner.ID, string(scanner.Role))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, 15 * 60, nil // 15 minutes in seconds
}

func (s *ScannerAuthService) recordLoginAttempt(ctx context.Context, scannerID *uuid.UUID, clientIP, userAgent string, success bool) {
	if scannerID == nil {
		return
	}

	loginHistory := &entities.ScannerLoginHistory{
		ID:        uuid.New(),
		ScannerID: *scannerID,
		IPAddress: &clientIP,
		UserAgent: &userAgent,
		Success:   success,
		LoginAt:   time.Now(),
	}

	s.repoManager.ScannerUsers().RecordLogin(ctx, loginHistory)
}

func (s *ScannerAuthService) logActivity(ctx context.Context, scannerID uuid.UUID, action string, sessionID *uuid.UUID, resourceType *string, resourceID *uuid.UUID, details map[string]interface{}) {
	log := &entities.ScannerAuditLog{
		ID:           uuid.New(),
		ScannerID:    scannerID,
		SessionID:    sessionID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		CreatedAt:    time.Now(),
	}

	s.repoManager.ScannerUsers().LogActivity(ctx, log)
}

// Demo validation helpers
func (s *ScannerAuthService) isDemoTicketCode(code string) bool {
	// Demo validation logic - accept codes that contain certain patterns
	code = strings.ToUpper(code)
	validPatterns := []string{
		"UDUXPASS",
		"VIP",
		"PREMIUM",
		"GENERAL",
		"DEMO",
		"TEST",
	}

	for _, pattern := range validPatterns {
		if strings.Contains(code, pattern) {
			return true
		}
	}

	return false
}

func (s *ScannerAuthService) extractTicketType(code string) *string {
	code = strings.ToUpper(code)
	if strings.Contains(code, "VIP") {
		ticketType := "VIP"
		return &ticketType
	}
	if strings.Contains(code, "PREMIUM") {
		ticketType := "Premium"
		return &ticketType
	}
	ticketType := "General Admission"
	return &ticketType
}

func (s *ScannerAuthService) extractHolderName(code string) *string {
	// Demo logic - return a demo name
	holderName := "Demo User"
	return &holderName
}

func (s *ScannerAuthService) getTicketPrice(code string) float64 {
	code = strings.ToUpper(code)
	if strings.Contains(code, "VIP") {
		return 10000.0 // NGN 10,000
	}
	if strings.Contains(code, "PREMIUM") {
		return 5000.0 // NGN 5,000
	}
	return 2500.0 // NGN 2,500
}

