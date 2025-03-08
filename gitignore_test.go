package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitignoreParsing(t *testing.T) {
	// Create a temporary .gitignore file
	tmpGitignore := filepath.Join(t.TempDir(), ".gitignore")
	content := `# Comment line
node_modules/
*.log
/dist
build/
!build/important.txt
`
	err := os.WriteFile(tmpGitignore, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .gitignore: %v", err)
	}

	// Reset global variables
	ignoredPaths = make(map[string]bool)
	ignoredPatterns = make([]string, 0)
	ignoredDirs = []string{".git"}

	// Parse the .gitignore file
	parseGitignore(tmpGitignore)

	// Check if patterns were parsed correctly
	expectedDirs := []string{".git", "node_modules", "build"}
	for _, dir := range expectedDirs {
		found := false
		for _, ignoreDir := range ignoredDirs {
			if ignoreDir == dir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected directory %s to be ignored, but it wasn't", dir)
		}
	}

	expectedPatterns := []string{"*.log", "/dist"}
	for _, pattern := range expectedPatterns {
		found := false
		for _, ignorePattern := range ignoredPatterns {
			if ignorePattern == pattern {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected pattern %s to be ignored, but it wasn't", pattern)
		}
	}
}

func TestShouldIgnoreFile(t *testing.T) {
	// Reset and set up test patterns
	ignoredPaths = make(map[string]bool)
	ignoredPatterns = []string{"*.log", "*.tmp", "secret.txt"}

	// Test cases
	testCases := []struct {
		path     string
		expected bool
	}{
		{"app.log", true},
		{"logs/error.log", true},
		{"temp.tmp", true},
		{"secret.txt", true},
		{"main.go", false},
		{"src/utils.js", false},
	}

	for _, tc := range testCases {
		if result := shouldIgnoreFile(tc.path); result != tc.expected {
			t.Errorf("shouldIgnoreFile(%s) = %v, expected %v", tc.path, result, tc.expected)
		}
	}
}

func TestShouldIgnoreDir(t *testing.T) {
	// Set up test directory patterns
	rootDir = "."
	ignoredDirs = []string{".git", "node_modules", "dist"}

	// Test cases
	testCases := []struct {
		path     string
		expected bool
	}{
		{".git", true},
		{"path/to/node_modules", true},
		{"project/dist", true},
		{"src", false},
		{"docs", false},
	}

	for _, tc := range testCases {
		if result := shouldIgnoreDir(tc.path); result != tc.expected {
			t.Errorf("shouldIgnoreDir(%s) = %v, expected %v", tc.path, result, tc.expected)
		}
	}
}
