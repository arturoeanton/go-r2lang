// Mocking and Spying Example for R2Lang
// This demonstrates Phase 2 features: mocking, spying, and fixtures

import "r2test" as test;
import "r2test/mocking" as mock;
import "r2test/fixtures" as fixtures;

// Example service that we want to test
let EmailService = {
    apiUrl: "https://api.email.com",
    
    sendEmail: func(to, subject, body) {
        // This would normally make an HTTP request
        return httpPost(this.apiUrl + "/send", {
            to: to,
            subject: subject,
            body: body
        });
    },
    
    validateEmail: func(email) {
        return email.indexOf("@") > 0 && email.indexOf(".") > 0;
    }
};

// Example database service
let DatabaseService = {
    users: [],
    
    saveUser: func(user) {
        this.users.push(user);
        return user.id;
    },
    
    getUserById: func(id) {
        for (let i = 0; i < len(this.users); i++) {
            if (this.users[i].id == id) {
                return this.users[i];
            }
        }
        return null;
    },
    
    getAllUsers: func() {
        return this.users;
    }
};

// User management service that depends on other services
let UserService = {
    create: func(userData) {
        if (!EmailService.validateEmail(userData.email)) {
            throw "Invalid email address";
        }
        
        let user = {
            id: generateId(),
            name: userData.name,
            email: userData.email,
            createdAt: new Date()
        };
        
        let userId = DatabaseService.saveUser(user);
        
        EmailService.sendEmail(
            user.email,
            "Welcome!",
            "Welcome to our platform, " + user.name + "!"
        );
        
        return userId;
    }
};

// Mock testing examples
test.describe("Mocking Examples", func() {
    
    test.describe("Basic Mocking", func() {
        test.it("should mock HTTP requests", func() {
            let httpMock = mock.createMock("httpPost");
            
            // Set up expectation
            httpMock.when("httpPost", "https://api.email.com/send", {
                to: "test@example.com",
                subject: "Test",
                body: "Test message"
            }).returns({status: 200, message: "Email sent"});
            
            // Call the method under test
            let result = EmailService.sendEmail("test@example.com", "Test", "Test message");
            
            // Verify the result
            test.assert("email result").that(result.status).equals(200);
            test.assert("email message").that(result.message).equals("Email sent");
            
            // Verify the mock was called
            test.assert("mock was called").that(httpMock.wasCalled("httpPost")).isTrue();
        });
        
        test.it("should handle mock errors", func() {
            let httpMock = mock.createMock("httpPost");
            
            httpMock.when("httpPost").returnsError("Network error");
            
            test.assert("network error").throws(func() {
                EmailService.sendEmail("test@example.com", "Test", "Test");
            }).withMessage("Network error");
        });
    });
    
    test.describe("Spy Examples", func() {
        test.it("should spy on method calls", func() {
            let emailSpy = mock.spyOn("EmailService.validateEmail", EmailService.validateEmail);
            emailSpy.callThrough();
            
            // Call the method
            let isValid = EmailService.validateEmail("test@example.com");
            
            // Verify the result and spy
            test.assert("email validation").that(isValid).isTrue();
            test.assert("spy was called").that(emailSpy.wasCalled()).isTrue();
            test.assert("spy called with").that(emailSpy.wasCalledWith("test@example.com")).isTrue();
        });
        
        test.it("should stub method behavior", func() {
            let dbSpy = mock.spyOn("DatabaseService.saveUser", DatabaseService.saveUser);
            dbSpy.and().returnValue(123);
            
            let result = DatabaseService.saveUser({id: 1, name: "Test"});
            
            test.assert("stubbed return").that(result).equals(123);
            test.assert("spy called").that(dbSpy.wasCalledTimes(1)).isTrue();
        });
    });
    
    test.describe("Mock Verification", func() {
        test.it("should verify call counts", func() {
            let dbMock = mock.createMock("database");
            
            dbMock.when("saveUser").returns(42).times(2);
            
            // Call twice
            dbMock.call("saveUser", {name: "User1"});
            dbMock.call("saveUser", {name: "User2"});
            
            // Verification should pass
            test.assert("mock verification").doesNotThrow(func() {
                dbMock.verify();
            });
        });
        
        test.it("should fail verification for wrong call count", func() {
            let dbMock = mock.createMock("database");
            
            dbMock.when("saveUser").returns(42).times(2);
            
            // Call only once
            dbMock.call("saveUser", {name: "User1"});
            
            // Verification should fail
            test.assert("mock verification failure").throws(func() {
                dbMock.verify();
            });
        });
    });
});

// Fixture testing examples
test.describe("Fixture Examples", func() {
    
    test.beforeAll(func() {
        fixtures.setBasePath("./test_fixtures");
    });
    
    test.it("should load JSON fixtures", func() {
        // Create a temporary JSON fixture
        let userData = {
            id: 1,
            name: "John Doe",
            email: "john@example.com",
            age: 30
        };
        
        fixtures.createTemporary("user_data.json", userData);
        
        // Load and verify
        let loadedData = fixtures.getData("user_data.json");
        
        test.assert("fixture name").that(loadedData.name).equals("John Doe");
        test.assert("fixture email").that(loadedData.email).equals("john@example.com");
        test.assert("fixture age").that(loadedData.age).equals(30);
    });
    
    test.it("should load text fixtures", func() {
        let emailTemplate = `
Hello {{name}},

Welcome to our platform! We're excited to have you.

Best regards,
The Team
        `;
        
        fixtures.createTemporary("email_template.txt", emailTemplate);
        
        let template = fixtures.getData("email_template.txt");
        
        test.assert("template contains").that(template).contains("Welcome to our platform");
        test.assert("template has placeholder").that(template).contains("{{name}}");
    });
    
    test.it("should load CSV fixtures", func() {
        let csvData = [
            {name: "Alice", age: "25", city: "New York"},
            {name: "Bob", age: "30", city: "Los Angeles"},
            {name: "Charlie", age: "35", city: "Chicago"}
        ];
        
        fixtures.createTemporary("users.csv", csvData);
        
        let users = fixtures.getData("users.csv");
        
        test.assert("csv length").that(len(users)).equals(3);
        test.assert("first user name").that(users[0].name).equals("Alice");
        test.assert("second user city").that(users[1].city).equals("Los Angeles");
    });
});

// Integration testing with mocks
test.describe("Integration Testing with Mocks", func() {
    
    test.beforeEach(func() {
        // Reset all mocks before each test
        mock.resetAllMocks();
        mock.resetAllSpies();
    });
    
    test.afterEach(func() {
        // Restore all mocks after each test
        mock.restoreAllMocks();
        mock.restoreAllSpies();
    });
    
    test.it("should create user with mocked dependencies", func() {
        // Mock the generateId function
        let idMock = mock.createMock("generateId");
        idMock.when("generateId").returns(12345);
        
        // Mock the database
        let dbMock = mock.createMock("DatabaseService.saveUser");
        dbMock.when("saveUser").returns(12345);
        
        // Mock the email service
        let emailMock = mock.createMock("EmailService.sendEmail");
        emailMock.when("sendEmail").returns({status: 200});
        
        // Create user
        let userId = UserService.create({
            name: "Test User",
            email: "test@example.com"
        });
        
        // Verify the result
        test.assert("user created").that(userId).equals(12345);
        
        // Verify all dependencies were called
        test.assert("id generated").that(idMock.wasCalled("generateId")).isTrue();
        test.assert("user saved").that(dbMock.wasCalled("saveUser")).isTrue();
        test.assert("email sent").that(emailMock.wasCalled("sendEmail")).isTrue();
        
        // Verify email was sent with correct parameters
        test.assert("email recipient").that(emailMock.wasCalledWith(
            "sendEmail",
            "test@example.com",
            "Welcome!",
            "Welcome to our platform, Test User!"
        )).isTrue();
    });
    
    test.it("should handle email validation failure", func() {
        test.assert("invalid email").throws(func() {
            UserService.create({
                name: "Test User",
                email: "invalid-email"
            });
        }).withMessage("Invalid email address");
    });
});

// Test isolation examples
test.describe("Test Isolation", func() {
    
    test.it("should run in isolated context", func() {
        mock.runInIsolation("isolated test", func(context) {
            // Create mocks within the isolation context
            let mock1 = context.createMock("service1");
            let mock2 = context.createMock("service2");
            
            mock1.when("method1").returns("result1");
            mock2.when("method2").returns("result2");
            
            // Use the mocks
            let result1 = mock1.call("method1");
            let result2 = mock2.call("method2");
            
            test.assert("isolated mock 1").that(result1).equals("result1");
            test.assert("isolated mock 2").that(result2).equals("result2");
            
            // Set global state within context
            context.setGlobalVariable("testVar", "testValue");
            let value = context.getGlobalVariable("testVar");
            
            test.assert("isolated global").that(value).equals("testValue");
        });
    });
});

// Run the tests
print("Running Mocking and Fixtures Tests...");
let results = test.runTests();

print("\n=== TEST RESULTS ===");
print("Total tests:", results.total);
print("Passed:", results.passed);
print("Failed:", results.failed);
print("Skipped:", results.skipped);

if (results.failed > 0) {
    print("\nFailed tests:");
    for (let i = 0; i < len(results.failures); i++) {
        let failure = results.failures[i];
        print("  -", failure.test + ":", failure.message);
    }
}