#!/usr/bin/env sh
set -eu

compose() {
	SECRET_KEY="$SECRET_KEY" PORT="$PORT" HOST_PORT="$HOST_PORT" docker compose "$@"
}

fail() {
	echo "E2E failed: $*" >&2
	compose logs --no-color >&2 || true
	exit 1
}

status() {
	url=$1
	output=$2
	shift 2
	curl -sS -o "$output" -w '%{http_code}' "$@" "$url"
}

SECRET_KEY="${SECRET_KEY:-$(python3 -c 'import secrets; print(secrets.token_urlsafe(48))' 2>/dev/null || true)}"
[ "${#SECRET_KEY}" -ge 32 ] || fail "SECRET_KEY must be at least 32 bytes (python3 is required to generate one automatically)"
PORT="${PORT:-8000}"
HOST_PORT="${HOST_PORT:-18000}"
BASE_URL="http://127.0.0.1:${HOST_PORT}"

compose_config="$(compose config)" || fail "docker compose config failed"
printf '%s\n' "$compose_config" | grep -Eq "published: \"?${HOST_PORT}\"?" || fail "Compose did not publish HOST_PORT=${HOST_PORT}"
printf '%s\n' "$compose_config" | grep -Eq "target: ${PORT}$" || fail "Compose did not target PORT=${PORT}"
tmpdir="$(mktemp -d)"
email="e2e-$(date +%s)-$$@example.test"
username="e2e$(date +%s)$$"

cleanup() {
	status=$?
	compose down --volumes --remove-orphans >/dev/null 2>&1 || true
	rm -rf "$tmpdir"
	exit "$status"
}
trap cleanup EXIT INT TERM

command -v curl >/dev/null 2>&1 || fail "curl is required"

echo "Starting Compose services on ${BASE_URL}..."
compose up --build -d --wait || fail "docker compose up failed"

for _ in $(seq 1 30); do
	if [ "$(curl -sS -o /dev/null -w '%{http_code}' "${BASE_URL}/healthz" || true)" = "200" ]; then
		break
	fi
	sleep 1
done
[ "$(curl -sS -o /dev/null -w '%{http_code}' "${BASE_URL}/healthz" || true)" = "200" ] || fail "healthz did not become ready"

signup_status="$(status "${BASE_URL}/api/auth/signup" "$tmpdir/signup.json" -H 'Content-Type: application/json' -d "{\"username\":\"${username}\",\"email\":\"${email}\",\"password\":\"correct-horse-battery-staple\"}")"
[ "$signup_status" = "201" ] || fail "signup returned HTTP ${signup_status}: $(cat "$tmpdir/signup.json")"

login_status="$(status "${BASE_URL}/api/auth/login" "$tmpdir/login.json" -H 'Content-Type: application/json' -d "{\"email\":\"${email}\",\"password\":\"correct-horse-battery-staple\"}")"
[ "$login_status" = "200" ] || fail "login returned HTTP ${login_status}: $(cat "$tmpdir/login.json")"

access_token="$(python3 - "$tmpdir/login.json" <<'PY'
import json, sys
value = json.load(open(sys.argv[1])).get("access_token")
if not isinstance(value, str) or not value:
    raise SystemExit("missing access_token in login response")
print(value)
PY
)" || fail "could not parse access token"
refresh_token="$(python3 - "$tmpdir/login.json" <<'PY'
import json, sys
value = json.load(open(sys.argv[1])).get("refresh_token")
if not isinstance(value, str) or not value:
    raise SystemExit("missing refresh_token in login response")
print(value)
PY
)" || fail "could not parse refresh token"

me_status="$(status "${BASE_URL}/api/user/me" "$tmpdir/me.json" -H "Authorization: Bearer ${access_token}")"
[ "$me_status" = "200" ] || fail "access token /api/user/me returned HTTP ${me_status}: $(cat "$tmpdir/me.json")"

refresh_status="$(status "${BASE_URL}/api/user/me" "$tmpdir/refresh-me.json" -H "Authorization: Bearer ${refresh_token}")"
[ "$refresh_status" = "401" ] || fail "refresh token /api/user/me returned HTTP ${refresh_status}: $(cat "$tmpdir/refresh-me.json")"

echo "E2E passed."
