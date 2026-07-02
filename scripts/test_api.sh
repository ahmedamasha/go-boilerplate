#!/usr/bin/env bash
# End-to-end API test — run with server on :8080
set -euo pipefail
BASE="${BASE_URL:-http://localhost:8080}"
ADMIN_KEY="${ADMIN_API_KEY:-12345}"
PHONE_LOGIN="+966501234567"
PHONE_REGISTER=$(python3 -c "import random; print(f'+9665{random.randint(10000000,99999999)}')")
LOG=/tmp/refda-api-test.log
PASS=0
FAIL=0

assert_ok() {
  local name="$1" body="$2"
  if echo "$body" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if d.get('success') else 1)" 2>/dev/null; then
    echo "  ✓ $name"
    PASS=$((PASS+1))
  else
    echo "  ✗ $name → $body"
    FAIL=$((FAIL+1))
  fi
}

assert_fail_code() {
  local name="$1" body="$2" code="$3"
  if echo "$body" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if d.get('error',{}).get('code')=='$code' else 1)" 2>/dev/null; then
    echo "  ✓ $name (expected $code)"
    PASS=$((PASS+1))
  else
    echo "  ✗ $name → $body"
    FAIL=$((FAIL+1))
  fi
}

# Restart API only if not healthy
if ! curl -sf "$BASE/health" >/dev/null 2>&1; then
  pkill -f 'build/refda-api' 2>/dev/null || true
  lsof -ti :8080 | xargs kill -9 2>/dev/null || true
  sleep 1
  ./build/refda-api > "$LOG" 2>&1 &
  sleep 2
else
  LOG=/tmp/refda-api.log
fi

echo "=== Refda E2E API tests @ $BASE ==="

# ── Health ──────────────────────────────────────────────────────────────────
R=$(curl -s "$BASE/health")
assert_ok "GET /health" "$R"

# ── Auth: login + verify ────────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/auth/login" -H "Content-Type: application/json" -d "{\"phone\":\"$PHONE_LOGIN\"}")
assert_ok "POST /auth/login" "$R"

R=$(curl -s -X POST "$BASE/api/v1/auth/verify" -H "Content-Type: application/json" \
  -d "{\"phone\":\"$PHONE_LOGIN\",\"code\":\"1234\",\"full_name\":\"أحمد محمد\"}")
assert_ok "POST /auth/verify" "$R"
TOKEN=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

# ── Auth: register new user ──────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/auth/register" -H "Content-Type: application/json" \
  -d "{\"phone\":\"$PHONE_REGISTER\",\"full_name\":\"سارة علي\"}")
assert_ok "POST /auth/register" "$R"
R=$(curl -s -X POST "$BASE/api/v1/auth/verify" -H "Content-Type: application/json" \
  -d "{\"phone\":\"$PHONE_REGISTER\",\"code\":\"1234\",\"full_name\":\"سارة علي\"}")
assert_ok "POST /auth/verify (register flow)" "$R"

# ── Auth: duplicate register ──────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/auth/register" -H "Content-Type: application/json" \
  -d "{\"phone\":\"$PHONE_LOGIN\",\"full_name\":\"Test\"}")
assert_fail_code "POST /auth/register (duplicate)" "$R" "USER_EXISTS"

# ── Auth: invalid OTP ────────────────────────────────────────────────────────
curl -s -X POST "$BASE/api/v1/auth/login" -H "Content-Type: application/json" \
  -d "{\"phone\":\"$PHONE_LOGIN\"}" > /dev/null
R=$(curl -s -X POST "$BASE/api/v1/auth/verify" -H "Content-Type: application/json" \
  -d "{\"phone\":\"$PHONE_LOGIN\",\"code\":\"0000\"}")
assert_fail_code "POST /auth/verify (invalid OTP)" "$R" "INVALID_OTP"

# ── Profile ─────────────────────────────────────────────────────────────────
R=$(curl -s "$BASE/api/v1/users/me" -H "Authorization: Bearer $TOKEN")
assert_ok "GET /users/me" "$R"

R=$(curl -s -X PUT "$BASE/api/v1/users/me" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"full_name":"أحمد محمد","location":"الرياض"}')
assert_ok "PUT /users/me (JSON)" "$R"

R=$(curl -s -X PUT "$BASE/api/v1/users/me" -H "Authorization: Bearer $TOKEN" \
  -F "full_name=أحمد محمد" -F "location=جدة")
assert_ok "PUT /users/me (form-data)" "$R"

# ── Create event (pending_review) ──────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/events" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{
  "title": "حفل زفاف أحمد وفاطمة",
  "event_type": "wedding",
  "event_date": "2026-08-15",
  "event_time": "19:00",
  "target_amount": 1000,
  "privacy": "public",
  "gifts": [
    {"name": "طقم ذهب", "type": "product", "target_amount": 800, "image_urls": []},
    {"name": "مساهمة نقدية", "type": "cash"}
  ]
}')
assert_ok "POST /events" "$R"
EVENT_ID=$(echo "$R" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print(d['id'])")
GIFT_PRODUCT=$(echo "$R" | python3 -c "import sys,json; gs=json.load(sys.stdin)['data']['gifts']; print([g['id'] for g in gs if g['type']=='product'][0])")
GIFT_CASH=$(echo "$R" | python3 -c "import sys,json; gs=json.load(sys.stdin)['data']['gifts']; print([g['id'] for g in gs if g['type']=='cash'][0])")
STATUS=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['status'])")
echo "    eventId=$EVENT_ID status=$STATUS"
if [ "$STATUS" != "pending_review" ]; then
  echo "  ✗ POST /events status should be pending_review → $STATUS"
  FAIL=$((FAIL+1))
else
  echo "  ✓ POST /events status=pending_review"
  PASS=$((PASS+1))
fi

# ── Contribute before admin approve (should fail) ───────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/gifts/contribute" -H "Content-Type: application/json" -d "{
  \"event_id\": \"$EVENT_ID\", \"amount\": 100
}")
assert_fail_code "POST /gifts/contribute (before open)" "$R" "EVENT_NOT_OPEN"

# ── Admin approve ────────────────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/admin/events/$EVENT_ID/review" \
  -H "X-Admin-Key: $ADMIN_KEY" -H "Content-Type: application/json" -d '{"approve":true}')
assert_ok "POST /admin/events/:id/review (approve)" "$R"

# ── Events list / get ────────────────────────────────────────────────────────
R=$(curl -s "$BASE/api/v1/events")
assert_ok "GET /events (public)" "$R"

R=$(curl -s "$BASE/api/v1/events?mine=true" -H "Authorization: Bearer $TOKEN")
assert_ok "GET /events?mine=true" "$R"

R=$(curl -s "$BASE/api/v1/events/$EVENT_ID")
assert_ok "GET /events/:id" "$R"

R=$(curl -s "$BASE/api/v1/events/$EVENT_ID?hide_amounts=true")
assert_ok "GET /events/:id?hide_amounts=true" "$R"

# ── Contribute partial product gift ──────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/gifts/contribute" -H "Content-Type: application/json" -d "{
  \"event_id\": \"$EVENT_ID\",
  \"gift_id\": \"$GIFT_PRODUCT\",
  \"amount\": 300,
  \"hide_amount\": false,
  \"message\": \"مبارك عليكم!\",
  \"reserve_full\": false
}")
assert_ok "POST /gifts/contribute (partial product)" "$R"
CONTRIB1=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['contribution_id'])")

R=$(curl -s -X POST "$BASE/api/v1/gifts/pay" -H "Content-Type: application/json" \
  -d "{\"contribution_id\":\"$CONTRIB1\",\"payment_method\":\"mada\"}")
assert_ok "POST /gifts/pay (mada)" "$R"

# ── Cash direct contribution ─────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/gifts/contribute" -H "Content-Type: application/json" -d "{
  \"event_id\": \"$EVENT_ID\",
  \"amount\": 200,
  \"hide_amount\": true,
  \"message\": \"هدية نقدية\"
}")
assert_ok "POST /gifts/contribute (cash direct)" "$R"
CONTRIB2=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['contribution_id'])")

R=$(curl -s -X POST "$BASE/api/v1/gifts/pay" -H "Content-Type: application/json" \
  -d "{\"contribution_id\":\"$CONTRIB2\",\"payment_method\":\"apple_pay\"}")
assert_ok "POST /gifts/pay (apple_pay)" "$R"

# ── Reserve full product gift ────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/gifts/contribute" -H "Content-Type: application/json" -d "{
  \"event_id\": \"$EVENT_ID\",
  \"gift_id\": \"$GIFT_PRODUCT\",
  \"amount\": 500,
  \"reserve_full\": true
}")
assert_ok "POST /gifts/contribute (reserve_full)" "$R"
CONTRIB3=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['contribution_id'])")

R=$(curl -s -X POST "$BASE/api/v1/gifts/pay" -H "Content-Type: application/json" \
  -d "{\"contribution_id\":\"$CONTRIB3\",\"payment_method\":\"mada\"}")
assert_ok "POST /gifts/pay (reserve_full)" "$R"

# ── Event should be collected (1000/1000) ──────────────────────────────────
R=$(curl -s "$BASE/api/v1/events/$EVENT_ID" -H "Authorization: Bearer $TOKEN")
STATUS=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['status'])")
COLLECTED=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['amount_collected'])")
echo "    status=$STATUS collected=$COLLECTED"
if [ "$STATUS" = "collected" ]; then
  echo "  ✓ GET /events/:id status=collected"
  PASS=$((PASS+1))
else
  echo "  ✗ GET /events/:id expected collected → $STATUS"
  FAIL=$((FAIL+1))
fi

# ── Withdraw partial ─────────────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/events/$EVENT_ID/withdraw" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" -d '{"amount":400}')
assert_ok "POST /events/:id/withdraw (partial)" "$R"
WS=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['event_status'])")
if [ "$WS" = "partial_withdrawn" ]; then
  echo "  ✓ withdraw → partial_withdrawn"
  PASS=$((PASS+1))
else
  echo "  ✗ withdraw expected partial_withdrawn → $WS"
  FAIL=$((FAIL+1))
fi

# ── Withdraw remainder ───────────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/events/$EVENT_ID/withdraw" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" -d '{"amount":600}')
assert_ok "POST /events/:id/withdraw (full)" "$R"
WS=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['event_status'])")
if [ "$WS" = "fully_withdrawn" ]; then
  echo "  ✓ withdraw → fully_withdrawn"
  PASS=$((PASS+1))
else
  echo "  ✗ withdraw expected fully_withdrawn → $WS"
  FAIL=$((FAIL+1))
fi

# ── Admin reject (new event) ─────────────────────────────────────────────────
R=$(curl -s -X POST "$BASE/api/v1/events" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"title":"Rejected event","event_type":"other","event_date":"2026-09-01","privacy":"public","gifts":[]}')
REJECT_ID=$(echo "$R" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])")
R=$(curl -s -X POST "$BASE/api/v1/admin/events/$REJECT_ID/review" \
  -H "X-Admin-Key: $ADMIN_KEY" -H "Content-Type: application/json" -d '{"approve":false}')
assert_ok "POST /admin/events/:id/review (reject)" "$R"

echo ""
echo "Results: $PASS passed, $FAIL failed"
exit "$FAIL"
