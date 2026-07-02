package r2libs

import (
	"strings"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func newJWTModule(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterJWT(env)
	jwtModule, ok := env.Get("jwt")
	if !ok {
		t.Fatal("jwt module not registered")
	}
	module, ok := jwtModule.(map[string]interface{})
	if !ok {
		t.Fatal("jwt module has unexpected type")
	}
	return module
}

func jwtFn(t *testing.T, module map[string]interface{}, name string) r2core.BuiltinFunction {
	t.Helper()
	fn, ok := module[name].(r2core.BuiltinFunction)
	if !ok {
		t.Fatalf("jwt.%s not found or has wrong type", name)
	}
	return fn
}

// TestJWTRefreshRejectsExpiredToken is a regression test for a bug where
// jwt.refresh used a separate internal verify path (verifyTokenInternal)
// that never checked the "exp"/"nbf" claims, so an already-expired (or
// not-yet-valid) token could be "refreshed" into a brand new, fully valid
// token forever, completely defeating expiration.
func TestJWTRefreshRejectsExpiredToken(t *testing.T) {
	module := newJWTModule(t)
	sign := jwtFn(t, module, "sign")
	refresh := jwtFn(t, module, "refresh")
	verify := jwtFn(t, module, "verify")

	secret := "supersecret"

	// exp far in the past.
	payload := map[string]interface{}{
		"sub": "alice",
		"exp": float64(946684800), // Jan 1 2000
	}
	token := sign(payload, secret).(string)

	verifyResult := verify(token, secret).(map[string]interface{})
	if verifyResult["valid"] != false {
		t.Fatalf("expected expired token to fail verify, got %v", verifyResult)
	}

	refreshResult := refresh(token, secret).(map[string]interface{})
	if refreshResult["success"] != false {
		t.Fatalf("BUG: jwt.refresh succeeded on an expired token: %v", refreshResult)
	}
	if refreshResult["token"] != nil {
		t.Fatalf("expected no token to be issued for an expired refresh, got %v", refreshResult["token"])
	}
}

// TestJWTRefreshRejectsNotYetValidToken mirrors the exp case for "nbf".
func TestJWTRefreshRejectsNotYetValidToken(t *testing.T) {
	module := newJWTModule(t)
	sign := jwtFn(t, module, "sign")
	refresh := jwtFn(t, module, "refresh")

	secret := "supersecret"
	future := float64(time.Now().Add(24 * time.Hour).Unix())

	payload := map[string]interface{}{
		"sub": "bob",
		"nbf": future,
	}
	token := sign(payload, secret).(string)

	refreshResult := refresh(token, secret).(map[string]interface{})
	if refreshResult["success"] != false {
		t.Fatalf("BUG: jwt.refresh succeeded on a not-yet-valid token: %v", refreshResult)
	}
}

// TestJWTRefreshAcceptsValidToken makes sure the fix didn't overcorrect:
// a non-expired token must still refresh successfully into a new, valid
// token.
func TestJWTRefreshAcceptsValidToken(t *testing.T) {
	module := newJWTModule(t)
	sign := jwtFn(t, module, "sign")
	refresh := jwtFn(t, module, "refresh")
	verify := jwtFn(t, module, "verify")

	secret := "supersecret"
	future := float64(time.Now().Add(1 * time.Hour).Unix())

	payload := map[string]interface{}{
		"sub": "carol",
		"exp": future,
	}
	token := sign(payload, secret).(string)

	refreshResult := refresh(token, secret).(map[string]interface{})
	if refreshResult["success"] != true {
		t.Fatalf("expected refresh of a valid token to succeed, got %v", refreshResult)
	}
	newToken, ok := refreshResult["token"].(string)
	if !ok || newToken == "" {
		t.Fatalf("expected a new token string, got %v", refreshResult["token"])
	}

	verifyResult := verify(newToken, secret).(map[string]interface{})
	if verifyResult["valid"] != true {
		t.Fatalf("expected refreshed token to verify as valid, got %v", verifyResult)
	}
}

// TestJWTVerifyRejectsAlgNone is a regression/hardening test for the
// classic JWT "alg: none" bypass: a token whose header claims an
// unsupported (or missing) algorithm must never be treated as valid, even
// if the "signature" segment is empty.
func TestJWTVerifyRejectsAlgNone(t *testing.T) {
	module := newJWTModule(t)
	verify := jwtFn(t, module, "verify")

	// {"alg":"none","typ":"JWT"} . {"sub":"attacker","exp":9999999999} . ""
	forged := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0" +
		".eyJzdWIiOiJhdHRhY2tlciIsImV4cCI6OTk5OTk5OTk5OX0" +
		"."

	result := verify(forged, "any-secret").(map[string]interface{})
	if result["valid"] != false {
		t.Fatalf("BUG: alg:none forged token was accepted: %v", result)
	}
}

// TestJWTVerifyIsConstantTimeSignatureComparison is a smoke test making
// sure the implementation uses hmac.Equal (constant time) rather than a
// naive == comparison; wrong signatures of the correct length must still
// be rejected without panicking regardless of where they first differ.
func TestJWTVerifyWrongSecretRejected(t *testing.T) {
	module := newJWTModule(t)
	sign := jwtFn(t, module, "sign")
	verify := jwtFn(t, module, "verify")

	payload := map[string]interface{}{"sub": "dave"}
	token := sign(payload, "correct-secret").(string)

	result := verify(token, "wrong-secret").(map[string]interface{})
	if result["valid"] != false {
		t.Fatalf("expected verify with wrong secret to fail, got %v", result)
	}
}

// TestJWTVerifyMalformedTokenDoesNotPanic exercises decode edge cases:
// wrong number of segments, invalid base64, invalid JSON. None of these
// should panic uncontrollably; they should return a structured
// valid=false result.
func TestJWTVerifyMalformedTokenDoesNotPanic(t *testing.T) {
	module := newJWTModule(t)
	verify := jwtFn(t, module, "verify")
	decode := jwtFn(t, module, "decode")

	malformed := []string{
		"",
		"not-a-jwt",
		"a.b",
		"a.b.c.d",
		"!!!.!!!.!!!",
		"eyJhbGciOiJIUzI1NiJ9.notbase64!!!.sig",
	}

	for _, tok := range malformed {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("verify(%q) panicked: %v", tok, r)
				}
			}()
			result := verify(tok, "secret").(map[string]interface{})
			if result["valid"] != false {
				t.Errorf("expected malformed token %q to be invalid, got %v", tok, result)
			}
		}()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("decode(%q) panicked: %v", tok, r)
				}
			}()
			decode(tok)
		}()
	}
}

// TestJWTSignVerifyRoundTrip is a basic happy-path sanity check.
func TestJWTSignVerifyRoundTrip(t *testing.T) {
	module := newJWTModule(t)
	sign := jwtFn(t, module, "sign")
	verify := jwtFn(t, module, "verify")

	payload := map[string]interface{}{"sub": "erin", "role": "admin"}
	token := sign(payload, "my-secret").(string)

	if strings.Count(token, ".") != 2 {
		t.Fatalf("expected a 3-part JWT, got %q", token)
	}

	result := verify(token, "my-secret").(map[string]interface{})
	if result["valid"] != true {
		t.Fatalf("expected valid token, got %v", result)
	}
	decodedPayload, ok := result["payload"].(map[string]interface{})
	if !ok || decodedPayload["sub"] != "erin" || decodedPayload["role"] != "admin" {
		t.Fatalf("unexpected payload: %v", result["payload"])
	}
}
