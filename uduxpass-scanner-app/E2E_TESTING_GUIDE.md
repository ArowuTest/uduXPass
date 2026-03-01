# uduXPass Scanner App - End-to-End Testing Guide

## Overview

This document provides a comprehensive testing checklist for the uduXPass scanner app. It covers all user flows, edge cases, and integration points with the backend.

---

## Prerequisites

### Backend Setup

1. ✅ **Backend Running**: uduXPass backend must be running on `http://localhost:3000`
2. ✅ **Database Seeded**: Database must have seed data (admin, events, users, scanners, tickets)
3. ✅ **CORS Configured**: Backend must allow requests from frontend origin

### Frontend Setup

1. ✅ **Dependencies Installed**: Run `pnpm install`
2. ✅ **Environment Configured**: Set `VITE_API_BASE_URL` in `.env` or `api.ts`
3. ✅ **Dev Server Running**: Run `pnpm dev`

---

## Test Scenarios

### 1. Scanner Authentication Flow

#### Test 1.1: Successful Login
**Steps:**
1. Navigate to `/login`
2. Enter valid scanner credentials:
   - Email: `scanner1`
   - Password: `Scanner@123!`
3. Click "Login to Scanner"

**Expected Results:**
- ✅ Loading spinner appears on button
- ✅ Success toast notification appears
- ✅ Redirected to `/dashboard`
- ✅ Scanner name appears in header
- ✅ Token saved to localStorage

**Edge Cases:**
- ❌ Empty email/password → Error toast "Please enter both email and password"
- ❌ Invalid credentials → Error toast "Invalid credentials"
- ❌ Backend offline → Error toast "Network error"

#### Test 1.2: Protected Routes
**Steps:**
1. Clear localStorage (logout)
2. Try to navigate directly to `/dashboard`

**Expected Results:**
- ✅ Redirected to `/login`
- ✅ After login, can access `/dashboard`

#### Test 1.3: Logout
**Steps:**
1. Login successfully
2. Click "Logout" button in dashboard header

**Expected Results:**
- ✅ Success toast "Logged out successfully"
- ✅ Redirected to `/login`
- ✅ Token removed from localStorage
- ✅ Cannot access protected routes

---

### 2. Dashboard Flow

#### Test 2.1: Dashboard with No Active Session
**Steps:**
1. Login as scanner
2. View dashboard (no active sessions)

**Expected Results:**
- ✅ Scanner name displayed in header
- ✅ "No active scanning session" message shown
- ✅ "Start New Session" button visible
- ✅ "Scan Ticket" button disabled
- ✅ "Session History" button visible

#### Test 2.2: Dashboard with Active Session
**Steps:**
1. Create a scanning session
2. Return to dashboard

**Expected Results:**
- ✅ Active session card displayed with:
  - Event name
  - Event date
  - Location
- ✅ Statistics cards showing:
  - Total tickets scanned
  - Valid tickets (green)
  - Invalid tickets (red)
- ✅ "Scan Ticket" button enabled (green)

---

### 3. Create Session Flow

#### Test 3.1: Successful Session Creation
**Steps:**
1. From dashboard, click "Start New Session"
2. Select an event from dropdown
3. Enter location: "Main Entrance"
4. Optionally add notes: "Test session"
5. Click "Start Scanning Session"

**Expected Results:**
- ✅ Events loaded in dropdown
- ✅ Loading spinner during creation
- ✅ Success toast "Scanning session created successfully!"
- ✅ Redirected to `/dashboard`
- ✅ Active session now visible on dashboard

**Edge Cases:**
- ❌ No event selected → Error toast "Please select an event"
- ❌ Empty location → Error toast "Please enter a location"
- ❌ No events available → Dropdown shows "No events available"
- ❌ Backend error → Error toast with message

---

### 4. QR Scanner Flow

#### Test 4.1: Camera Access Granted
**Steps:**
1. From dashboard with active session, click "Scan Ticket"
2. Allow camera access when prompted

**Expected Results:**
- ✅ Camera viewfinder appears full-screen
- ✅ Animated blue scanning frame overlay visible
- ✅ Instructions card at bottom: "Position QR code within the frame"
- ✅ "Enter Ticket Code Manually" link visible
- ✅ Back button in top bar
- ✅ Info icon in top bar

#### Test 4.2: Camera Access Denied
**Steps:**
1. Navigate to `/scan`
2. Deny camera access when prompted

**Expected Results:**
- ✅ Error message displayed: "Camera Access Required"
- ✅ Error details shown
- ✅ "Try Again" button visible
- ✅ Manual entry fallback available

#### Test 4.3: Successful QR Scan (Valid Ticket)
**Steps:**
1. Open scanner
2. Point camera at valid ticket QR code
3. Wait for auto-scan

**Expected Results:**
- ✅ Haptic feedback (vibration) on scan
- ✅ Scanner stops automatically
- ✅ Redirected to `/validation-success`
- ✅ Green gradient background
- ✅ Animated checkmark icon
- ✅ "Valid Ticket" heading
- ✅ Ticket details card showing:
  - Attendee name
  - Ticket type
  - Event name
  - Location
  - Scan timestamp
- ✅ "Scan Next Ticket" button
- ✅ Auto-redirect after 5 seconds

#### Test 4.4: Successful QR Scan (Invalid Ticket)
**Steps:**
1. Open scanner
2. Point camera at invalid/already-scanned ticket QR code

**Expected Results:**
- ✅ Haptic feedback (different pattern)
- ✅ Scanner stops automatically
- ✅ Redirected to `/validation-error`
- ✅ Red gradient background
- ✅ Warning triangle icon
- ✅ "Invalid Ticket" heading
- ✅ Error details card showing:
  - Reason for rejection
  - Previous scan info (if applicable)
  - Ticket ID
- ✅ "Override & Allow Entry" button (outline)
- ✅ "Scan Next Ticket" button (solid)

#### Test 4.5: Invalid QR Code
**Steps:**
1. Open scanner
2. Point camera at non-ticket QR code (e.g., URL)

**Expected Results:**
- ✅ Error toast "Failed to validate ticket"
- ✅ Scanner restarts automatically after 1 second

#### Test 4.6: Network Error During Validation
**Steps:**
1. Open scanner
2. Disconnect network
3. Scan a ticket QR code

**Expected Results:**
- ✅ Error toast "Failed to validate ticket"
- ✅ Scanner restarts after error

---

### 5. Validation Result Flows

#### Test 5.1: Valid Ticket - Scan Next
**Steps:**
1. Scan a valid ticket
2. View success screen
3. Click "Scan Next Ticket"

**Expected Results:**
- ✅ Redirected back to `/scan`
- ✅ Camera starts automatically
- ✅ Ready to scan next ticket

#### Test 5.2: Valid Ticket - Auto Redirect
**Steps:**
1. Scan a valid ticket
2. Wait 5 seconds without clicking

**Expected Results:**
- ✅ Auto-redirected to `/scan` after 5 seconds
- ✅ "Redirecting in 5 seconds..." message visible

#### Test 5.3: Invalid Ticket - Override (Admin)
**Steps:**
1. Scan an invalid ticket
2. Click "Override & Allow Entry"

**Expected Results:**
- ✅ Alert dialog appears (placeholder)
- ✅ Message: "Override functionality requires admin permissions"
- ❌ TODO: Implement actual override logic

#### Test 5.4: Invalid Ticket - Scan Next
**Steps:**
1. Scan an invalid ticket
2. Click "Scan Next Ticket"

**Expected Results:**
- ✅ Redirected back to `/scan`
- ✅ Camera starts automatically

---

### 6. Session History Flow

#### Test 6.1: View Session History
**Steps:**
1. From dashboard, click "Session History"

**Expected Results:**
- ✅ List of all scanning sessions displayed
- ✅ Each session card shows:
  - Event name with icon
  - Date and time range
  - Location
  - Statistics (scanned, valid, invalid)
  - Status badge (Active/Completed)
  - Notes (if any)
- ✅ Active sessions at top with blue badge
- ✅ Completed sessions with green badge

#### Test 6.2: Empty Session History
**Steps:**
1. Login as new scanner with no sessions
2. Navigate to `/history`

**Expected Results:**
- ✅ Empty state displayed
- ✅ Icon and message: "No sessions found"
- ✅ Helpful text: "Start a new session to begin scanning tickets"

---

### 7. Mobile Responsiveness

#### Test 7.1: Portrait Mode (Primary)
**Steps:**
1. Open app on mobile device in portrait mode
2. Navigate through all screens

**Expected Results:**
- ✅ All screens fit within viewport
- ✅ No horizontal scrolling
- ✅ Touch targets ≥ 44x44px
- ✅ Text readable without zooming
- ✅ Buttons within thumb reach

#### Test 7.2: Landscape Mode
**Steps:**
1. Rotate device to landscape
2. Test scanner screen

**Expected Results:**
- ✅ Scanner adjusts to landscape
- ✅ Viewfinder fills screen appropriately
- ✅ Instructions still visible

#### Test 7.3: Small Screens (iPhone SE)
**Steps:**
1. Test on small screen device (320px width)

**Expected Results:**
- ✅ All content visible
- ✅ No layout breaks
- ✅ Buttons still accessible

---

### 8. Performance & UX

#### Test 8.1: Loading States
**Steps:**
1. Test all actions with slow network (throttle to 3G)

**Expected Results:**
- ✅ Login button shows spinner during auth
- ✅ Dashboard shows loading for sessions
- ✅ Create session button shows spinner
- ✅ Event dropdown shows loading state
- ✅ Session history shows loading spinner

#### Test 8.2: Error Handling
**Steps:**
1. Test all actions with backend offline

**Expected Results:**
- ✅ Clear error messages for all failures
- ✅ No silent failures
- ✅ Retry options where appropriate
- ✅ Graceful degradation

#### Test 8.3: Animations
**Steps:**
1. Navigate through all screens
2. Observe animations

**Expected Results:**
- ✅ Scanning frame pulses smoothly
- ✅ Success checkmark animates in
- ✅ Screen transitions smooth
- ✅ No janky animations
- ✅ Respects prefers-reduced-motion

---

### 9. PWA Features

#### Test 9.1: Installability
**Steps:**
1. Open app in mobile browser
2. Look for install prompt

**Expected Results:**
- ✅ Install prompt appears (iOS/Android)
- ✅ Can add to home screen
- ✅ App icon appears on home screen
- ✅ Opens in standalone mode

#### Test 9.2: Offline Behavior
**Steps:**
1. Install PWA
2. Disconnect network
3. Open app

**Expected Results:**
- ✅ App loads (cached)
- ✅ Shows offline indicator
- ✅ Graceful error messages for API calls

---

### 10. Security

#### Test 10.1: Token Persistence
**Steps:**
1. Login successfully
2. Close browser
3. Reopen app

**Expected Results:**
- ✅ Still logged in
- ✅ Token persists in localStorage
- ✅ Can access protected routes

#### Test 10.2: Token Expiration
**Steps:**
1. Login successfully
2. Manually expire token in backend
3. Try to access protected route

**Expected Results:**
- ✅ 401 error from backend
- ✅ Redirected to login
- ✅ Token cleared from localStorage

---

## Edge Cases & Error Scenarios

### Camera Issues
- ✅ Camera permission denied → Show error and manual entry
- ✅ No camera available → Show manual entry only
- ✅ Camera in use by another app → Show error message
- ✅ Poor lighting → QR code not detected → Manual entry fallback

### Network Issues
- ✅ Backend offline → Clear error messages
- ✅ Slow network → Loading states
- ✅ Timeout → Retry option
- ✅ Intermittent connection → Queue scans for later (future)

### Data Issues
- ✅ No events available → Disable session creation
- ✅ No active session → Disable scan button
- ✅ Invalid QR data → Error message
- ✅ Malformed API response → Graceful error

### User Flow Issues
- ✅ Direct URL access → Proper redirects
- ✅ Back button navigation → Maintains state
- ✅ Refresh during scan → Restart scanner
- ✅ Multiple tabs → Independent sessions

---

## Integration Testing Checklist

### Backend API Integration
- ✅ POST `/api/v1/scanner/auth/login` - Authentication works
- ✅ GET `/api/v1/scanner/events` - Events load correctly
- ✅ POST `/api/v1/scanner/sessions` - Sessions create successfully
- ✅ GET `/api/v1/scanner/sessions` - Sessions list correctly
- ✅ GET `/api/v1/scanner/sessions?status=active` - Active sessions filter
- ✅ GET `/api/v1/scanner/sessions/:id/stats` - Statistics load
- ✅ POST `/api/v1/scanner/validate` - Validation works correctly
- ✅ PATCH `/api/v1/scanner/sessions/:id/end` - End session works

### Data Flow
- ✅ Login → Token saved → API calls authenticated
- ✅ Create session → Session ID stored → Used in validation
- ✅ Scan ticket → Validation request → Result displayed
- ✅ Statistics update → Dashboard reflects changes

---

## Performance Benchmarks

### Load Times
- ✅ Initial page load < 2 seconds
- ✅ Route transitions < 300ms
- ✅ API calls < 1 second (normal network)
- ✅ QR scan detection < 500ms

### Lighthouse Scores (Target)
- ✅ Performance: > 90
- ✅ Accessibility: > 95
- ✅ Best Practices: > 90
- ✅ SEO: > 90
- ✅ PWA: 100

---

## Browser Compatibility

### Mobile
- ✅ iOS Safari (latest 2 versions)
- ✅ Android Chrome (latest 2 versions)
- ✅ iOS Chrome (latest)
- ✅ Android Firefox (latest)

### Desktop (Secondary)
- ✅ Chrome (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ✅ Edge (latest)

---

## Test Results Summary

| Category | Tests | Passed | Failed | Notes |
|---|---|---|---|---|
| Authentication | 3 | - | - | To be tested |
| Dashboard | 2 | - | - | To be tested |
| Create Session | 1 | - | - | To be tested |
| QR Scanner | 6 | - | - | To be tested |
| Validation Results | 4 | - | - | To be tested |
| Session History | 2 | - | - | To be tested |
| Mobile Responsive | 3 | - | - | To be tested |
| Performance | 3 | - | - | To be tested |
| PWA | 2 | - | - | To be tested |
| Security | 2 | - | - | To be tested |
| **Total** | **28** | **0** | **0** | **Pending** |

---

## Known Issues & TODOs

### High Priority
- ❌ Override functionality not implemented (requires admin permissions)
- ❌ Manual ticket entry not implemented
- ❌ Session statistics not fetching from backend (hardcoded to 0)
- ❌ Offline queue for scans not implemented

### Medium Priority
- ❌ Pull-to-refresh not implemented
- ❌ Session end functionality not exposed in UI
- ❌ Filter/search in session history
- ❌ Export session data

### Low Priority
- ❌ Dark mode support
- ❌ Multiple language support
- ❌ Sound effects for scans
- ❌ Scan history within session

---

## Next Steps

1. ✅ **Configure Backend URL** - Set correct backend API URL
2. ✅ **Start Backend** - Ensure backend is running with seed data
3. ✅ **Run Frontend** - Start development server
4. ✅ **Execute Test Plan** - Go through all test scenarios
5. ✅ **Document Results** - Update test results table
6. ✅ **Fix Issues** - Address any failures
7. ✅ **Retest** - Verify fixes
8. ✅ **Production Build** - Test production build
9. ✅ **Deploy** - Deploy to production environment

---

**Testing Status**: Ready for comprehensive E2E testing  
**Last Updated**: February 3, 2026  
**Tester**: Pending assignment
