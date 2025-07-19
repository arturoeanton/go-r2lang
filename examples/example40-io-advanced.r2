// Example 40: Advanced IO Operations
// This example demonstrates advanced file I/O operations using r2lang

func main() {
    std.print("=== Example 40: Advanced IO Operations ===");

    // 1. Directory operations
    std.print("✓ Directory operations:");
    io.mkdir("test_dir");
    std.print("  Directory created: test_dir");
    
    let dirContents = io.listDir(".");
    std.print("  Current directory contains", std.len(dirContents), "items");

    // 2. Binary file operations
    std.print("✓ Binary file operations:");
    let binaryData = [72, 101, 108, 108, 111]; // "Hello" in ASCII
    io.writeFileBytes("binary_test.bin", binaryData);
    let readBinary = io.readFileBytes("binary_test.bin");
    std.print("  Binary data written and read:", readBinary);

    // 3. Lines operations
    std.print("✓ Lines operations:");
    let lines = ["First line", "Second line", "Third line"];
    io.writeLines("lines_test.txt", lines);
    let readLinesResult = io.readLines("lines_test.txt");
    std.print("  Lines written and read:", readLinesResult);

    // 4. File information operations
    std.print("✓ File information operations:");
    io.writeFile("info_test.txt", "Test content for file info");
    
    let fileExists = io.exists("info_test.txt");
    std.print("  File exists:", fileExists);
    
    let fileSize = io.fileSize("info_test.txt");
    std.print("  File size:", fileSize, "bytes");
    
    let isFile = io.isFile("info_test.txt");
    std.print("  Is file:", isFile);
    
    let isDir = io.isDir("test_dir");
    std.print("  test_dir is directory:", isDir);

    // 5. File path operations
    std.print("✓ File path operations:");
    let absolutePath = io.absPath("info_test.txt");
    std.print("  Absolute path:", absolutePath);
    
    let baseName = io.baseName("/path/to/file.txt");
    std.print("  Base name of '/path/to/file.txt':", baseName);
    
    let dirName = io.dirName("/path/to/file.txt");
    std.print("  Directory name of '/path/to/file.txt':", dirName);
    
    let extName = io.extName("file.txt");
    std.print("  Extension of 'file.txt':", extName);

    // 6. File operations with permissions
    std.print("✓ File operations with permissions:");
    let tempFile = io.tempFile();
    std.print("  Temporary file created:", tempFile);
    
    let workingDir = io.workingDir();
    std.print("  Working directory:", workingDir);

    // 7. File checksums
    std.print("✓ File checksums:");
    let md5sum = io.checksum("info_test.txt", "md5");
    let sha256sum = io.checksum("info_test.txt", "sha256");
    std.print("  MD5:", md5sum);
    std.print("  SHA256:", sha256sum);

    // 8. Advanced file operations
    std.print("✓ Advanced file operations:");
    io.copyFile("info_test.txt", "copy_test.txt");
    std.print("  File copied to copy_test.txt");
    
    io.moveFile("copy_test.txt", "moved_test.txt");
    std.print("  File moved to moved_test.txt");

    // 9. Batch operations
    std.print("✓ Batch operations:");
    let globResults = io.glob("*.txt");
    std.print("  Found .txt files:", std.len(globResults));
    
    let findResults = io.findFiles(".", "*.txt");
    std.print("  Found files matching *.txt:", std.len(findResults));

    // 10. File metadata
    std.print("✓ File metadata:");
    let metadata = io.getMetadata("info_test.txt");
    std.print("  Metadata name:", metadata["name"]);
    std.print("  Metadata size:", metadata["size"]);
    std.print("  Metadata isDir:", metadata["isDir"]);

    // 11. Cleanup
    std.print("✓ Cleanup:");
    io.rmFile("binary_test.bin");
    io.rmFile("lines_test.txt");
    io.rmFile("info_test.txt");
    io.rmFile("moved_test.txt");
    io.rmFile(tempFile);
    io.rmDir("test_dir");
    std.print("  All test files cleaned up");

    std.print("=== Example 40 completed ===");
}