// example34-soap-enterprise.r2: Comprehensive enterprise SOAP client examples
// Demonstrates customizable headers, response parsing, HTTPS/SSL, and authentication

std.print("🏢 === ENTERPRISE SOAP CLIENT EXAMPLES ===");
std.print("");

// Example 1: Basic SOAP Client with Custom Headers
std.print("1️⃣ Basic Client with Custom Headers");
try {
    // Create client with custom headers
    let customHeaders = {};
    customHeaders["User-Agent"] = "MyEnterprise-SOAP.client/2.0";
    customHeaders["X-Company"] = "My Corporation";
    customHeaders["X-Version"] = "1.0";

    
    let client = soap.client("http://www.dneonline.com/calculator.asmx?WSDL", customHeaders);
    std.print("   ✅ Client created with custom headers");
    
    // Check current headers
    let headers = client.getHeaders();
    let keys = std.keys(headers);
    let length = std.len(keys); 
    std.print("   📋 Current headers:", length, "headers set");
    
    // Simple operation call
    let result = client.callSimple("Add", {"intA": 100, "intB": 200});
    std.print("   🧮 Add(100,200) =", result);
    
} catch (error) {
    std.print("   ❌ Error:", error);
}

std.print("");

// Example 2: Response Parsing - Full vs Simple
std.print("2️⃣ Response Parsing Examples");
try {
    let client = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Full response with all metadata
    std.print("   🔍 Full response parsing:");
    let fullResponse = client.call("Multiply", {"intA": 15, "intB": 4});
    if (std.typeOf(fullResponse) == "map") {
        std.print("     - Success:", fullResponse.success);
        std.print("     - Result:", fullResponse.result);
        std.print("     - Values:", fullResponse.values);
    } else {
        std.print("     - Direct result:", fullResponse);
    }
    
    // Simplified response (just the result)
    std.print("   🎯 Simple response:");
    let simpleResult = client.callSimple("Divide", {"intA": 100, "intB": 4});
    std.print("     - Divide(100,4) =", simpleResult);
    
    // Raw XML response
    std.print("   📄 Raw XML response:");
    let rawResponse = client.callRaw("Subtract", {"intA": 50, "intB": 25});
    std.print("     - Raw response length:", std.len(rawResponse), "characters");
    
} catch (error) {
    std.print("   ❌ Error:", error);
}

std.print("");

// Example 3: HTTPS/SSL Configuration
std.print("3️⃣ HTTPS/SSL Configuration");
try {
    // Note: Using HTTP example, but shows how HTTPS would work
    let secureClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Configure TLS settings for enterprise security
    let tlsConfig = {}
    tlsConfig["minVersion"] = "1.2";  // Require TLS 1.2 minimum
    tlsConfig["skipVerify"] = false;   // Always verify certificates in production
    secureClient.setTLSConfig(tlsConfig);
    
    std.print("   🔒 TLS configuration set (TLS 1.2 minimum)");
    
    // For testing with self-signed certificates (NOT for production)
    let testClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    let tlsConfig = {};
    tlsConfig["skipVerify"] = true;  // ⚠️ Only for testing!
    tlsConfig["minVersion"] = "1.2";  // Still require TLS
    testClient.setTLSConfig(tlsConfig);
    
    std.print("   ⚠️  Test client configured (skip verify for self-signed certs)");
    
    let result = secureClient.callSimple("Add", {"intA": 25, "intB": 75});
    std.print("   ✅ Secure call successful:", result);
    
} catch (error) {
    std.print("   ❌ Error:", error);
}

std.print("");

// Example 4: Authentication Examples
std.print("4️⃣ Authentication Configuration");
try {
    let authClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Basic Authentication (username/password)
    std.print("   🔐 Basic Authentication setup:");
    let authConfig = {};
    authConfig["type"] = "basic";
    authConfig["username"] = "enterprise_user";
    authConfig["password"] = "secure_password_123";
    authClient.setAuth(authConfig);
    std.print("     - Basic auth configured");
    
    // Bearer Token Authentication (for OAuth/JWT)
    std.print("   🎫 Bearer Token setup:");
    let tokenClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    let tokenConfig = {};
    tokenConfig["type"] = "bearer";
    tokenConfig["token"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.example.token";  // Example token
    tokenClient.setAuth(tokenConfig);
    std.print("     - Bearer token configured");
    
    // Note: These won't work with the test service, but show configuration
    std.print("     ✅ Authentication methods demonstrated");
    
} catch (error) {
    std.print("   ❌ Error:", error);
}

std.print("");

// Example 5: Enterprise Error Handling and Resilience
std.print("5️⃣ Enterprise Error Handling");

func handleSOAPCall(client, operation, params, description) {
    std.print("   🔄 Attempting:", description);
    try {
        let result = client.callSimple(operation, params);
        std.print("     ✅ Success:", result);
        return result;
    } catch (error) {
        let errorStr = "" + error;
        
        if (string.indexOf(errorStr, "timeout") != -1) {
            std.print("     ⏱️  Timeout - increase timeout or check network");
            return null;
        }
        if (string.indexOf(errorStr, "connection") != -1) {
            std.print("     🌐 Connection error - check connectivity");
            return null;
        }
        if (string.indexOf(errorStr, "certificate") != -1) {
            std.print("     🔒 Certificate error - check TLS configuration");
            return null;
        }
        if (string.indexOf(errorStr, "authentication") != -1) {
            std.print("     🔐 Authentication error - check credentials");
            return null;
        } 
        std.print("     ❌ Other error:", error);
        
        return null;
    }
}

try {
    let resilientClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    
    // Configure for enterprise use
    resilientClient.setTimeout(60.0);  // 60 second timeout
    std.print("   ⏳ Timeout set to 60 seconds");
    // Set custom headers for resilience
    resilientClient.setHeader("X-Request-ID", "REQ-" + math.floor(math.random() * 10000));
    resilientClient.setHeader("X-Enterprise-Client", "true");
    std.print("   ✅ Resilient client created with timeout and headers");
    // Test multiple operations with error handling
    handleSOAPCall(resilientClient, "Add", {"intA": 500, "intB": 250}, "Large number addition");
    handleSOAPCall(resilientClient, "Multiply", {"intA": 12, "intB": 12}, "Multiplication test");
    handleSOAPCall(resilientClient, "Divide", {"intA": 144, "intB": 12}, "Division operation");
    
} catch (error) {
    std.print("   ❌ Client creation failed:", error);
}

std.print("");

// Example 6: Header Management for Corporate Environments
std.print("6️⃣ Corporate Header Management");
try {
    let corpClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    

    // Set multiple corporate headers at once
    corpClient.setHeader("X-Company-ID", "CORP-12345");
    corpClient.setHeader("X-Department", "Finance");
    corpClient.setHeader("X-User-Role", "Manager");
    corpClient.setHeader("X-Session-ID", "SESSION-" + math.floor(math.random() * 100000));
    corpClient.setHeader("X-Compliance", "SOX-Approved");
    

    
    std.print("   📋 Corporate headers configured");
    
    // Show current headers
    let allHeaders = corpClient.getHeaders();
    std.print("   📊 Headers configured:", allHeaders);
    
    // Remove sensitive headers for logging
    corpClient.removeHeader("Authorization");
    std.print("   🧹 Cleaned sensitive headers");
    
    // Reset to defaults if needed
    corpClient.resetHeaders();
    std.print("   🔄 Headers reset to browser-like defaults");
    
    let result = corpClient.callSimple("Add", {"intA": 999, "intB": 1});
    std.print("   ✅ Corporate call successful:", result);
    
} catch (error) {
    std.print("   ❌ Error:", error);
}

std.print("");

// Example 7: Service Discovery and Metadata
std.print("7️⃣ Service Discovery");
try {
    let discoveryClient = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    
    std.print("   🔍 Service Information:");
    std.print("     - WSDL URL:", discoveryClient.wsdlURL);
    std.print("     - Service URL:", discoveryClient.serviceURL);
    std.print("     - Namespace:", discoveryClient.namespace);
    
    std.print("   📋 Available Operations:");
    let operations = discoveryClient.listOperations();
    for (i in operations) {
        op = operations[i];
        std.print("     -", op);
        
        // Get detailed operation info
        let opInfo = discoveryClient.getOperation(op);
        std.print("       * SOAP Action:", opInfo.soapAction);
        std.print("       * Message:", opInfo.message);
    }
    
} catch (error) {
    std.print("   ❌ Error:", error);
}

std.print("");

// Summary
std.print("🎉 === ENTERPRISE SOAP CLIENT CAPABILITIES ===");
std.print("✅ Customizable HTTP headers (default browser-like)");
std.print("✅ Response parsing to R2Lang objects and maps");
std.print("✅ HTTPS/SSL support with configurable TLS versions");  
std.print("✅ Authentication (Basic Auth, Bearer tokens)");
std.print("✅ Enterprise error handling and resilience");
std.print("✅ Corporate header management");
std.print("✅ Service discovery and metadata extraction");
std.print("✅ Multiple response formats (full, simple, raw)");
std.print("");
std.print("🏢 r2soap is now enterprise-ready for corporate environments!");
std.print("   Use with internal web services, partner APIs, and secure endpoints.");