package r2libs

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterEncoding(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"base64Encode":    encodingBase64Encode,
		"base64Decode":    encodingBase64Decode,
		"base64UrlEncode": encodingBase64UrlEncode,
		"base64UrlDecode": encodingBase64UrlDecode,
		"hexEncode":       encodingHexEncode,
		"hexDecode":       encodingHexDecode,
		"urlEncode":       encodingUrlEncode,
		"urlDecode":       encodingUrlDecode,
		"urlParse":        encodingUrlParse,
	}

	RegisterModule(env, "encoding", functions)

	RegisterModule(env, "uuid", map[string]r2core.BuiltinFunction{
		"v4":      uuidV4,
		"isValid": uuidIsValid,
	})
}

var encodingBase64Encode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.base64Encode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.base64Encode: arg must be string")
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
})

var encodingBase64Decode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.base64Decode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.base64Decode: arg must be string")
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("encoding.base64Decode: invalid base64 input: %v", err))
	}
	return string(data)
})

var encodingBase64UrlEncode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.base64UrlEncode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.base64UrlEncode: arg must be string")
	}
	return base64.URLEncoding.EncodeToString([]byte(s))
})

var encodingBase64UrlDecode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.base64UrlDecode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.base64UrlDecode: arg must be string")
	}
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("encoding.base64UrlDecode: invalid base64url input: %v", err))
	}
	return string(data)
})

var encodingHexEncode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.hexEncode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.hexEncode: arg must be string")
	}
	return hex.EncodeToString([]byte(s))
})

var encodingHexDecode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.hexDecode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.hexDecode: arg must be string")
	}
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("encoding.hexDecode: invalid hex input: %v", err))
	}
	return string(data)
})

var encodingUrlEncode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.urlEncode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.urlEncode: arg must be string")
	}
	return url.QueryEscape(s)
})

var encodingUrlDecode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.urlDecode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.urlDecode: arg must be string")
	}
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		panic(fmt.Sprintf("encoding.urlDecode: invalid encoded input: %v", err))
	}
	return decoded
})

var encodingUrlParse = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("encoding.urlParse needs (urlString)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("encoding.urlParse: arg must be string")
	}
	parsed, err := url.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("encoding.urlParse: invalid url: %v", err))
	}

	query := make(map[string]interface{})
	for key, values := range parsed.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	result := make(map[string]interface{})
	result["scheme"] = parsed.Scheme
	result["host"] = parsed.Host
	result["path"] = parsed.Path
	result["query"] = query
	return result
})

var uuidV4 = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	buf := make([]byte, 16)
	// crypto/rand, not math/rand: UUIDs must be unpredictable, not just distinct.
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("uuid.v4: failed to generate random bytes: %v", err))
	}
	buf[6] = (buf[6] & 0x0f) | 0x40 // version 4
	buf[8] = (buf[8] & 0x3f) | 0x80 // RFC 4122 variant

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
})

var uuidRegexp = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

var uuidIsValid = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("uuid.isValid needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		return false
	}
	return uuidRegexp.MatchString(s)
})
