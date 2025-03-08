package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainFunctionality(t *testing.T) {
	// Save original working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"main.go":    "package main\n\nfunc main() {}\n",
		"README.md":  "# Test Project\n",
		".gitignore": "*.log\ntmp/\n",
		"app.log":    "Some log content", // Should be ignored
		"go.mod":     "module testproject\n\ngo 1.17\n",
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Create tmp directory (should be ignored)
	if err := os.Mkdir(filepath.Join(tempDir, "tmp"), 0755); err != nil {
		t.Fatalf("Failed to create tmp directory: %v", err)
	}

	// Change to the test directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Restore original directory when the test completes
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Fatalf("Failed to restore original directory: %v", err)
		}
	}()

	// Set output file name for this test
	testOutputFile := "test-output.md"
	outputFile = testOutputFile

	// Run main (partially - we can't actually call main() in a test)
	rootDir, _ = os.Getwd()
	parseGitignore(".gitignore")

	// Create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Error creating output file: %v", err)
	}

	// Process the project files
	err = processProject(rootDir, file)
	if err != nil {
		t.Fatalf("Error processing project: %v", err)
	}
	file.Close()

	// Verify the output file exists
	if _, err := os.Stat(testOutputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}

	// Read the output file
	content, err := os.ReadFile(testOutputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check content for expected files
	expectedFiles := []string{"main.go", "README.md", "go.mod"}
	for _, expectedFile := range expectedFiles {
		if !contains(string(content), "`"+expectedFile+"`") {
			t.Errorf("Expected output to contain %s, but it didn't", expectedFile)
		}
	}

	// Check that ignored files are not included
	unexpectedFiles := []string{"app.log"}
	for _, unexpectedFile := range unexpectedFiles {
		if contains(string(content), "`"+unexpectedFile+"`") {
			t.Errorf("Expected output to NOT contain %s, but it did", unexpectedFile)
		}
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
