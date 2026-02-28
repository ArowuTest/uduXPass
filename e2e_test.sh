#!/bin/bash
# uduXPass Full E2E Test Suite
# Tests the complete ticket purchase and validation flow
# Usage: bash e2e_test.sh [BASE_URL]

BASE_URL="${1:-http://localhost:3000}"
TS=$(date +%s)
PASS=0
FAIL=0

check() {
  local name="$1"
  local result="$2"
  local expected="$3"
  if echo "$result" | python3 -c "import sys,json; d=json.load(sys.stdin); assert $expected" 2>/dev/null; then
    echo "  PASS: $name"
    PASS=$((PASS+1))
  else
    echo "  FAIL: $name"
    echo "    Response: $(echo $result | cut -c1-400)"
    FAIL=$((FAIL+1))
  fi
}

echo "================================================================"
echo "uduXPass E2E Test Suite"
echo "Base URL: $BASE_URL"
echo "================================================================"

echo ""
echo "--- Phase 1: Admin & User Authentication ---"

ADMIN_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@uduxpass.com","password":"Admin@123!"}')
ADMIN_TOKEN=$(echo "$ADMIN_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('access_token','') or d.get('data',{}).get('access_token',''))" 2>/dev/null)
check "Admin login" "$ADMIN_RESP" "d.get('success') == True and (d.get('access_token') or d.get('data',{}).get('access_token'))"

# Register new user - NOTE: firstName/lastName are camelCase (not snake_case)
USER_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/auth/email/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"e2e_${TS}@test.com\",\"password\":\"Test@123!\",\"firstName\":\"Test\",\"lastName\":\"User\",\"phone\":\"+234${TS}\"}")
USER_TOKEN=$(echo "$USER_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('access_token','') or d.get('data',{}).get('access_token',''))" 2>/dev/null)
# Registration response does NOT include 'success' field - it directly returns access_token
check "User registration (camelCase fields)" "$USER_RESP" "bool(d.get('access_token'))"

echo ""
echo "--- Phase 2: Event Discovery ---"

EVENTS_RESP=$(curl -s --max-time 10 "$BASE_URL/v1/events")
EVENT_ID=$(echo "$EVENTS_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); events=d.get('data',{}).get('events',[]); print(events[0]['id'] if events else '')" 2>/dev/null)
check "Get events (public)" "$EVENTS_RESP" "d.get('success') == True and len(d.get('data',{}).get('events',[])) > 0"

# Ticket tiers are embedded in event detail, NOT a separate /tiers endpoint
EVENT_RESP=$(curl -s --max-time 10 "$BASE_URL/v1/events/$EVENT_ID")
TIER_ID=$(echo "$EVENT_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); tiers=d.get('data',{}).get('ticket_tiers',[]); print(tiers[0]['id'] if tiers else '')" 2>/dev/null)
check "Get event detail with ticket_tiers" "$EVENT_RESP" "d.get('success') == True and len(d.get('data',{}).get('ticket_tiers',[])) > 0"

echo ""
echo "--- Phase 3: Order & Payment ---"

ORDER_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"event_id\":\"$EVENT_ID\",\"items\":[{\"ticket_tier_id\":\"$TIER_ID\",\"quantity\":2}]}")
ORDER_ID=$(echo "$ORDER_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('order',{}).get('id','') or d.get('order',{}).get('id',''))" 2>/dev/null)
check "Create order" "$ORDER_RESP" "d.get('success') == True"

CONFIRM_RESP=$(curl -s --max-time 15 -X POST "$BASE_URL/v1/admin/orders/$ORDER_ID/confirm-payment" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"payment_method\":\"bank_transfer\",\"payment_reference\":\"E2E_${TS}\",\"notes\":\"E2E test\"}")
check "Admin confirm payment (generates tickets)" "$CONFIRM_RESP" "d.get('success') == True"

echo ""
echo "--- Phase 4: Ticket Retrieval ---"

TICKETS_RESP=$(curl -s --max-time 10 "$BASE_URL/v1/orders/$ORDER_ID/tickets" \
  -H "Authorization: Bearer $USER_TOKEN")
# Tickets endpoint returns data.items (not data.tickets)
TICKET_CODE=$(echo "$TICKETS_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); items=d.get('data',{}).get('items',[]); print(items[0]['qr_code_data'] if items else '')" 2>/dev/null)
check "Get order tickets (JWT QR codes)" "$TICKETS_RESP" "d.get('success') == True and len(d.get('data',{}).get('items',[])) > 0"

echo ""
echo "--- Phase 5: Scanner Operations ---"

SCANNER_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/scanner/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"scanner1","password":"Scanner@123!"}')
SCANNER_TOKEN=$(echo "$SCANNER_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('access_token',''))" 2>/dev/null)
check "Scanner login" "$SCANNER_RESP" "d.get('success') == True and d.get('access_token')"

SESSION_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/scanner/session/start" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SCANNER_TOKEN" \
  -d "{\"event_id\":\"$EVENT_ID\"}")
check "Start scanner session" "$SESSION_RESP" "d.get('success') == True"

# Validate ticket - uses ticket_code and event_id fields
VALIDATE_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/scanner/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SCANNER_TOKEN" \
  -d "{\"ticket_code\":\"$TICKET_CODE\",\"event_id\":\"$EVENT_ID\"}")
check "Validate ticket (first scan)" "$VALIDATE_RESP" "d.get('success') == True and d.get('valid') == True"

# Validate same ticket again - should detect already_validated
REVALIDATE_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/scanner/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SCANNER_TOKEN" \
  -d "{\"ticket_code\":\"$TICKET_CODE\",\"event_id\":\"$EVENT_ID\"}")
check "Validate ticket (duplicate scan detection)" "$REVALIDATE_RESP" "d.get('already_validated') == True"

STATS_RESP=$(curl -s --max-time 10 "$BASE_URL/v1/scanner/stats" \
  -H "Authorization: Bearer $SCANNER_TOKEN")
check "Get scanner stats" "$STATS_RESP" "d.get('success') == True and d.get('data',{}).get('total_scans',0) >= 1"

END_RESP=$(curl -s --max-time 10 -X POST "$BASE_URL/v1/scanner/session/end" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SCANNER_TOKEN" \
  -d '{}')
check "End scanner session" "$END_RESP" "d.get('success') == True"

echo ""
echo "================================================================"
echo "RESULTS: $PASS passed, $FAIL failed out of $((PASS+FAIL)) tests"
if [ $FAIL -eq 0 ]; then
  echo "STATUS: ALL TESTS PASSED"
else
  echo "STATUS: $FAIL TEST(S) FAILED"
fi
echo "================================================================"
