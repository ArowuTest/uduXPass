package scanner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	pkgjwt "github.com/uduxpass/backend/pkg/jwt"
)

// ticketJWTClaims mirrors the claims structure used by the payment service when
// signing ticket QR codes. The scanner verifies the signature before any DB lookup.
type ticketJWTClaims struct {
	TicketID     string `json:"tid"`
	EventID      string `json:"eid"`
	SerialNumber string `json:"sn"`
	OrderLineID  string `json:"olid"`
	jwt.RegisteredClaims
}

// ScannerAuthService handles scanner authentication, session management, and ticket validation.
type ScannerAuthService struct {
	repoManager     repositories.RepositoryManager
	jwtService      pkgjwt.Service
	ticketJWTSecret []byte // same secret used by PaymentService to sign ticket JWTs
}

// NewScannerAuthService creates a new scanner authentication service.
// jwtSecret is used both for scanner session JWTs and for verifying ticket QR code JWTs.
func NewScannerAuthService(repoManager repositories.RepositoryManager, jwtSecret string) *ScannerAuthService {
	jwtService := pkgjwt.NewJWTService(jwtSecret, 15*time.Minute, 7*24*time.Hour, "uduxpass-scanner")
	return &ScannerAuthService{
		repoManager:     repoManager,
		jwtService:      jwtService,
		ticketJWTSecret: []byte(jwtSecret),
	}
}

// Login authenticates a scanner user and returns JWT tokens.
func (s *ScannerAuthService) Login(ctx context.Context, username, password, clientIP, userAgent string) (*entities.ScannerLoginResponse, error) {
	fmt.Printf("DEBUG: Attempting to find scanner with username: %s\n", username)
	scanner, err := s.repoManager.ScannerUsers().GetByUsername(ctx, username)
	if err != nil {
		fmt.Printf("DEBUG: GetByUsername failed with error: %v\n", err)
		s.recordLoginAttempt(ctx, nil, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	fmt.Printf("DEBUG: Found scanner: %s, Status: %s\n", scanner.Username, scanner.Status)

	// Check if account is locked
	if scanner.IsLocked() {
		s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Account is temporarily locked due to multiple failed login attempts",
		}, nil
	}

	// Check if account is active
	if scanner.Status != entities.ScannerStatusActive {
		s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, false)
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Account is not active",
		}, nil
	}

	// Verify password
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

	s.recordLoginAttempt(ctx, &scanner.ID, clientIP, userAgent, true)
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

// RefreshToken generates new tokens using a refresh token.
func (s *ScannerAuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.ScannerLoginResponse, error) {
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

	scanner, err := s.repoManager.ScannerUsers().GetByID(ctx, scannerID)
	if err != nil {
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Scanner not found",
		}, nil
	}

	if scanner.Status != entities.ScannerStatusActive {
		return &entities.ScannerLoginResponse{
			Success: false,
			Message: "Scanner account is not active",
		}, nil
	}

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

// Logout handles scanner logout and ends any active session.
func (s *ScannerAuthService) Logout(ctx context.Context, scannerID uuid.UUID) error {
	session, err := s.repoManager.ScannerUsers().GetActiveSession(ctx, scannerID)
	if err == nil && session != nil {
		s.EndSession(ctx, session.ID)
	}
	s.logActivity(ctx, scannerID, "logout", nil, nil, nil, nil)
	return nil
}

// StartSession starts a new scanning session for an event.
// The scanner must be assigned to the event before a session can be started.
func (s *ScannerAuthService) StartSession(ctx context.Context, scannerID, eventID uuid.UUID) (*entities.ScannerSession, error) {
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

	// End any existing active session before starting a new one
	if activeSession, err := s.repoManager.ScannerUsers().GetActiveSession(ctx, scannerID); err == nil && activeSession != nil {
		s.EndSession(ctx, activeSession.ID)
	}

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

	resourceType := "scanner_session"
	s.logActivity(ctx, scannerID, "session_start", &session.ID, &resourceType, &session.ID, map[string]interface{}{
		"event_id": eventID,
	})

	return session, nil
}

// EndSession ends a scanning session.
func (s *ScannerAuthService) EndSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.repoManager.ScannerUsers().EndSession(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}
	return nil
}

// ValidateTicket validates a ticket QR code and records the scan result.
//
// Validation flow (enterprise-grade, per SRS FR-3.2):
//  1. Verify the JWT signature of the QR code data — reject tampered codes immediately
//  2. Extract ticket_id and event_id from the JWT claims
//  3. Verify the event_id in the JWT matches the scanner's active session event
//  4. Look up the ticket in the DB by ID
//  5. Check ticket status: active → valid; redeemed → duplicate; voided/other → invalid
//  6. On valid scan: atomically mark the ticket as redeemed in the DB
//  7. Record the validation event in ticket_validations
//  8. Update session statistics
func (s *ScannerAuthService) ValidateTicket(ctx context.Context, scannerID, sessionID, eventID uuid.UUID, ticketCode string, notes *string) (*entities.TicketValidationResponse, error) {
	response := &entities.TicketValidationResponse{
		ValidationTime: time.Now(),
	}

	// --- Step 1: Verify JWT signature ---
	claims, err := s.verifyTicketJWT(ticketCode)
	if err != nil {
		// JWT is invalid or tampered — record as invalid scan and return
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket: QR code signature verification failed"
		response.AlreadyValidated = false

		s.recordInvalidScan(ctx, scannerID, sessionID, notes, "invalid_signature")
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil
	}

	// --- Step 2: Parse ticket ID and event ID from claims ---
	ticketID, err := uuid.Parse(claims.TicketID)
	if err != nil {
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket: malformed ticket ID in QR code"
		s.recordInvalidScan(ctx, scannerID, sessionID, notes, "malformed_claims")
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil
	}

	claimedEventID, err := uuid.Parse(claims.EventID)
	if err != nil {
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket: malformed event ID in QR code"
		s.recordInvalidScan(ctx, scannerID, sessionID, notes, "malformed_claims")
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil
	}

	// --- Step 3: Verify the ticket belongs to this event ---
	if claimedEventID != eventID {
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket: this ticket is for a different event"
		s.recordValidationEvent(ctx, ticketID, scannerID, sessionID, "wrong_event", notes)
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil
	}

	// --- Step 4: Look up ticket in the database ---
	ticket, err := s.repoManager.Tickets().GetByID(ctx, ticketID)
	if err != nil {
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket: ticket not found in system"
		s.recordInvalidScan(ctx, scannerID, sessionID, notes, "not_found")
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil
	}

	// --- Step 5: Check ticket status ---
	switch ticket.Status {
	case entities.TicketStatusRedeemed:
		// Duplicate scan — ticket already used
		response.Success = true
		response.Valid = false
		response.AlreadyValidated = true
		response.Message = fmt.Sprintf("Ticket already redeemed at %s", ticket.RedeemedAt.Format("02 Jan 2006 15:04:05"))
		s.recordValidationEvent(ctx, ticketID, scannerID, sessionID, "already_redeemed", notes)
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil

	case entities.TicketStatusVoided:
		response.Success = true
		response.Valid = false
		response.Message = "Invalid ticket: this ticket has been voided"
		s.recordValidationEvent(ctx, ticketID, scannerID, sessionID, "voided", notes)
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil

	case entities.TicketStatusActive:
		// Valid ticket — proceed to redeem

	default:
		response.Success = true
		response.Valid = false
		response.Message = fmt.Sprintf("Invalid ticket: unexpected status '%s'", ticket.Status)
		s.recordValidationEvent(ctx, ticketID, scannerID, sessionID, "invalid_status", notes)
		s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 0, 1, 0)
		return response, nil
	}

	// --- Step 6: Mark ticket as redeemed ---
	redeemedBy := scannerID.String()
	if err := ticket.Redeem(redeemedBy); err != nil {
		// Should not happen since we already checked status == active, but handle defensively
		return nil, fmt.Errorf("failed to redeem ticket %s: %w", ticketID, err)
	}

	if err := s.repoManager.Tickets().MarkRedeemed(ctx, ticketID, redeemedBy); err != nil {
		return nil, fmt.Errorf("failed to persist ticket redemption for %s: %w", ticketID, err)
	}

	// --- Step 7: Record the validation event ---
	s.recordValidationEvent(ctx, ticketID, scannerID, sessionID, "valid", notes)

	// --- Step 8: Update session statistics ---
	s.repoManager.ScannerUsers().UpdateSessionStats(ctx, sessionID, 1, 1, 0, 0)

	// Log audit trail
	ticketResourceType := "ticket"
	s.logActivity(ctx, scannerID, "ticket_validation", &sessionID, &ticketResourceType, &ticketID, map[string]interface{}{
		"validation_result": "valid",
		"serial_number":     claims.SerialNumber,
		"event_id":          eventID,
	})

	// Build success response
	response.Success = true
	response.Valid = true
	response.AlreadyValidated = false
	response.Message = "Ticket validated successfully"
	serialNumber := claims.SerialNumber
	response.SerialNumber = &serialNumber

	return response, nil
}

// verifyTicketJWT parses and verifies the HMAC-SHA256 signature of a ticket QR code JWT.
// Returns the claims on success, or an error if the token is invalid or tampered.
func (s *ScannerAuthService) verifyTicketJWT(tokenString string) (*ticketJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ticketJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Enforce HMAC signing method — reject asymmetric key attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.ticketJWTSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("ticket JWT verification failed: %w", err)
	}

	claims, ok := token.Claims.(*ticketJWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid ticket JWT claims")
	}

	// Validate required claims are present
	if claims.TicketID == "" || claims.EventID == "" {
		return nil, fmt.Errorf("ticket JWT missing required claims")
	}

	return claims, nil
}

// recordValidationEvent records a ticket validation event in the ticket_validations table.
func (s *ScannerAuthService) recordValidationEvent(ctx context.Context, ticketID, scannerID, sessionID uuid.UUID, result string, notes *string) {
	validation := &entities.TicketValidation{
		ID:                  uuid.New(),
		TicketID:            ticketID,
		ScannerID:           scannerID,
		SessionID:           sessionID,
		ValidationResult:    result,
		ValidationTimestamp: time.Now(),
		Notes:               notes,
	}
	if err := s.repoManager.ScannerUsers().ValidateTicket(ctx, validation); err != nil {
		fmt.Printf("Warning: failed to record validation event for ticket %s: %v\n", ticketID, err)
	}
}

// recordInvalidScan records a validation event for a scan that failed before a ticket ID was known.
// Uses a nil UUID for the ticket ID since no valid ticket was identified.
func (s *ScannerAuthService) recordInvalidScan(ctx context.Context, scannerID, sessionID uuid.UUID, notes *string, result string) {
	s.recordValidationEvent(ctx, uuid.Nil, scannerID, sessionID, result, notes)
}

// Helper methods

func (s *ScannerAuthService) generateTokens(scanner *entities.ScannerUser) (string, string, int64, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(scanner.ID, string(scanner.Role))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

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
