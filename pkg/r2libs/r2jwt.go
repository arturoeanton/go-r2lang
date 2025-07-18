package r2libs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterJWT(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"sign": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JWT.sign needs (payload, secret)")
			}
			payload, ok1 := args[0].(map[string]interface{})
			secret, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JWT.sign: arguments must be (object, string)")
			}

			algorithm := "HS256"
			if len(args) > 2 {
				if alg, ok := args[2].(string); ok {
					algorithm = alg
				}
			}

			// Create header
			header := map[string]interface{}{
				"alg": algorithm,
				"typ": "JWT",
			}

			// Add timestamps if not present
			now := time.Now().Unix()
			if _, exists := payload["iat"]; !exists {
				payload["iat"] = float64(now)
			}
			if _, exists := payload["nbf"]; !exists {
				payload["nbf"] = float64(now)
			}

			// Encode header and payload
			headerJSON, _ := json.Marshal(header)
			payloadJSON, _ := json.Marshal(payload)

			headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
			payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

			// Create signature
			message := headerEncoded + "." + payloadEncoded
			var signature string

			switch algorithm {
			case "HS256":
				h := hmac.New(sha256.New, []byte(secret))
				h.Write([]byte(message))
				signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
			default:
				panic(fmt.Sprintf("JWT.sign: unsupported algorithm '%s'", algorithm))
			}

			return message + "." + signature
		}),

		"verify": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JWT.verify needs (token, secret)")
			}
			token, ok1 := args[0].(string)
			secret, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JWT.verify: arguments must be (string, string)")
			}

			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Invalid token format",
					"payload": nil,
				}
			}

			headerEncoded := parts[0]
			payloadEncoded := parts[1]
			signatureProvided := parts[2]

			// Decode header
			headerJSON, err := base64.RawURLEncoding.DecodeString(headerEncoded)
			if err != nil {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Invalid header encoding",
					"payload": nil,
				}
			}

			var header map[string]interface{}
			if err := json.Unmarshal(headerJSON, &header); err != nil {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Invalid header JSON",
					"payload": nil,
				}
			}

			// Get algorithm
			algorithm, exists := header["alg"].(string)
			if !exists {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Missing algorithm in header",
					"payload": nil,
				}
			}

			// Verify signature
			message := headerEncoded + "." + payloadEncoded
			var expectedSignature string

			switch algorithm {
			case "HS256":
				h := hmac.New(sha256.New, []byte(secret))
				h.Write([]byte(message))
				expectedSignature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
			default:
				return map[string]interface{}{
					"valid":   false,
					"error":   fmt.Sprintf("Unsupported algorithm '%s'", algorithm),
					"payload": nil,
				}
			}

			if signatureProvided != expectedSignature {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Invalid signature",
					"payload": nil,
				}
			}

			// Decode payload
			payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
			if err != nil {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Invalid payload encoding",
					"payload": nil,
				}
			}

			var payload map[string]interface{}
			if err := json.Unmarshal(payloadJSON, &payload); err != nil {
				return map[string]interface{}{
					"valid":   false,
					"error":   "Invalid payload JSON",
					"payload": nil,
				}
			}

			// Check expiration
			if exp, exists := payload["exp"]; exists {
				if expFloat, ok := exp.(float64); ok {
					if time.Now().Unix() > int64(expFloat) {
						return map[string]interface{}{
							"valid":   false,
							"error":   "Token expired",
							"payload": payload,
						}
					}
				}
			}

			// Check not before
			if nbf, exists := payload["nbf"]; exists {
				if nbfFloat, ok := nbf.(float64); ok {
					if time.Now().Unix() < int64(nbfFloat) {
						return map[string]interface{}{
							"valid":   false,
							"error":   "Token not yet valid",
							"payload": payload,
						}
					}
				}
			}

			return map[string]interface{}{
				"valid":   true,
				"error":   nil,
				"payload": payload,
			}
		}),

		"decode": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JWT.decode needs (token)")
			}
			token, ok := args[0].(string)
			if !ok {
				panic("JWT.decode: token must be string")
			}

			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				return map[string]interface{}{
					"header":  nil,
					"payload": nil,
					"error":   "Invalid token format",
				}
			}

			// Decode header
			headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
			if err != nil {
				return map[string]interface{}{
					"header":  nil,
					"payload": nil,
					"error":   "Invalid header encoding",
				}
			}

			var header map[string]interface{}
			if err := json.Unmarshal(headerJSON, &header); err != nil {
				return map[string]interface{}{
					"header":  nil,
					"payload": nil,
					"error":   "Invalid header JSON",
				}
			}

			// Decode payload
			payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				return map[string]interface{}{
					"header":  header,
					"payload": nil,
					"error":   "Invalid payload encoding",
				}
			}

			var payload map[string]interface{}
			if err := json.Unmarshal(payloadJSON, &payload); err != nil {
				return map[string]interface{}{
					"header":  header,
					"payload": nil,
					"error":   "Invalid payload JSON",
				}
			}

			return map[string]interface{}{
				"header":  header,
				"payload": payload,
				"error":   nil,
			}
		}),

		"createPayload": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JWT.createPayload needs (data)")
			}
			data, ok := args[0].(map[string]interface{})
			if !ok {
				panic("JWT.createPayload: data must be object")
			}

			// Default expiration: 1 hour
			expireInSeconds := 3600
			if len(args) > 1 {
				if exp, ok := args[1].(float64); ok {
					expireInSeconds = int(exp)
				}
			}

			issuer := ""
			if len(args) > 2 {
				if iss, ok := args[2].(string); ok {
					issuer = iss
				}
			}

			subject := ""
			if len(args) > 3 {
				if sub, ok := args[3].(string); ok {
					subject = sub
				}
			}

			audience := ""
			if len(args) > 4 {
				if aud, ok := args[4].(string); ok {
					audience = aud
				}
			}

			now := time.Now().Unix()
			payload := make(map[string]interface{})

			// Copy user data
			for key, value := range data {
				payload[key] = value
			}

			// Standard claims
			payload["iat"] = float64(now)
			payload["nbf"] = float64(now)
			payload["exp"] = float64(now + int64(expireInSeconds))

			if issuer != "" {
				payload["iss"] = issuer
			}
			if subject != "" {
				payload["sub"] = subject
			}
			if audience != "" {
				payload["aud"] = audience
			}

			return payload
		}),

		"isExpired": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JWT.isExpired needs (token)")
			}
			token, ok := args[0].(string)
			if !ok {
				panic("JWT.isExpired: token must be string")
			}

			decoded := decodeTokenInternal(token)
			if decoded == nil {
				return true
			}

			payload, ok := decoded["payload"].(map[string]interface{})
			if !ok {
				return true
			}

			if exp, exists := payload["exp"]; exists {
				if expFloat, ok := exp.(float64); ok {
					return time.Now().Unix() > int64(expFloat)
				}
			}

			return false
		}),

		"getExpiration": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JWT.getExpiration needs (token)")
			}
			token, ok := args[0].(string)
			if !ok {
				panic("JWT.getExpiration: token must be string")
			}

			decoded := decodeTokenInternal(token)
			if decoded == nil {
				return nil
			}

			payload, ok := decoded["payload"].(map[string]interface{})
			if !ok {
				return nil
			}

			if exp, exists := payload["exp"]; exists {
				return exp
			}

			return nil
		}),

		"getClaims": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JWT.getClaims needs (token)")
			}
			token, ok := args[0].(string)
			if !ok {
				panic("JWT.getClaims: token must be string")
			}

			decoded := decodeTokenInternal(token)
			if decoded == nil {
				return nil
			}

			payload, ok := decoded["payload"].(map[string]interface{})
			if !ok {
				return nil
			}

			// Extract standard claims
			claims := make(map[string]interface{})
			standardClaims := []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}

			for _, claim := range standardClaims {
				if value, exists := payload[claim]; exists {
					claims[claim] = value
				}
			}

			return claims
		}),

		"getHeader": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JWT.getHeader needs (token)")
			}
			token, ok := args[0].(string)
			if !ok {
				panic("JWT.getHeader: token must be string")
			}

			decoded := decodeTokenInternal(token)
			if decoded == nil {
				return nil
			}

			return decoded["header"]
		}),

		"refresh": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JWT.refresh needs (token, secret)")
			}
			token, ok1 := args[0].(string)
			secret, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JWT.refresh: arguments must be (string, string)")
			}

			// Verify the token first
			verifyResult := verifyTokenInternal(token, secret)
			verifyMap, ok := verifyResult.(map[string]interface{})
			if !ok {
				return map[string]interface{}{
					"success": false,
					"error":   "Failed to verify token",
					"token":   nil,
				}
			}

			valid, _ := verifyMap["valid"].(bool)
			if !valid {
				return map[string]interface{}{
					"success": false,
					"error":   verifyMap["error"],
					"token":   nil,
				}
			}

			// Get payload and create new token
			payload, ok := verifyMap["payload"].(map[string]interface{})
			if !ok {
				return map[string]interface{}{
					"success": false,
					"error":   "Invalid payload",
					"token":   nil,
				}
			}

			// Remove old timestamps
			delete(payload, "iat")
			delete(payload, "nbf")
			delete(payload, "exp")

			// Create new token with same payload but new timestamps
			newToken := createTokenInternal(payload, secret, "HS256")

			return map[string]interface{}{
				"success": true,
				"error":   nil,
				"token":   newToken,
			}
		}),

		"createRefreshToken": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JWT.createRefreshToken needs (userId, secret)")
			}
			userId, ok1 := args[0].(string)
			secret, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JWT.createRefreshToken: arguments must be (string, string)")
			}

			// Refresh tokens typically have longer expiration (30 days)
			expireInSeconds := 30 * 24 * 3600
			if len(args) > 2 {
				if exp, ok := args[2].(float64); ok {
					expireInSeconds = int(exp)
				}
			}

			payload := map[string]interface{}{
				"sub":  userId,
				"type": "refresh",
			}

			now := time.Now().Unix()
			payload["iat"] = float64(now)
			payload["nbf"] = float64(now)
			payload["exp"] = float64(now + int64(expireInSeconds))

			return createTokenInternal(payload, secret, "HS256")
		}),
	}

	RegisterModule(env, "jwt", functions)
}

// Helper functions
func decodeTokenInternal(token string) map[string]interface{} {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}

	// Decode header
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil
	}

	return map[string]interface{}{
		"header":  header,
		"payload": payload,
	}
}

func verifyTokenInternal(token, secret string) interface{} {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return map[string]interface{}{
			"valid":   false,
			"error":   "Invalid token format",
			"payload": nil,
		}
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureProvided := parts[2]

	// Decode header
	headerJSON, err := base64.RawURLEncoding.DecodeString(headerEncoded)
	if err != nil {
		return map[string]interface{}{
			"valid":   false,
			"error":   "Invalid header encoding",
			"payload": nil,
		}
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return map[string]interface{}{
			"valid":   false,
			"error":   "Invalid header JSON",
			"payload": nil,
		}
	}

	// Verify signature
	message := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if signatureProvided != expectedSignature {
		return map[string]interface{}{
			"valid":   false,
			"error":   "Invalid signature",
			"payload": nil,
		}
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return map[string]interface{}{
			"valid":   false,
			"error":   "Invalid payload encoding",
			"payload": nil,
		}
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return map[string]interface{}{
			"valid":   false,
			"error":   "Invalid payload JSON",
			"payload": nil,
		}
	}

	return map[string]interface{}{
		"valid":   true,
		"error":   nil,
		"payload": payload,
	}
}

func createTokenInternal(payload map[string]interface{}, secret, algorithm string) string {
	// Create header
	header := map[string]interface{}{
		"alg": algorithm,
		"typ": "JWT",
	}

	// Add timestamps if not present
	now := time.Now().Unix()
	if _, exists := payload["iat"]; !exists {
		payload["iat"] = float64(now)
	}
	if _, exists := payload["nbf"]; !exists {
		payload["nbf"] = float64(now)
	}

	// Encode header and payload
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature
	message := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return message + "." + signature
}
