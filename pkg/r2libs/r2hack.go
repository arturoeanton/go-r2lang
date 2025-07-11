package r2libs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2hack.go: Funciones de "seguridad", "forense" y "análisis" para R2.
// Enfoque didáctico, no pretende ser una suite de hacking real.

func RegisterHack(env *r2core.Environment) {
	hackFunctions := []struct {
		Name string
		Func r2core.BuiltinFunction
	}{
		{"hashMD5", hashMD5},
		{"hashSHA1", hashSHA1},
		{"hashSHA256", hashSHA256},
		{"base64Encode", base64Encode},
		{"base64Decode", base64Decode},
		{"portScan", portScan},
		{"whois", whois},
		{"hexdump", hexdump},
		{"hmacSHA256", hmacSHA256},
		{"aesEncrypt", aesEncrypt},
		{"aesDecrypt", aesDecrypt},
		{"dnsLookup", dnsLookup},
		{"dnsLookupAddr", dnsLookupAddr},
		{"simplePing", simplePing},
		{"quickRSA", quickRSA},
		{"rsaEncrypt", rsaEncrypt},
		{"rsaDecrypt", rsaDecrypt},
	}

	for _, f := range hackFunctions {
		env.Set(f.Name, f.Func)
	}
}

var hashMD5 = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("hashMD5 needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("hashMD5: arg must be string")
	}
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", sum)
})

var hashSHA1 = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("hashSHA1 needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("hashSHA1: arg must be string")
	}
	sum := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", sum)
})

var hashSHA256 = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("hashSHA256 needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("hashSHA256: arg must be string")
	}
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
})

var base64Encode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("base64Encode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("base64Encode: arg must be string")
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
})

var base64Decode = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("base64Decode needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("base64Decode: arg must be string")
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Sprintf("base64Decode: error => %v", err)
	}
	return string(decoded)
})

var portScan = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 3 {
		panic("portScan needs (host, startPort, endPort)")
	}
	host, ok1 := args[0].(string)
	startP := int(toFloat(args[1]))
	endP := int(toFloat(args[2]))
	if !ok1 {
		panic("portScan: host debe ser string")
	}
	if startP < 1 || endP > 65535 || endP < startP {
		panic("portScan: invalid port range")
	}
	var openPorts []interface{}
	for port := startP; port <= endP; port++ {
		portStr := strconv.Itoa(port)
		address := net.JoinHostPort(host, portStr)
		conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
		if err == nil {
			// Conexión exitosa => puerto abierto
			conn.Close()
			openPorts = append(openPorts, float64(port))
		}
	}
	return openPorts
})

var whois = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("whois needs (domain)")
	}
	domain, ok := args[0].(string)
	if !ok {
		panic("whois: arg must be string")
	}
	// Podrías llamar un "execCmd" si tienes uno, o un net.Dial("tcp", "whois.server.com:43")
	// Simplificado => net.Dial con whois.verisign-grs.com:43
	conn, err := net.Dial("tcp", "whois.verisign-grs.com:43")
	if err != nil {
		return fmt.Sprintf("Error connecting to whois: %v", err)
	}
	defer conn.Close()
	// Enviar dominio +
	conn.Write([]byte(domain + ""))
	// Leer respuesta
	var sb strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			sb.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	return sb.String()
})

var hexdump = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("hexdump needs (str)")
	}
	s, ok := args[0].(string)
	if !ok {
		panic("hexdump: arg must be a string (or base64 binary)")
	}
	data := []byte(s)
	// Generar volcado
	var sb strings.Builder
	for i := 0; i < len(data); i += 16 {
		line := data[i:]
		if len(line) > 16 {
			line = line[:16]
		}
		// offset
		sb.WriteString(fmt.Sprintf("%08x  ", i))
		// hex
		for j := 0; j < 16; j++ {
			if j < len(line) {
				sb.WriteString(fmt.Sprintf("%02x ", line[j]))
			} else {
				sb.WriteString("   ")
			}
			if j == 7 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(" |")
		// ASCII
		for _, b := range line {
			if b >= 32 && b < 127 {
				sb.WriteByte(b)
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteString("|")
	}
	return sb.String()
})

var hmacSHA256 = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("hmacSHA256 needs (key, message)")
	}
	k, ok1 := args[0].(string)
	m, ok2 := args[1].(string)
	if !(ok1 && ok2) {
		panic("hmacSHA256: args must be strings")
	}
	mac := hmac.New(sha256.New, []byte(k))
	mac.Write([]byte(m))
	return fmt.Sprintf("%x", mac.Sum(nil))
})

var aesEncrypt = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("aesEncrypt needs (key, plaintext)")
	}
	key, ok1 := args[0].(string)
	pt, ok2 := args[1].(string)
	if !(ok1 && ok2) {
		panic("aesEncrypt: args must be strings")
	}
	// crear block AES
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return fmt.Sprintf("aesEncrypt: error => %v", err)
	}
	// generar IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Sprintf("aesEncrypt: error generating IV => %v", err)
	}
	// modo CFB (simplificado), hay muchos modos
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(pt))
	stream.XORKeyStream(ciphertext, []byte(pt))
	// concatenamos IV + ciphertext en hex
	combined := append(iv, ciphertext...)
	return hex.EncodeToString(combined)
})

var aesDecrypt = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("aesDecrypt needs (key, hexCipher)")
	}
	key, ok1 := args[0].(string)
	hexCiph, ok2 := args[1].(string)
	if !(ok1 && ok2) {
		panic("aesDecrypt: args must be strings")
	}
	data, err := hex.DecodeString(hexCiph)
	if err != nil {
		return fmt.Sprintf("aesDecrypt: error decode hex => %v", err)
	}
	if len(data) < aes.BlockSize {
		return "aesDecrypt: data too short"
	}
	iv := data[:aes.BlockSize]
	ciph := data[aes.BlockSize:]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return fmt.Sprintf("aesDecrypt: error => %v", err)
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciph))
	stream.XORKeyStream(plaintext, ciph)
	return string(plaintext)
})

var dnsLookup = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("dnsLookup needs (host)")
	}
	host, ok := args[0].(string)
	if !ok {
		panic("dnsLookup: arg must be string")
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Sprintf("dnsLookup: error => %v", err)
	}
	var arr []interface{}
	for _, ip := range ips {
		arr = append(arr, ip.String())
	}
	return arr
})

var dnsLookupAddr = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("dnsLookupAddr needs (ip)")
	}
	ip, ok := args[0].(string)
	if !ok {
		panic("dnsLookupAddr: arg must be string")
	}
	names, err := net.LookupAddr(ip)
	if err != nil {
		return fmt.Sprintf("dnsLookupAddr: error => %v", err)
	}
	var arr []interface{}
	for _, nm := range names {
		arr = append(arr, nm)
	}
	return arr
})

var simplePing = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("simplePing needs (host)")
	}
	host, ok := args[0].(string)
	if !ok {
		panic("simplePing: arg must be string")
	}
	address := fmt.Sprintf("%s:80", host)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}
	return false
})

var quickRSA = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	var bitSize int = 2048
	if len(args) >= 1 {
		bitSize = int(toFloat(args[0]))
	}
	priv, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return fmt.Sprintf("quickRSA: error => %v", err)
	}
	pub := &priv.PublicKey

	// Armamos un string toy:
	// ej. "RSA-PUB (N=12345, E=65537)"
	pubStr := fmt.Sprintf("RSA-PUB (N=%v, E=%v)", pub.N.String(), pub.E)

	// Incluimos N, E y D en la privada:
	privStr := fmt.Sprintf("RSA-PRIV (N=%v, E=%v, D=%v)",
		pub.N.String(), pub.E, priv.D.String())

	return []interface{}{pubStr, privStr}
})

var rsaEncrypt = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("rsaEncrypt needs (pubStr, plaintext)")
	}
	pubString, ok1 := args[0].(string)
	msg, ok2 := args[1].(string)
	if !(ok1 && ok2) {
		panic("rsaEncrypt => (string, string)")
	}
	// parse "RSA-PUB (N=..., E=...)"
	// muy naive parse
	nVal, eVal := parseRsaPubString(pubString)
	if nVal == nil {
		return "rsaEncrypt: error parse pubKey"
	}
	pubKey := rsa.PublicKey{N: nVal, E: eVal}
	ciph, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pubKey, []byte(msg), nil)
	if err != nil {
		return fmt.Sprintf("rsaEncrypt error => %v", err)
	}
	return hex.EncodeToString(ciph)
})

var rsaDecrypt = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
	// rsaDecrypt(privString, hexCipher) => plaintext
	if len(args) < 2 {
		panic("rsaDecrypt(privString, hexCiph)")
	}
	privString, ok1 := args[0].(string)
	hexCiph, ok2 := args[1].(string)
	if !(ok1 && ok2) {
		panic("rsaDecrypt => (string, string)")
	}

	ciphData, err := hex.DecodeString(hexCiph)
	if err != nil {
		return fmt.Sprintf("rsaDecrypt: decode error => %v", err)
	}

	// parse => "RSA-PRIV (N=..., E=..., D=...)"
	nVal, eVal, dVal := parseRsaPrivString(privString)
	if nVal == nil || dVal == nil {
		return "rsaDecrypt: error parse privKey"
	}

	// Construir PrivateKey "incompleto"
	priv := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: nVal,
			E: eVal,
		},
		D: dVal,
		// Primes: nil => no tendremos la optimización de descifrado,
		// pero con set de Go 1.20+, es posible que Decrypt se queje
		// si no tenemos primes.
		// Haremos un fallback "toy" si no se queja:
	}
	// Precompute => a veces requiere primes
	// priv.Precompute() // Suele fallar si no tienes p y q

	// Decrypt:
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphData, nil)
	if err != nil {
		return fmt.Sprintf("rsaDecrypt error => %v", err)
	}
	return string(plaintext)
})

func parseRsaPrivString(privStr string) (nVal *big.Int, eVal int, dVal *big.Int) {
	// Esperamos "RSA-PRIV (N=..., E=..., D=...)"
	start := strings.Index(privStr, "(")
	end := strings.LastIndex(privStr, ")")
	if start < 0 || end < 0 || end <= start {
		return nil, 0, nil
	}
	inside := privStr[start+1 : end] // "N=..., E=..., D=..."
	// Spliteamos por coma
	parts := strings.Split(inside, ",")
	// Esperamos 3 partes: N=..., E=..., D=...
	if len(parts) != 3 {
		return nil, 0, nil
	}
	var nStr, eStr, dStr string
	for _, p := range parts {
		p = strings.TrimSpace(p) // "N=xxx"
		if strings.HasPrefix(p, "N=") {
			nStr = strings.TrimPrefix(p, "N=")
		} else if strings.HasPrefix(p, "E=") {
			eStr = strings.TrimPrefix(p, "E=")
		} else if strings.HasPrefix(p, "D=") {
			dStr = strings.TrimPrefix(p, "D=")
		}
	}
	// parse nStr => big.Int
	nVal_ := new(big.Int)
	_, ok := nVal_.SetString(strings.TrimSpace(nStr), 10)
	if !ok {
		return nil, 0, nil
	}
	// parse eStr => int
	eVal_ := 0
	fmt.Sscanf(strings.TrimSpace(eStr), "%d", &eVal_)
	// parse dStr => big.Int
	dVal_ := new(big.Int)
	_, ok2 := dVal_.SetString(strings.TrimSpace(dStr), 10)
	if !ok2 {
		return nil, 0, nil
	}
	return nVal_, eVal_, dVal_
}

func parseRsaPubString(pubStr string) (*big.Int, int) {
	// Verificamos que empiece con "RSA-PUB"
	if !strings.HasPrefix(pubStr, "RSA-PUB") {
		return nil, 0
	}
	// Buscamos el paréntesis de apertura "(" y el de cierre ")"
	start := strings.Index(pubStr, "(")
	end := strings.LastIndex(pubStr, ")")
	if start < 0 || end < 0 || end <= start {
		return nil, 0
	}
	// Extraemos la parte interior: "N=..., E=..."
	inside := pubStr[start+1 : end]
	// Dividimos por comas
	parts := strings.Split(inside, ",")
	if len(parts) != 2 {
		return nil, 0
	}
	var nStr, eStr string
	for _, p := range parts {
		p = strings.TrimSpace(p) // "N=xxx"
		if strings.HasPrefix(p, "N=") {
			nStr = strings.TrimPrefix(p, "N=")
		} else if strings.HasPrefix(p, "E=") {
			eStr = strings.TrimPrefix(p, "E=")
		}
	}
	// Convertir nStr a *big.Int
	nVal := new(big.Int)
	_, ok := nVal.SetString(strings.TrimSpace(nStr), 10)
	if !ok {
		return nil, 0
	}
	// Convertir eStr a int
	eVal := 0
	fmt.Sscanf(strings.TrimSpace(eStr), "%d", &eVal)
	return nVal, eVal
}
