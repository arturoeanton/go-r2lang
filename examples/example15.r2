func main() {
    println("=== r2hack Demo ===");

    let text = "Hola Mundo";
    println("hashMD5 =>", hashMD5(text));
    println("hashSHA1 =>", hashSHA1(text));
    println("hashSHA256 =>", hashSHA256(text));

    let b64 = base64Encode(text);
    println("base64Encode =>", b64);
    let dec = base64Decode(b64);
    println("base64Decode =>", dec);

    // portScan
    let openPorts = portScan("127.0.0.1", 1, 10000);
    println("Puertos abiertos =>", openPorts);

    // whois
    let ws = whois("goole.com");
    println("whois =>", ws);

    // hexdump
    let hex = hexdump("Test\n\x00\xff!");
    println("hexdump =>\n", hex);


     // 1) hmac
        let hm = hmacSHA256("key123", "hola");
        println("hmacSHA256 =>", hm);

        // 2) AES
        let cipher = aesEncrypt("1234567890123456", "Mensaje secreto");
        println("aesEncrypt =>", cipher);
        let plain = aesDecrypt("1234567890123456", cipher);
        println("aesDecrypt =>", plain);

        // 3) DNS
        let ips = dnsLookup("example.com");
        println("dnsLookup =>", ips);

        // 4) simplePing
        let alive = simplePing("google.com");
        println("simplePing(google.com) =>", alive);

        // 5) RSA
        let keys = quickRSA(1024);
        println("quickRSA =>", keys);

        // parse a manual
        let pub = keys[0];
        let priv = keys[1];
        let ciphRSA = rsaEncrypt(pub, "HolaRSA");
        println("rsaEncrypt =>", ciphRSA);
        let decRSA = rsaDecrypt(priv, ciphRSA);
        println("rsaDecrypt =>", decRSA);

    println("=== Fin ===");
}