package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsBinaryFile(t *testing.T) {
	// Create a temporary text file
	textFile := filepath.Join(t.TempDir(), "sample.txt")
	err := os.WriteFile(textFile, []byte("This is a text file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test text file: %v", err)
	}

	// Create a temporary binary file
	binFile := filepath.Join(t.TempDir(), "sample.bin")
	data := []byte{0, 1, 2, 3, 0, 5}
	err = os.WriteFile(binFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test binary file: %v", err)
	}

	// Create a temporary image file (just by extension)
	imgFile := filepath.Join(t.TempDir(), "sample.png")
	err = os.WriteFile(imgFile, []byte("fake png content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	// Test cases
	testCases := []struct {
		path     string
		expected bool
	}{
		{textFile, false},
		{binFile, true},
		{imgFile, true}, // Should be true based on extension
	}

	for _, tc := range testCases {
		if result := isBinaryFile(tc.path); result != tc.expected {
			t.Errorf("isBinaryFile(%s) = %v, expected %v", tc.path, result, tc.expected)
		}
	}
}

func TestDetermineCodeFence(t *testing.T) {
	testCases := []struct {
		path     string
		expected string
	}{
		{"file.go", "go"},
		{"script.js", "javascript"},
		{"style.css", "css"},
		{"document.md", "markdown"}, // Changed from "text" to "markdown"
		{"unknown.xyz", "text"},     // Changed from "" to "text"
	}

	for _, tc := range testCases {
		if result := determineCodeFence(tc.path); result != tc.expected {
			t.Errorf("determineCodeFence(%s) = %s, expected %s", tc.path, result, tc.expected)
		}
	}
}

func TestProcessProject(t *testing.T) {
	// Create a temporary project structure
	tempDir := t.TempDir()

	// Create some test files (using OS-specific path joins for creating the files)
	files := map[string]string{
		"main.go":                           "package main\n\nfunc main() {}\n",
		filepath.Join("utils", "helper.go"): "package utils\n\nfunc Helper() {}\n", // OS-specific path join
		"README.md":                         "# Test Project\n",
		"docs.md":                           "# Documentation\n\n```go\nfunc example() {}\n```\n",
		filepath.Join("build", "binary"):    string([]byte{0, 1, 2, 3}), // Binary file
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		// Create directory if needed
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Write file
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Create output file
	outputPath := filepath.Join(tempDir, "output.md")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Reset global variables
	ignoredPaths = make(map[string]bool)
	ignoredPatterns = []string{"README.md"} // Ignore README.md but process other .md files
	ignoredDirs = []string{".git", "build"}

	// Save the old rootDir and restore it after the test
	oldRootDir := rootDir
	rootDir = tempDir
	defer func() { rootDir = oldRootDir }()

	// Process the project
	err = processProject(tempDir, outputFile)
	if err != nil {
		t.Fatalf("processProject failed: %v", err)
	}

	// Close the file to ensure all data is written
	outputFile.Close()

	// Read the output file
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(output)

	// Check if expected files are in the output (always using forward slashes for comparison)
	expectedFiles := []string{
		"main.go",
		"utils/helper.go", // Always use forward slashes here since we normalize in the output
		"docs.md",         // Should be included and handled properly
	}

	for _, file := range expectedFiles {
		if !strings.Contains(outputStr, "`"+file+"`") {
			t.Errorf("Expected output to contain %s, but it didn't", file)
			t.Logf("Output content: %s", outputStr)
		}
	}

	// Specifically check markdown handling
	if strings.Contains(outputStr, "docs.md") {
		if !strings.Contains(outputStr, "(Markdown file)") {
			t.Errorf("Expected docs.md to be marked as Markdown file")
		}

		// Check for tildes instead of HTML comment markers in the new approach
		if !strings.Contains(outputStr, "~~~go") {
			t.Errorf("Expected docs.md to use tilde code blocks instead of backtick code blocks")
		}
	}

	// Check if ignored files are NOT in the output
	unexpectedFiles := []string{"build/binary", "README.md"} // Always use forward slashes here
	for _, file := range unexpectedFiles {
		if strings.Contains(outputStr, "`"+file+"`") {
			t.Errorf("Expected output to NOT contain %s, but it did", file)
		}
	}
}
