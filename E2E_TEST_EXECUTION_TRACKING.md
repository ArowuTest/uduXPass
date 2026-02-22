# uduXPass Platform - E2E Test Execution Tracking

**Test Script Version:** 1.0 (Production Certification)  
**Project:** Tems Nigerian Tour / uduXPass Launch  
**Test Date:** February 22, 2026  
**Tester:** Manus AI Agent  
**Approach:** Enterprise-Grade, No Shortcuts, No Assumptions

---

## Test Environment Requirements

### URLs
- [ ] Staging URL: https://staging.uduxpass.com (or local equivalent)
- [ ] Admin Dashboard: /admin
- [ ] Scanner PWA: /scanner

### Payment Sandboxes
- [ ] MoMo PSB Sandbox (API keys active)
- [ ] Paystack Test Mode (API keys active)

### Devices Needed
- [ ] Desktop (Admin Dashboard)
- [ ] Smartphone A (User/Fan experience)
- [ ] Smartphone B (Staff Scanner experience)

---

## Module 1: Admin Command Centre (Setup & Logic)

**Goal:** Verify the organizer can correctly configure the tour.

| ID | Test Action | Expected Result | Status | Notes |
|----|-------------|-----------------|--------|-------|
| 1.1 | Log in to /admin | Access granted to the admin dashboard | ⏳ PENDING | |
| 1.2 | Create a "Tour" and 5 specific "Events" for each city | Tour/Events appear in the database and public listing | ⏳ PENDING | |
| 1.3 | Define Tiers for Lagos: VVIP (₦500k), VIP (₦100k), Regular (₦20k) | Pricing and inventory quantities saved correctly | ⏳ PENDING | |
| 1.4 | Set a "Max 4 per transaction" limit on Regular tickets | Frontend prevents selecting >4 | ⏳ PENDING | |
| 1.5 | Payment Toggle: Enable ONLY MoMo for Abuja event | Abuja checkout shows only MoMo option | ⏳ PENDING | |
| 1.6 | Payment Toggle: Enable BOTH for Lagos event | Lagos checkout shows MoMo and Paystack | ⏳ PENDING | |

---

## Module 2: The Fan Journey - MoMo Flow (Strategic Acquisition)

**Goal:** Validate the seamless account acquisition via MoMo PSB.

| ID | Test Action | Expected Result | Status | Notes |
|----|-------------|-----------------|--------|-------|
| 2.1 | Browse Lagos event on Smartphone A (Unauthenticated) | uduX dark-mode UI renders; event details visible | ⏳ PENDING | |
| 2.2 | Select 2 VIP tickets and click "Get Tickets" | Checkout page shows 10:00 reservation timer | ⏳ PENDING | |
| 2.3 | Select "Pay with MoMo" and enter MoMo Phone Number | System hits MoMo API; UI shows "Awaiting Approval" | ⏳ PENDING | |
| 2.4 | Approve payment in MoMo Sandbox/Simulator | Redirected to "Purchase Successful" page with Order ID | ⏳ PENDING | |
| 2.5 | Check "My Tickets" dashboard (Auto-login via MoMo ID) | User is logged in; 2 tickets visible with valid QR codes | ⏳ PENDING | |

---

## Module 3: The Fan Journey - Paystack Flow (Guest Experience)

**Goal:** Validate email verification and standard payment.

| ID | Test Action | Expected Result | Status | Notes |
|----|-------------|-----------------|--------|-------|
| 3.1 | Select 1 Regular ticket for Lagos event | Proceed to checkout | ⏳ PENDING | |
| 3.2 | Select "Pay with Card/Bank" | System prompts for Email Address | ⏳ PENDING | |
| 3.3 | Enter email and click "Verify Email" | System sends verification link to the inbox | ⏳ PENDING | |
| 3.4 | Click link in email and complete Paystack test payment | Order finalized; "Thank You" page displayed | ⏳ PENDING | |

---

## Module 4: Fulfillment & Communication

**Goal:** Ensure digital assets are delivered correctly.

| ID | Test Action | Expected Result | Status | Notes |
|----|-------------|-----------------|--------|-------|
| 4.1 | Verify Email Inbox for both MoMo and Paystack orders | PDF tickets received; uduX branding present; QR codes clear | ⏳ PENDING | |
| 4.2 | Log in to Dashboard using Email + Password | Dashboard displays purchased tickets correctly | ⏳ PENDING | |
| 4.3 | Download PDF ticket and verify Ticket ID vs Dashboard | IDs match; PDF layout is printer-friendly | ⏳ PENDING | |

---

## Module 5: On-Site Entry (Scanner PWA)

**Goal:** Verify the most critical failure point - the gate.

| ID | Test Action | Expected Result | Status | Notes |
|----|-------------|-----------------|--------|-------|
| 5.1 | Open Scanner URL on Smartphone B; Install as PWA | PWA launches in full-screen (no browser bars) | ⏳ PENDING | |
| 5.2 | Log in with Staff credentials for "Lagos" event | Camera activates; Stats show 0/Total scanned | ⏳ PENDING | |
| 5.3 | Valid Scan: Scan the QR code from Smartphone A | Screen flashes GREEN; Success sound; Vibration | ⏳ PENDING | |
| 5.4 | Duplicate Scan: Scan the same QR code again | Screen flashes RED; Error sound; Shows "ALREADY USED AT [Time]" | ⏳ PENDING | |
| 5.5 | Invalid Scan: Scan a random QR code (non-system) | Screen flashes YELLOW; Error sound; Shows "INVALID TICKET" | ⏳ PENDING | |
| 5.6 | Offline Test: Disable internet on Staff Phone. Scan ticket | Scanner validates against local cache (PWA logic) | ⏳ PENDING | |
| 5.7 | Sync Test: Re-enable internet | Redemptions stored offline sync to Admin Dashboard | ⏳ PENDING | |

---

## Module 6: Security & Data Integrity

**Goal:** Ensure data matches and access is restricted.

| ID | Test Action | Expected Result | Status | Notes |
|----|-------------|-----------------|--------|-------|
| 6.1 | Compare Admin Analytics vs Test Transactions | Revenue totals and inventory counts match perfectly | ⏳ PENDING | |
| 6.2 | Attempt to access /admin via Smartphone A | Access denied; redirection to login | ⏳ PENDING | |
| 6.3 | Export "Lagos" Sales Report to CSV | CSV contains all customer IDs (MoMo/Email) and payment refs | ⏳ PENDING | |

---

## Issues Found

### Critical Issues
(To be populated during testing)

### Major Issues
(To be populated during testing)

### Minor Issues
(To be populated during testing)

---

## Strategic Fixes Required

### Enterprise-Grade Solutions
(To be populated after analysis)

---

## Test Summary

**Total Tests:** 27  
**Passed:** 0  
**Failed:** 0  
**Pending:** 27  
**Pass Rate:** 0%

---

**Status Legend:**
- ⏳ PENDING - Not yet tested
- ✅ PASS - Test passed
- ❌ FAIL - Test failed
- ⚠️ PARTIAL - Partially working
- 🔧 FIXED - Issue found and fixed
