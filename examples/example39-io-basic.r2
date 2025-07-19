// Example 39: Basic IO Operations
// This example demonstrates basic file I/O operations using r2lang

func main() {
    std.print("=== Example 39: Basic IO Operations ===");

    // 1. Basic File Writing
    let content = "Hello, R2Lang!\nThis is a test file.\nLine 3 with special chars: áéíóú";
    io.writeFile("test_basic.txt", content);
    std.print("✓ File written: test_basic.txt");

    // 2. Basic File Reading
    let readContent = io.readFile("test_basic.txt");
    std.print("✓ File content read:");
    std.print(readContent);

    // 3. File Existence Check
    let fileExists = io.exists("test_basic.txt");
    std.print("✓ File exists:", fileExists);

    // 4. File Size
    let size = io.fileSize("test_basic.txt");
    std.print("✓ File size:", size, "bytes");

    // 5. File Information
    let info = io.getMetadata("test_basic.txt");
    std.print("✓ File metadata:");
    std.print("  Name:", info["name"]);
    std.print("  Size:", info["size"]);
    std.print("  Is Directory:", info["isDir"]);
    std.print("  Extension:", info["ext"]);

    // 6. Append to File
    io.appendFile("test_basic.txt", "\nAppended line");
    let newContent = io.readFile("test_basic.txt");
    std.print("✓ After append:");
    std.print(newContent);

    // 7. Copy File
    io.copyFile("test_basic.txt", "test_copy.txt");
    std.print("✓ File copied to test_copy.txt");

    // 8. Rename File
    io.renameFile("test_copy.txt", "test_renamed.txt");
    std.print("✓ File renamed to test_renamed.txt");

    // 9. Cleanup
    io.rmFile("test_basic.txt");
    io.rmFile("test_renamed.txt");
    std.print("✓ Files cleaned up");

    std.print("=== Example 39 completed ===");
}