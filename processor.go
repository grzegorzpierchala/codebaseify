package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// List of binary file extensions to skip
var binaryExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".pdf": true, ".doc": true, ".docx": true, ".ppt": true,
	".pptx": true, ".xls": true, ".xlsx": true, ".zip": true,
	".gz": true, ".tar": true, ".mp3": true, ".mp4": true,
	".mov": true, ".avi": true, ".exe": true, ".dll": true,
	".so": true, ".o": true, ".a": true, ".dylib": true,
	".class": true, ".jar": true, ".pyc": true, ".pyo": true,
}

// isBinaryFile does a simple check if a file might be binary
func isBinaryFile(path string) bool {
	// Check by extension first
	ext := filepath.Ext(path)
	if binaryExtensions[ext] {
		return true
	}

	// Read a sample of the file
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the first 512 bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil || n == 0 {
		return false
	}

	// Check for null bytes (common in binary files)
	return bytes.ContainsRune(buf[:n], 0)
}

// determineCodeFence selects the appropriate code fence for the file type
func determineCodeFence(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".java":
		return "java"
	case ".c", ".cpp", ".h", ".hpp":
		return "cpp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".sh":
		return "bash"
	case ".md":
		return "markdown"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".sql":
		return "sql"
	case ".rs":
		return "rust"
	case ".kt":
		return "kotlin"
	case ".swift":
		return "swift"
	case ".cs":
		return "csharp"
	case ".fs":
		return "fsharp"
	default:
		return "text" // Default to text for unknown file types
	}
}

// processProject walks through the project directory and processes files
func processProject(rootDir string, outputFile *os.File) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			if shouldIgnoreDir(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Process files
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// Skip files that match ignore patterns
		if shouldIgnoreFile(relPath) {
			return nil
		}

		// Skip binary files
		if isBinaryFile(path) {
			fmt.Printf("Skipping binary file: %s\n", relPath)
			return nil
		}

		// Read and include file content
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}

		// Write to the output file
		normalizedPath := filepath.ToSlash(relPath) // Convert OS-specific path to forward slashes

		// Get the file extension
		ext := filepath.Ext(path)

		// Special handling for markdown files
		if ext == ".md" {
			// For markdown files, we need to escape any existing code blocks
			// We'll replace triple backticks with a placeholder
			mdContent := string(content)

			// Replace actual backtick code blocks with HTML code blocks for display
			mdContent = strings.ReplaceAll(mdContent, "```", "~~~")

			_, err = outputFile.WriteString(fmt.Sprintf("`%s` (Markdown file)\n\n%s\n\n", normalizedPath, mdContent))
		} else {
			// For regular code files, use the appropriate code fence
			codeFence := determineCodeFence(path)
			if codeFence == "" {
				codeFence = "text" // Default to text for unknown file types
			}

			_, err = outputFile.WriteString(fmt.Sprintf("`%s`\n```%s\n%s\n```\n\n", normalizedPath, codeFence, string(content)))
		}

		if err != nil {
			return fmt.Errorf("error writing to output file: %v", err)
		}
		fmt.Printf("Added %s to documentation\n", relPath)

		return nil
	})
}
