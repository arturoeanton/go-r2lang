// Example 42: Advanced OS Operations
// This example demonstrates advanced OS operations using r2lang

func main() {
    std.print("=== Example 42: Advanced OS Operations ===");

    // 1. Command Execution with Timeout
    std.print("✓ Command with Timeout:");
    let timeoutResult = os.execWithTimeout("echo 'Quick command'", 5);
    std.print("  Timeout command success:", timeoutResult["success"]);
    std.print("  Timeout command output:", timeoutResult["output"]);

    // 2. Command with Environment
    std.print("✓ Command with Custom Environment:");
    let envMap = {
        "CUSTOM_VAR": "CustomValue",
        "ANOTHER_VAR": "AnotherValue"
    };
    let envResult = os.execWithEnv("echo $CUSTOM_VAR", envMap);
    std.print("  Environment command success:", envResult["success"]);
    std.print("  Environment command output:", envResult["output"]);

    // 3. System Load and Memory (Linux/Unix only)
    std.print("✓ System Monitoring:");
    if (os.getPlatform() == "linux") {
        let loadAvg = os.getLoadAvg();
        if (loadAvg["error"] == null) {
            std.print("  Load Average 1min:", loadAvg["1min"]);
            std.print("  Load Average 5min:", loadAvg["5min"]);
            std.print("  Load Average 15min:", loadAvg["15min"]);
        } else {
            std.print("  Load Average:", loadAvg["error"]);
        }

        let memInfo = os.getMemoryInfo();
        if (memInfo["error"] == null) {
            std.print("  Memory Total:", memInfo["MemTotal"]);
            std.print("  Memory Free:", memInfo["MemFree"]);
            std.print("  Memory Available:", memInfo["MemAvailable"]);
        } else {
            std.print("  Memory Info:", memInfo["error"]);
        }

        let uptime = os.getUptime();
        if (std.typeOf(uptime) == "float64") {
            std.print("  System Uptime:", uptime, "seconds");
        } else {
            std.print("  Uptime:", uptime["error"]);
        }
    } else {
        std.print("  System monitoring features available on Linux only");
    }

    // 4. Disk Usage
    std.print("✓ Disk Usage:");
    let diskUsage = os.getDiskUsage(".");
    if (diskUsage["error"] == null) {
        std.print("  Total space:", diskUsage["total"]);
        std.print("  Free space:", diskUsage["free"]);
        std.print("  Available space:", diskUsage["available"]);
    } else {
        std.print("  Disk usage error:", diskUsage["error"]);
    }

    // 5. Process Information Extended
    std.print("✓ Extended Process Information:");
    let selfPid = os.getPid();
    std.print("  Current process PID:", selfPid);
    let parentPid = os.getParentPid();
    std.print("  Parent process PID:", parentPid);

    // 6. Process Management
    std.print("✓ Process Management:");
    let proc = os.runProcess("sleep 1 && echo 'Background process completed'");
    std.print("  Background process started");
    
    let waitResult = os.waitProcess(proc);
    std.print("  Process wait result:", waitResult);

    // 7. Process with kill demonstration
    std.print("✓ Process Kill Demo:");
    let longProc = os.runProcess("sleep 10");
    std.print("  Started long-running process");
    
    // Wait a moment then kill it
    std.sleep(100); // 100ms
    let killResult = os.killProcess(longProc);
    std.print("  Kill result:", killResult);

    // 8. Working Directory Management
    std.print("✓ Working Directory Management:");
    let originalDir = os.currentDir();
    std.print("  Original directory:", originalDir);
    
    // Change to temp directory and back
    let tempDir = os.getTempDir();
    os.chDir(tempDir);
    let newDir = os.currentDir();
    std.print("  Changed to:", newDir);
    
    os.chDir(originalDir);
    let backDir = os.currentDir();
    std.print("  Changed back to:", backDir);

    // 9. Extended System Information
    std.print("✓ Extended System Information:");
    let sysTime = os.getSystemTime();
    std.print("  System time ISO:", sysTime["iso"]);
    std.print("  System timezone:", sysTime["timezone"]);
    
    let userInfo = os.getUser();
    std.print("  Current user:", userInfo["username"]);
    std.print("  User home:", userInfo["homeDir"]);

    // 10. Environment Variables Management
    std.print("✓ Environment Variables Management:");
    os.setEnv("R2LANG_ADVANCED_TEST", "Advanced Value");
    let advancedVar = os.getEnv("R2LANG_ADVANCED_TEST");
    std.print("  Advanced test variable:", advancedVar);
    
    let allEnvs = os.envList();
    std.print("  Total environment variables:", std.len(std.keys(allEnvs)));

    std.print("=== Example 42 completed ===");
}