// Coverage and Reporting Example for R2Lang
// This demonstrates Phase 3 features: coverage collection and reporting

import "r2test" as test;
import "r2test/coverage" as coverage;
import "r2test/reporters" as reporters;

// Example module to test for coverage
let StringUtils = {
    // Function 1: String manipulation
    capitalize: func(str) {
        if (str == null || str == "") {
            return "";
        }
        return str.charAt(0).toUpperCase() + str.substring(1).toLowerCase();
    },
    
    // Function 2: String validation
    isValidEmail: func(email) {
        if (email == null || email == "") {
            return false;
        }
        
        let atIndex = email.indexOf("@");
        if (atIndex <= 0) {
            return false;
        }
        
        let dotIndex = email.lastIndexOf(".");
        if (dotIndex <= atIndex) {
            return false;
        }
        
        return true;
    },
    
    // Function 3: String processing with branches
    processText: func(text, options) {
        if (text == null) {
            return null;
        }
        
        let result = text;
        
        if (options && options.trim) {
            result = result.trim();
        }
        
        if (options && options.uppercase) {
            result = result.toUpperCase();
        } else if (options && options.lowercase) {
            result = result.toLowerCase();
        }
        
        if (options && options.removeSpaces) {
            result = result.replace(/\s+/g, "");
        }
        
        return result;
    },
    
    // Function 4: Complex string analysis
    analyzeString: func(str) {
        if (str == null || str == "") {
            return {
                length: 0,
                words: 0,
                characters: 0,
                vowels: 0,
                consonants: 0
            };
        }
        
        let words = str.split(" ").filter(w => w.length > 0);
        let characters = str.replace(/\s/g, "").length;
        let vowels = 0;
        let consonants = 0;
        
        for (let i = 0; i < str.length; i++) {
            let char = str.charAt(i).toLowerCase();
            if (char.match(/[aeiou]/)) {
                vowels++;
            } else if (char.match(/[bcdfghjklmnpqrstvwxyz]/)) {
                consonants++;
            }
        }
        
        return {
            length: str.length,
            words: words.length,
            characters: characters,
            vowels: vowels,
            consonants: consonants
        };
    }
};

// Math utilities for additional coverage testing
let MathUtils = {
    factorial: func(n) {
        if (n < 0) {
            throw "Factorial of negative number is undefined";
        }
        if (n == 0 || n == 1) {
            return 1;
        }
        return n * this.factorial(n - 1);
    },
    
    isPrime: func(n) {
        if (n < 2) {
            return false;
        }
        if (n == 2) {
            return true;
        }
        if (n % 2 == 0) {
            return false;
        }
        
        for (let i = 3; i * i <= n; i += 2) {
            if (n % i == 0) {
                return false;
            }
        }
        return true;
    },
    
    gcd: func(a, b) {
        a = Math.abs(a);
        b = Math.abs(b);
        
        while (b != 0) {
            let temp = b;
            b = a % b;
            a = temp;
        }
        return a;
    }
};

// Configure coverage collection
test.describe("Coverage Configuration", func() {
    
    test.beforeAll(func() {
        // Enable coverage collection
        coverage.enable();
        coverage.start();
        
        // Set coverage base path
        coverage.setBasePath("./");
        
        // Add exclude patterns for files we don't want to track
        coverage.addExcludeGlob("*_test.r2");
        coverage.addExcludeGlob("node_modules/*");
        coverage.addExcludeGlob("coverage/*");
        
        print("Coverage collection enabled");
    });
    
    test.it("should configure coverage settings", func() {
        test.assert("coverage enabled").that(coverage.isEnabled()).isTrue();
    });
});

// Test suite with good coverage
test.describe("StringUtils Coverage Tests", func() {
    
    test.describe("capitalize function", func() {
        test.it("should capitalize normal strings", func() {
            coverage.recordHit("string_utils.r2", 5); // Function entry
            coverage.recordHit("string_utils.r2", 6); // Null check
            coverage.recordHit("string_utils.r2", 9); // Main logic
            
            let result = StringUtils.capitalize("hello world");
            test.assert("capitalized").that(result).equals("Hello world");
        });
        
        test.it("should handle empty strings", func() {
            coverage.recordHit("string_utils.r2", 5); // Function entry
            coverage.recordHit("string_utils.r2", 6); // Null check
            coverage.recordHit("string_utils.r2", 7); // Empty return
            
            let result = StringUtils.capitalize("");
            test.assert("empty string").that(result).equals("");
        });
        
        test.it("should handle null input", func() {
            coverage.recordHit("string_utils.r2", 5); // Function entry
            coverage.recordHit("string_utils.r2", 6); // Null check
            coverage.recordHit("string_utils.r2", 7); // Empty return
            
            let result = StringUtils.capitalize(null);
            test.assert("null input").that(result).equals("");
        });
    });
    
    test.describe("isValidEmail function", func() {
        test.it("should validate correct emails", func() {
            coverage.recordHit("string_utils.r2", 14); // Function entry
            coverage.recordHit("string_utils.r2", 15); // Null check
            coverage.recordHit("string_utils.r2", 19); // @ check
            coverage.recordHit("string_utils.r2", 23); // . check
            coverage.recordHit("string_utils.r2", 27); // Return true
            
            let result = StringUtils.isValidEmail("test@example.com");
            test.assert("valid email").that(result).isTrue();
        });
        
        test.it("should reject emails without @", func() {
            coverage.recordHit("string_utils.r2", 14); // Function entry
            coverage.recordHit("string_utils.r2", 15); // Null check
            coverage.recordHit("string_utils.r2", 19); // @ check
            coverage.recordHit("string_utils.r2", 20); // Return false
            
            let result = StringUtils.isValidEmail("testexample.com");
            test.assert("no @ symbol").that(result).isFalse();
        });
        
        test.it("should reject emails without domain", func() {
            coverage.recordHit("string_utils.r2", 14); // Function entry
            coverage.recordHit("string_utils.r2", 15); // Null check
            coverage.recordHit("string_utils.r2", 19); // @ check
            coverage.recordHit("string_utils.r2", 23); // . check
            coverage.recordHit("string_utils.r2", 24); // Return false
            
            let result = StringUtils.isValidEmail("test@example");
            test.assert("no domain").that(result).isFalse();
        });
    });
    
    test.describe("processText function", func() {
        test.it("should process text with all options", func() {
            coverage.recordHit("string_utils.r2", 32); // Function entry
            coverage.recordHit("string_utils.r2", 36); // Main logic
            coverage.recordHit("string_utils.r2", 38); // Trim check
            coverage.recordHit("string_utils.r2", 39); // Trim action
            coverage.recordHit("string_utils.r2", 42); // Uppercase check
            coverage.recordHit("string_utils.r2", 43); // Uppercase action
            coverage.recordHit("string_utils.r2", 48); // Remove spaces check
            coverage.recordHit("string_utils.r2", 49); // Remove spaces action
            
            let result = StringUtils.processText("  hello world  ", {
                trim: true,
                uppercase: true,
                removeSpaces: true
            });
            test.assert("processed text").that(result).equals("HELLOWORLD");
        });
        
        test.it("should handle lowercase option", func() {
            coverage.recordHit("string_utils.r2", 32); // Function entry
            coverage.recordHit("string_utils.r2", 36); // Main logic
            coverage.recordHit("string_utils.r2", 42); // Uppercase check
            coverage.recordHit("string_utils.r2", 44); // Lowercase check
            coverage.recordHit("string_utils.r2", 45); // Lowercase action
            
            let result = StringUtils.processText("HELLO WORLD", {
                lowercase: true
            });
            test.assert("lowercase text").that(result).equals("hello world");
        });
        
        test.it("should handle null input", func() {
            coverage.recordHit("string_utils.r2", 32); // Function entry
            coverage.recordHit("string_utils.r2", 33); // Null check
            coverage.recordHit("string_utils.r2", 34); // Return null
            
            let result = StringUtils.processText(null);
            test.assert("null input").that(result).isNull();
        });
    });
});

// Test suite with partial coverage (demonstrating uncovered branches)
test.describe("MathUtils Partial Coverage Tests", func() {
    
    test.describe("factorial function", func() {
        test.it("should calculate factorial of positive numbers", func() {
            coverage.recordHit("math_utils.r2", 5);  // Function entry
            coverage.recordHit("math_utils.r2", 9);  // n == 0 check
            coverage.recordHit("math_utils.r2", 12); // Recursive case
            
            let result = MathUtils.factorial(5);
            test.assert("factorial 5").that(result).equals(120);
        });
        
        test.it("should handle factorial of 0", func() {
            coverage.recordHit("math_utils.r2", 5);  // Function entry
            coverage.recordHit("math_utils.r2", 9);  // n == 0 check
            coverage.recordHit("math_utils.r2", 10); // Return 1
            
            let result = MathUtils.factorial(0);
            test.assert("factorial 0").that(result).equals(1);
        });
        
        // Note: We're not testing negative numbers, leaving that branch uncovered
    });
    
    test.describe("isPrime function", func() {
        test.it("should identify prime numbers", func() {
            coverage.recordHit("math_utils.r2", 20); // Function entry
            coverage.recordHit("math_utils.r2", 21); // n < 2 check
            coverage.recordHit("math_utils.r2", 24); // n == 2 check
            coverage.recordHit("math_utils.r2", 30); // Loop check
            coverage.recordHit("math_utils.r2", 35); // Return true
            
            let result = MathUtils.isPrime(17);
            test.assert("17 is prime").that(result).isTrue();
        });
        
        test.it("should identify 2 as prime", func() {
            coverage.recordHit("math_utils.r2", 20); // Function entry
            coverage.recordHit("math_utils.r2", 21); // n < 2 check
            coverage.recordHit("math_utils.r2", 24); // n == 2 check
            coverage.recordHit("math_utils.r2", 25); // Return true
            
            let result = MathUtils.isPrime(2);
            test.assert("2 is prime").that(result).isTrue();
        });
        
        // Note: We're not testing even numbers > 2, leaving that branch uncovered
    });
});

// Coverage reporting tests
test.describe("Coverage Reporting", func() {
    
    test.it("should generate coverage statistics", func() {
        let stats = coverage.getStats();
        
        test.assert("stats object").that(stats).isNotNull();
        test.assert("total files").that(stats.totalFiles).isGreaterThan(0);
        test.assert("line percentage").that(stats.linePercentage).isGreaterThanOrEqual(0);
        test.assert("line percentage max").that(stats.linePercentage).isLessThanOrEqual(100);
        
        print("Coverage Statistics:");
        print("  Total Files:", stats.totalFiles);
        print("  Total Lines:", stats.totalLines);
        print("  Covered Lines:", stats.coveredLines);
        print("  Line Coverage:", stats.linePercentage.toFixed(2) + "%");
        print("  Statement Coverage:", stats.statementPercentage.toFixed(2) + "%");
        print("  Function Coverage:", stats.functionPercentage.toFixed(2) + "%");
    });
    
    test.it("should check coverage thresholds", func() {
        // Test different threshold levels
        let meets50 = coverage.meetsThreshold(50.0);
        let meets80 = coverage.meetsThreshold(80.0);
        let meets95 = coverage.meetsThreshold(95.0);
        
        test.assert("50% threshold").that(meets50).isBoolean();
        test.assert("80% threshold").that(meets80).isBoolean();
        test.assert("95% threshold").that(meets95).isBoolean();
        
        print("Threshold Results:");
        print("  Meets 50%:", meets50);
        print("  Meets 80%:", meets80);
        print("  Meets 95%:", meets95);
    });
    
    test.it("should identify uncovered lines", func() {
        let uncoveredLines = coverage.getUncoveredLines("math_utils.r2");
        
        test.assert("uncovered lines").that(uncoveredLines).isArray();
        
        if (len(uncoveredLines) > 0) {
            print("Uncovered lines in math_utils.r2:", uncoveredLines);
        } else {
            print("All lines covered in math_utils.r2");
        }
    });
});

// Generate reports
test.describe("Report Generation", func() {
    
    test.afterAll(func() {
        // Generate different types of reports
        print("Generating coverage reports...");
        
        // Generate HTML report
        try {
            reporters.generateHTMLReport("./coverage/html", test.getResults());
            print("✓ HTML report generated in ./coverage/html/");
        } catch (e) {
            print("✗ Failed to generate HTML report:", e);
        }
        
        // Generate JSON report
        try {
            reporters.generateJSONReport("./coverage/coverage.json", test.getResults());
            print("✓ JSON report generated: ./coverage/coverage.json");
        } catch (e) {
            print("✗ Failed to generate JSON report:", e);
        }
        
        // Generate JUnit XML report
        try {
            reporters.generateJUnitReport("./coverage/junit.xml", test.getResults());
            print("✓ JUnit report generated: ./coverage/junit.xml");
        } catch (e) {
            print("✗ Failed to generate JUnit report:", e);
        }
        
        // Print final coverage summary
        let finalStats = coverage.getStats();
        print("\n=== FINAL COVERAGE SUMMARY ===");
        print("Line Coverage:", finalStats.linePercentage.toFixed(2) + "%");
        print("Statement Coverage:", finalStats.statementPercentage.toFixed(2) + "%");
        print("Branch Coverage:", finalStats.branchPercentage.toFixed(2) + "%");
        print("Function Coverage:", finalStats.functionPercentage.toFixed(2) + "%");
        
        if (finalStats.linePercentage >= 80) {
            print("🎉 Excellent coverage! Target achieved.");
        } else if (finalStats.linePercentage >= 60) {
            print("👍 Good coverage, but room for improvement.");
        } else {
            print("⚠️  Coverage below recommended threshold.");
        }
    });
});

// Run the tests
print("Running Coverage and Reporting Tests...");
let results = test.runTests();

print("\n=== TEST RESULTS ===");
print("Total tests:", results.total);
print("Passed:", results.passed);
print("Failed:", results.failed);
print("Success rate:", (results.passed / results.total * 100).toFixed(2) + "%");

if (results.failed > 0) {
    print("\nFailed tests:");
    for (let i = 0; i < len(results.failures); i++) {
        let failure = results.failures[i];
        print("  -", failure.test + ":", failure.message);
    }
}