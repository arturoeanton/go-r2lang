package r2lang

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
	"strings"
	"time"
)

// r2hack.go: Funciones de "seguridad", "forense" y "análisis" para R2.
// Enfoque didáctico, no pretende ser una suite de hacking real.

func RegisterHack(env *Environment) {

	// 1) hashMD5(str) => string (hex)
	env.Set("hashMD5", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("hashMD5 needs (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("hashMD5: arg must be string")
		}
		sum := md5.Sum([]byte(s))
		return fmt.Sprintf("%x", sum)
	}))

	// 2) hashSHA1(str) => string (hex)
	env.Set("hashSHA1", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("hashSHA1 needs (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("hashSHA1: arg must be string")
		}
		sum := sha1.Sum([]byte(s))
		return fmt.Sprintf("%x", sum)
	}))

	// 3) hashSHA256(str) => string (hex)
	env.Set("hashSHA256", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("hashSHA256 needs (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("hashSHA256: arg must be string")
		}
		sum := sha256.Sum256([]byte(s))
		return fmt.Sprintf("%x", sum)
	}))

	// 4) base64Encode(str) => string
	env.Set("base64Encode", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("base64Encode needs (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("base64Encode: arg must be string")
		}
		return base64.StdEncoding.EncodeToString([]byte(s))
	}))

	// 5) base64Decode(str) => string
	env.Set("base64Decode", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 6) portScan(host, startPort, endPort) => array de puertos abiertos
	// *Muy* simplificado, hace un connect con timeout
	env.Set("portScan", BuiltinFunction(func(args ...interface{}) interface{} {
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
			address := fmt.Sprintf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
			if err == nil {
				// Conexión exitosa => puerto abierto
				conn.Close()
				openPorts = append(openPorts, float64(port))
			}
		}
		return openPorts
	}))

	// 7) whois(domain) => string con respuesta (requiere whois servidor, simplificado)
	// Esto NO es un whois real con servidores especializados, sino un "whois <domain>" shell.
	env.Set("whois", BuiltinFunction(func(args ...interface{}) interface{} {
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
		// Enviar dominio + \r\n
		conn.Write([]byte(domain + "\r\n"))
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
	}))

	// 8) hexdump(str) => string con un volcado hex
	env.Set("hexdump", BuiltinFunction(func(args ...interface{}) interface{} {
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
			sb.WriteString("|\n")
		}
		return sb.String()
	}))

	// 1) hmacSHA256(key, message) => string (hex)
	env.Set("hmacSHA256", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 2) aesEncrypt(key, plaintext) => hex string
	//   - key => 16 bytes (AES-128), 24 bytes (AES-192), o 32 bytes (AES-256)
	//   - generamos iv random
	env.Set("aesEncrypt", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 3) aesDecrypt(key, hexCombined) => plaintext
	env.Set("aesDecrypt", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 4) dnsLookup(host) => array de IPs en string
	env.Set("dnsLookup", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 5) dnsLookupAddr(ip) => array de hosts
	env.Set("dnsLookupAddr", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 6) simplePing(host) => bool indica si hay respuesta
	//   - Con DialTimeout TCP en puerto 80 (no es ICMP real, sino un trick)
	env.Set("simplePing", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// quickRSA(bits) => array [ pubString, privString ]
	env.Set("quickRSA", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 8) rsaEncrypt(pubStr, plaintext) => ciphertext
	//   (demasiado simplificado, parsear pub.N, pub.E)
	env.Set("rsaEncrypt", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	env.Set("rsaDecrypt", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))
}
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
