// Example 41: Basic OS Operations
// This example demonstrates basic operating system operations using r2lang


func main() {
    std.print("=== Example 41: Basic OS Operations ===");

    // 1. System Information
    std.print("✓ System Information:");
    std.print("  Platform:", os.getPlatform());
    std.print("  Architecture:", os.getArch());
    std.print("  CPU Count:", os.getNumCPU());
    std.print("  Hostname:", os.getHostname());

    // 2. Directory Operations
    let currentDir = os.currentDir();
    std.print("✓ Current Directory:", currentDir);

    let tempDir = os.getTempDir();
    std.print("✓ Temp Directory:", tempDir);

    let homeDir = os.getHomeDir();
    std.print("✓ Home Directory:", homeDir);

    // 3. Environment Variables
    std.print("✓ Environment Variables:");
    let pathVar = os.getEnv("PATH");
    std.print("  PATH length:", std.len(pathVar));
    
    // Set and get custom environment variable
    os.setEnv("R2LANG_TEST", "Hello from R2Lang");
    let customVar = os.getEnv("R2LANG_TEST");
    std.print("  Custom variable:", customVar);

    // 4. Environment listing (showing first few)
    let allEnvs = os.envList();
    let envKeys = std.keys(allEnvs);
    std.print("✓ Total environment variables:", std.len(envKeys));
    std.print("  First 5 env vars:");
    let i = 0;
    while (i < 5 && i < std.len(envKeys)) {
        let key = envKeys[i];
        std.print("    " + key + " = " + string.substr(allEnvs[key], 0, 50) + "...");
        i = i + 1;
    }

    // 5. Process Information
    std.print("✓ Process Information:");
    std.print("  Process ID:", os.getPid());
    std.print("  Parent Process ID:", os.getParentPid());

    // 6. User Information
    let userInfo = os.getUser();
    std.print("✓ User Information:");
    std.print("  Username:", userInfo["username"]);
    std.print("  Name:", userInfo["name"]);
    std.print("  UID:", userInfo["uid"]);
    std.print("  Home Dir:", userInfo["homeDir"]);

    // 7. System Time
    let sysTime = os.getSystemTime();
    std.print("✓ System Time:");
    std.print("  Unix timestamp:", sysTime["unix"]);
    std.print("  ISO format:", sysTime["iso"]);
    std.print("  Local time:", sysTime["local"]);
    std.print("  UTC time:", sysTime["utc"]);
    std.print("  Timezone:", sysTime["timezone"]);

    // 8. Simple Command Execution
    let result = os.execCmd("echo 'Hello from shell command'");
    std.print("✓ Command execution result:");
    std.print(result);

    // 9. Directory Listing
    let files = os.listDir(".");
    std.print("✓ Files in current directory:", std.len(files));
    let fileCount = 0;
    while (fileCount < 5 && fileCount < std.len(files)) {
        std.print("  " + files[fileCount]);
        fileCount = fileCount + 1;
    }

    std.print("=== Example 41 completed ===");
}