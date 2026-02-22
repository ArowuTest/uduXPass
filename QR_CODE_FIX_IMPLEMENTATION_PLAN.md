# Strategic QR Code Implementation Plan
**Date:** February 13, 2026  
**Developer:** Manus AI (Champion Developer)  
**Priority:** CRITICAL - BLOCKING PRODUCTION

---

## Problem Statement

The uduXPass ticketing platform generates QR code **data strings** but has no mechanism to display them as scannable QR code **images** to end users. This blocks the entire ticketing flow.

---

## Strategic Solution (Multi-Layer Approach)

### Layer 1: Frontend QR Code Display (IMMEDIATE)
**Priority:** 🔴 CRITICAL  
**Effort:** 2-3 hours  
**Impact:** Unblocks user ticket display

**Implementation:**
1. Install production-grade QR code library (`qrcode.react`)
2. Create reusable `TicketQRCode` component
3. Implement ticket display page with QR code
4. Add download/share functionality

### Layer 2: Backend QR Code Image Generation (STRATEGIC)
**Priority:** 🟡 HIGH  
**Effort:** 3-4 hours  
**Impact:** Improves performance, enables offline viewing

**Implementation:**
1. Install Go QR code library (`github.com/skip2/go-qrcode`)
2. Generate QR code images on ticket creation
3. Store images in S3 or file storage
4. Return image URLs in API responses
5. Frontend displays pre-generated images (faster, more reliable)

### Layer 3: Hybrid Approach (PRODUCTION-READY)
**Priority:** 🟢 RECOMMENDED  
**Effort:** 4-5 hours  
**Impact:** Best of both worlds

**Implementation:**
- Backend generates and stores QR images
- Frontend has fallback client-side generation
- API returns both: `qr_code_data` (string) and `qr_code_image_url` (URL)
- Frontend prefers image URL, falls back to client-side generation

---

## Implementation Steps

### Phase 1: Frontend QR Display (Now)

#### Step 1: Install QR Library
```bash
cd /home/ubuntu/frontend
pnpm add qrcode.react
pnpm add @types/qrcode.react --save-dev
```

#### Step 2: Create TicketQRCode Component
**File:** `/home/ubuntu/frontend/src/components/tickets/TicketQRCode.tsx`

Features:
- Display QR code from data string
- Configurable size
- Download as PNG
- Share functionality
- Error handling

#### Step 3: Create Ticket Display Component
**File:** `/home/ubuntu/frontend/src/components/tickets/TicketCard.tsx`

Features:
- Show ticket details (event, tier, serial number)
- Display QR code
- Status indicator (active/redeemed)
- Download/share buttons

#### Step 4: Create User Tickets Page
**File:** `/home/ubuntu/frontend/src/pages/UserTicketsPage.tsx`

Features:
- List all user tickets
- Filter by status (active/redeemed)
- Search by event name
- Individual ticket view

#### Step 5: Update API Service
**File:** `/home/ubuntu/frontend/src/services/api.ts`

Add endpoints:
- `getUserTickets()` - Get all user tickets
- `getTicketById(id)` - Get single ticket
- `downloadTicket(id)` - Download ticket as PDF

---

### Phase 2: Backend QR Image Generation (Strategic)

#### Step 1: Install Go QR Library
```bash
cd /home/ubuntu/backend
go get github.com/skip2/go-qrcode
```

#### Step 2: Create QR Service
**File:** `/home/ubuntu/backend/pkg/qrcode/generator.go`

Functions:
- `GenerateQRCode(data string) ([]byte, error)` - Generate QR image
- `GenerateQRCodeFile(data, filepath string) error` - Save to file
- `GenerateQRCodeBase64(data string) (string, error)` - Base64 encoded

#### Step 3: Update Ticket Creation
**File:** `/home/ubuntu/backend/internal/usecases/payments/payment_service.go`

After generating `qr_code_data`:
1. Generate QR code image
2. Upload to S3 or save to storage
3. Store image URL in database (add `qr_code_image_url` column)

#### Step 4: Update Database Schema
**Migration:** `005_add_qr_image_url.sql`

```sql
ALTER TABLE tickets ADD COLUMN qr_code_image_url VARCHAR(500);
CREATE INDEX idx_tickets_qr_image_url ON tickets(qr_code_image_url);
```

#### Step 5: Update API Responses
Return both:
```json
{
  "qr_code_data": "uduxpass://...",
  "qr_code_image_url": "https://cdn.uduxpass.com/qr/ticket-123.png"
}
```

---

## Testing Strategy

### Unit Tests
- QR code generation (various data formats)
- Component rendering
- Download functionality

### Integration Tests
- Ticket creation → QR generation → Storage
- API endpoints returning correct data
- Frontend displaying QR codes

### E2E Tests
1. User purchases ticket
2. Ticket appears in dashboard with QR code
3. QR code is scannable
4. Scanner validates ticket
5. Ticket cannot be reused

---

## Rollout Plan

### Stage 1: Frontend Fix (Immediate)
- Deploy frontend with client-side QR generation
- Users can immediately see and use tickets
- **Timeline:** Today

### Stage 2: Backend Enhancement (This Week)
- Add backend QR image generation
- Migrate existing tickets to have images
- Update API responses
- **Timeline:** 2-3 days

### Stage 3: Optimization (Next Week)
- Performance testing
- CDN configuration
- Caching strategy
- **Timeline:** 1 week

---

## Success Metrics

### Immediate (Phase 1)
- ✅ Users can view tickets with QR codes
- ✅ QR codes are scannable
- ✅ Download functionality works

### Strategic (Phase 2)
- ✅ QR images load in <200ms
- ✅ 99.9% QR code generation success rate
- ✅ Offline viewing capability
- ✅ Reduced frontend bundle size

---

## Risk Mitigation

### Risk 1: QR Code Not Scannable
**Mitigation:**
- Test with multiple scanner apps
- Validate QR data format
- Add error correction level (High)
- Provide manual entry fallback

### Risk 2: Performance Issues
**Mitigation:**
- Generate QR codes asynchronously
- Use CDN for image delivery
- Implement caching
- Monitor generation time

### Risk 3: Storage Costs
**Mitigation:**
- Compress QR images (PNG with optimization)
- Set expiration policy (delete after event + 30 days)
- Use cost-effective storage tier

---

## Implementation Order

1. ✅ Install frontend QR library
2. ✅ Create TicketQRCode component
3. ✅ Create TicketCard component
4. ✅ Create UserTicketsPage
5. ✅ Update API service
6. ✅ Test frontend QR display
7. ⏳ Install backend QR library
8. ⏳ Create QR service
9. ⏳ Update ticket creation logic
10. ⏳ Add database migration
11. ⏳ Update API responses
12. ⏳ End-to-end testing

---

**Status:** Ready to implement  
**Next Action:** Execute Phase 1 (Frontend QR Display)
