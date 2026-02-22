# uduXPass Platform - Comprehensive Testing Checklist

## Phase 1: Backend API & QR Code Generation
- [ ] Verify all API endpoints are registered and accessible
- [ ] Test event creation via API
- [ ] Test ticket tier creation
- [ ] Test order creation
- [ ] **CRITICAL: Verify QR code is generated when ticket is created**
- [ ] Verify QR code format and data structure
- [ ] Test ticket retrieval with QR code data

## Phase 2: Complete User Flow
- [ ] Test user registration via API
- [ ] Test user login and token generation
- [ ] Test event browsing
- [ ] Test ticket purchase flow
- [ ] Verify order confirmation
- [ ] Verify ticket delivery to user
- [ ] **CRITICAL: Verify user receives ticket with valid QR code**

## Phase 3: Scanner App Functionality
- [ ] Verify scanner user creation
- [ ] Test scanner login
- [ ] Test scanner session management
- [ ] **CRITICAL: Test QR code scanning functionality**
- [ ] Test ticket validation logic
- [ ] Verify validation records are created
- [ ] Test validation history retrieval

## Phase 4: Anti-Reuse Protection
- [ ] Validate ticket first time (should succeed)
- [ ] Attempt to validate same ticket again (should fail)
- [ ] Verify database constraints prevent duplicate validation
- [ ] Verify ticket status changes after validation
- [ ] Test edge cases (expired tickets, invalid QR codes)

## Phase 5: Integration Testing
- [ ] Test complete flow: Create event → User buys → Scanner validates
- [ ] Verify data consistency across all tables
- [ ] Test error handling and edge cases
- [ ] Verify API response formats
- [ ] Test concurrent operations

## Critical Issues to Investigate
- [ ] QR code generation library/method
- [ ] QR code data format and encoding
- [ ] Scanner app camera integration
- [ ] QR code decoding in scanner app
- [ ] Network connectivity handling
- [ ] Offline mode functionality
