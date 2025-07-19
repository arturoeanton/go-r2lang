func main() {
    std.print("=== r2hack Demo ===");

    let text = "Hola Mundo";
    std.print("hashMD5 =>", hack.hashMD5(text));
    std.print("hashSHA1 =>", hack.hashSHA1(text));
    std.print("hashSHA256 =>", hack.hashSHA256(text));

    let b64 = hack.base64Encode(text);
    std.print("base64Encode =>", b64);
    let dec = hack.base64Decode(b64);
    std.print("base64Decode =>", dec);

    // portScan
    let openPorts = hack.portScan("127.0.0.1", 1, 10000);
    std.print("Puertos abiertos =>", openPorts);

    // whois
    let ws = hack.whois("goole.com");
    std.print("whois =>", ws);

    // hexdump
    let hex = hack.hexdump("Test\n\x00\xff!");
    std.print("hexdump =>\n", hex);


     // 1) hmac
        let hm = hack.hmacSHA256("key123", "hola");
        std.print("hmacSHA256 =>", hm);

        // 2) AES
        let cipher = hack.aesEncrypt("1234567890123456", "Mensaje secreto");
        std.print("aesEncrypt =>", cipher);
        let plain = hack.aesDecrypt("1234567890123456", cipher);
        std.print("aesDecrypt =>", plain);

        // 3) DNS
        let ips = hack.dnsLookup("example.com");
        std.print("dnsLookup =>", ips);

        // 4) simplePing
        let alive = hack.simplePing("google.com");
        std.print("simplePing(google.com) =>", alive);

        // 5) RSA
        let keys = hack.quickRSA(1024);
        std.print("quickRSA =>", keys);

        // parse a manual
        let pub = keys[0];
        let priv = keys[1];
        let ciphRSA = hack.rsaEncrypt(pub, "HolaRSA");
        std.print("rsaEncrypt =>", ciphRSA);
        let decRSA = hack.rsaDecrypt(priv, ciphRSA);
        std.print("rsaDecrypt =>", decRSA);

    std.print("=== Fin ===");
}