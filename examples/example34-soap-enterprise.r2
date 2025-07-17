// example34-soap-enterprise.r2: Comprehensive enterprise SOAP client examples
// Demonstrates customizable headers, response parsing, HTTPS/SSL, and authentication

print("🏢 === ENTERPRISE SOAP CLIENT EXAMPLES ===");
print("");

// Example 1: Basic SOAP Client with Custom Headers
print("1️⃣ Basic Client with Custom Headers");
try {
    // Create client with custom headers
    let customHeaders = {};
    customHeaders["User-Agent"] = "MyEnterprise-SOAPClient/2.0";
    customHeaders["X-Company"] = "My Corporation";
    customHeaders["X-Version"] = "1.0";

    
    let client = soapClient("http://www.dneonline.com/calculator.asmx?WSDL", customHeaders);
    print("   ✅ Client created with custom headers");
    
    // Check current headers
    let headers = client.getHeaders();
    print("   📋 Current headers:", Object.keys(headers).length, "headers set");
    
    // Simple operation call
    let result = client.callSimple("Add", {"intA": 100, "intB": 200});
    print("   🧮 Add(100,200) =", result);
    
} catch (error) {
    print("   ❌ Error:", error);
}

print("");

// Example 2: Response Parsing - Full vs Simple
print("2️⃣ Response Parsing Examples");
try {
    let client = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Full response with all metadata
    print("   🔍 Full response parsing:");
    let fullResponse = client.call("Multiply", {"intA": 15, "intB": 4});
    if (typeOf(fullResponse) == "map") {
        print("     - Success:", fullResponse.success);
        print("     - Result:", fullResponse.result);
        print("     - Values:", fullResponse.values);
    } else {
        print("     - Direct result:", fullResponse);
    }
    
    // Simplified response (just the result)
    print("   🎯 Simple response:");
    let simpleResult = client.callSimple("Divide", {"intA": 100, "intB": 4});
    print("     - Divide(100,4) =", simpleResult);
    
    // Raw XML response
    print("   📄 Raw XML response:");
    let rawResponse = client.callRaw("Subtract", {"intA": 50, "intB": 25});
    print("     - Raw response length:", len(rawResponse), "characters");
    
} catch (error) {
    print("   ❌ Error:", error);
}

print("");

// Example 3: HTTPS/SSL Configuration
print("3️⃣ HTTPS/SSL Configuration");
try {
    // Note: Using HTTP example, but shows how HTTPS would work
    let secureClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Configure TLS settings for enterprise security
    let tlsConfig = {}
    tlsConfig["minVersion"] = "1.2";  // Require TLS 1.2 minimum
    tlsConfig["skipVerify"] = false;   // Always verify certificates in production
    secureClient.setTLSConfig(tlsConfig);
    
    print("   🔒 TLS configuration set (TLS 1.2 minimum)");
    
    // For testing with self-signed certificates (NOT for production)
    let testClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    let tlsConfig = {};
    tlsConfig["skipVerify"] = true;  // ⚠️ Only for testing!
    tlsConfig["minVersion"] = "1.2";  // Still require TLS
    testClient.setTLSConfig(tlsConfig);
    
    print("   ⚠️  Test client configured (skip verify for self-signed certs)");
    
    let result = secureClient.callSimple("Add", {"intA": 25, "intB": 75});
    print("   ✅ Secure call successful:", result);
    
} catch (error) {
    print("   ❌ Error:", error);
}

print("");

// Example 4: Authentication Examples
print("4️⃣ Authentication Configuration");
try {
    let authClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Basic Authentication (username/password)
    print("   🔐 Basic Authentication setup:");
    let authConfig = {};
    authConfig["type"] = "basic";
    authConfig["username"] = "enterprise_user";
    authConfig["password"] = "secure_password_123";
    authClient.setAuth(authConfig);
    print("     - Basic auth configured");
    
    // Bearer Token Authentication (for OAuth/JWT)
    print("   🎫 Bearer Token setup:");
    let tokenClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    let tokenConfig = {};
    tokenConfig["type"] = "bearer";
    tokenConfig["token"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.example.token";  // Example token
    tokenClient.setAuth(tokenConfig);
    print("     - Bearer token configured");
    
    // Note: These won't work with the test service, but show configuration
    print("     ✅ Authentication methods demonstrated");
    
} catch (error) {
    print("   ❌ Error:", error);
}

print("");

// Example 5: Enterprise Error Handling and Resilience
print("5️⃣ Enterprise Error Handling");

func handleSOAPCall(client, operation, params, description) {
    print("   🔄 Attempting:", description);
    try {
        let result = client.callSimple(operation, params);
        print("     ✅ Success:", result);
        return result;
    } catch (error) {
        let errorStr = "" + error;
        
        if (indexOf(errorStr, "timeout") != -1) {
            print("     ⏱️  Timeout - increase timeout or check network");
            return null;
        }
        if (indexOf(errorStr, "connection") != -1) {
            print("     🌐 Connection error - check connectivity");
            return null;
        }
        if (indexOf(errorStr, "certificate") != -1) {
            print("     🔒 Certificate error - check TLS configuration");
            return null;
        }
        if (indexOf(errorStr, "authentication") != -1) {
            print("     🔐 Authentication error - check credentials");
            return null;
        } 
        print("     ❌ Other error:", error);
        
        return null;
    }
}

try {
    let resilientClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Configure for enterprise use
    resilientClient.setTimeout(60.0);  // 60 second timeout
    resilientClient.setHeader("X-Request-ID", "REQ-" + Math.floor(Math.random() * 10000));
    resilientClient.setHeader("X-Enterprise-Client", "true");
    
    // Test multiple operations with error handling
    handleSOAPCall(resilientClient, "Add", {"intA": 500, "intB": 250}, "Large number addition");
    handleSOAPCall(resilientClient, "Multiply", {"intA": 12, "intB": 12}, "Multiplication test");
    handleSOAPCall(resilientClient, "Divide", {"intA": 144, "intB": 12}, "Division operation");
    
} catch (error) {
    print("   ❌ Client creation failed:", error);
}

print("");

// Example 6: Header Management for Corporate Environments
print("6️⃣ Corporate Header Management");
try {
    let corpClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Set multiple corporate headers at once
    let corpHeaders = {}
    corpHeaders["X-Company-ID"] = "CORP-12345";
    corpHeaders["X-Department"] = "Finance";
    corpHeaders["X-User-Role"] = "Manager";
    corpHeaders["X-Session-ID"] = "SESSION-" + Math.floor(Math.random() * 100000);
    corpHeaders["X-Compliance"] = "SOX-Approved";
    corpClient.setHeaders(corpHeaders);
    
    print("   📋 Corporate headers configured");
    
    // Show current headers
    let allHeaders = corpClient.getHeaders();
    print("   📊 Headers configured:", allHeaders);
    
    // Remove sensitive headers for logging
    corpClient.removeHeader("Authorization");
    print("   🧹 Cleaned sensitive headers");
    
    // Reset to defaults if needed
    corpClient.resetHeaders();
    print("   🔄 Headers reset to browser-like defaults");
    
    let result = corpClient.callSimple("Add", {"intA": 999, "intB": 1});
    print("   ✅ Corporate call successful:", result);
    
} catch (error) {
    print("   ❌ Error:", error);
}

print("");

// Example 7: Service Discovery and Metadata
print("7️⃣ Service Discovery");
try {
    let discoveryClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    
    print("   🔍 Service Information:");
    print("     - WSDL URL:", discoveryClient.wsdlURL);
    print("     - Service URL:", discoveryClient.serviceURL);
    print("     - Namespace:", discoveryClient.namespace);
    
    print("   📋 Available Operations:");
    let operations = discoveryClient.listOperations();
    for (let op in operations) {
        print("     -", op);
        
        // Get detailed operation info
        let opInfo = discoveryClient.getOperation(op);
        print("       * SOAP Action:", opInfo.soapAction);
        print("       * Message:", opInfo.message);
    }
    
} catch (error) {
    print("   ❌ Error:", error);
}

print("");

// Summary
print("🎉 === ENTERPRISE SOAP CLIENT CAPABILITIES ===");
print("✅ Customizable HTTP headers (default browser-like)");
print("✅ Response parsing to R2Lang objects and maps");
print("✅ HTTPS/SSL support with configurable TLS versions");  
print("✅ Authentication (Basic Auth, Bearer tokens)");
print("✅ Enterprise error handling and resilience");
print("✅ Corporate header management");
print("✅ Service discovery and metadata extraction");
print("✅ Multiple response formats (full, simple, raw)");
print("");
print("🏢 r2soap is now enterprise-ready for corporate environments!");
print("   Use with internal web services, partner APIs, and secure endpoints.");