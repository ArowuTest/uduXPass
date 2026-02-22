# Entity Structure Audit - uduXPass Backend

## Critical Findings

### 1. Ticket Entity
**Location:** `internal/domain/entities/ticket.go`

**Actual Structure:**
```go
type Ticket struct {
    ID              uuid.UUID     // NOT an array, single UUID
    OrderLineID     uuid.UUID
    SerialNumber    string
    QRCodeData      string
    QRCodeImageURL  *string
    Status          TicketStatus
    RedeemedAt      *time.Time
    RedeemedBy      *string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**Key Issues:**
- ❌ NO `TierName` field
- ❌ NO `Price` field  
- ✅ `ID` is `uuid.UUID` (NOT array)
- ✅ Has `OrderLineID` to get tier info

**Solution:** Must join with OrderLine → TicketTier to get tier name and price

---

### 2. Order Entity
**Location:** `internal/domain/entities/order.go`

**Actual Structure:**
```go
type Order struct {
    ID       uuid.UUID
    EventID  string      // ⚠️ STRING, not UUID!
    // ... other fields
}
```

**Key Issue:**
- ⚠️ `EventID` is **string**, not `uuid.UUID`
- This explains the type mismatch error

---

### 3. OrderLine Entity
**Location:** `internal/domain/entities/order_line.go`

**Actual Structure:**
```go
type OrderLine struct {
    ID           uuid.UUID
    OrderID      uuid.UUID
    TicketTierID uuid.UUID
    Quantity     int
    UnitPrice    float64
    Subtotal     float64
    
    // Relations
    TicketTier   *TicketTier  // Can access tier info
}
```

**Key Insight:**
- OrderLine has `TicketTier` relation
- Can get tier name and price through this

---

### 4. TicketTier Entity
**Location:** `internal/domain/entities/ticket_tier.go`

**Actual Structure:**
```go
type TicketTier struct {
    ID          uuid.UUID
    EventID     uuid.UUID
    Name        string      // ✅ This is the tier name
    Price       float64     // ✅ This is the price
    Description *string
    Quota       int
    Sold        int
    // ... other fields
}
```

**Key Insight:**
- `Name` field contains tier name (e.g., "VIP", "Regular")
- `Price` field contains ticket price

---

## Data Flow for PDF Generation

To generate PDF tickets with tier name and price:

```
Ticket → OrderLine → TicketTier
  ↓         ↓            ↓
  ID    TicketTierID   Name, Price
```

**Correct Query Pattern:**
1. Get Ticket by ID
2. Get OrderLine by Ticket.OrderLineID
3. Get TicketTier by OrderLine.TicketTierID
4. Use TicketTier.Name and TicketTier.Price in PDF

---

## Fixes Required

### Fix #1: PDF Generation Service
**File:** `internal/infrastructure/email/send_ticket_pdf.go`

**Current (Wrong):**
```go
ticket.ID         // Treated as array
ticket.TierName   // Doesn't exist
ticket.Price      // Doesn't exist
```

**Should Be:**
```go
ticket.ID.String()           // Convert UUID to string
orderLine.TicketTier.Name    // Get from joined relation
orderLine.TicketTier.Price   // Get from joined relation
```

### Fix #2: Payment Service
**File:** `internal/usecases/payments/payment_service.go`

**Current (Wrong):**
```go
order.EventID  // Passed as uuid.UUID
```

**Should Be:**
```go
order.EventID  // Already a string, use directly
// OR parse to UUID if function requires it
eventUUID, _ := uuid.Parse(order.EventID)
```

---

## Action Plan

1. ✅ **Phase 1:** Entity audit complete
2. ⏭️ **Phase 2:** Fix PDF generation to use OrderLine relations
3. ⏭️ **Phase 3:** Fix payment service EventID type handling
4. ⏭️ **Phase 4:** Fix PasswordService interface issues
5. ⏭️ **Phase 5:** Compile and verify all fixes
6. ⏭️ **Phase 6:** Commit to GitHub
7. ⏭️ **Phase 7:** Deploy and test

---

**Champion Developer Note:**  
This audit reveals the root cause - we made assumptions about entity structure without reading the actual code. The fix requires proper database joins to access related data, not direct field access.
