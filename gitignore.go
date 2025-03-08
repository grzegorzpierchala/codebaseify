package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ignoredPaths    = make(map[string]bool)
	ignoredPatterns = []string{"*.log", "tmp/", "*.md"}       // Default patterns to ignore
	ignoredDirs     = []string{".git", ".gitignore", "build"} // Default directories to ignore
)

func parseGitignore(gitignorePath string) {
	file, err := os.Open(gitignorePath)
	if err != nil {
		fmt.Printf("Warning: Error opening .gitignore: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle negation (files that shouldn't be ignored)
		if strings.HasPrefix(line, "!") {
			// This is a complex case, simplifying for now
			continue
		}

		// Add to ignored patterns
		pattern := line
		if strings.HasSuffix(pattern, "/") {
			// This is a directory pattern
			ignoredDirs = append(ignoredDirs, strings.TrimSuffix(pattern, "/"))
		} else {
			ignoredPatterns = append(ignoredPatterns, pattern)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading .gitignore: %v\n", err)
	}
}

func shouldIgnoreDir(path string) bool {
	dir := filepath.Base(path)

	// Check if it's in our list of ignored directories
	for _, ignoreDir := range ignoredDirs {
		if matched, _ := filepath.Match(ignoreDir, dir); matched {
			return true
		}
		// Check for paths like "dir/**/something"
		if strings.Contains(ignoreDir, "**/") {
			pattern := strings.Replace(ignoreDir, "**", "*", -1)
			relPath, _ := filepath.Rel(rootDir, path)
			if matched, _ := filepath.Match(pattern, relPath); matched {
				return true
			}
		}
	}
	return false
}

func shouldIgnoreFile(relPath string) bool {
	// Normalize the path for consistent matching
	normalizedPath := filepath.ToSlash(relPath)

	// Check explicit path matches
	if ignoredPaths[normalizedPath] {
		return true
	}

	// Check if the file matches any gitignore pattern
	for _, pattern := range ignoredPatterns {
		if strings.Contains(pattern, "/") {
			// This is a path pattern
			if matched, _ := filepath.Match(pattern, relPath); matched {
				return true
			}
		} else {
			// This is a filename or extension pattern
			if matched, _ := filepath.Match(pattern, filepath.Base(relPath)); matched {
				return true
			}
			// Handle patterns like "*.txt"
			if strings.HasPrefix(pattern, "*") {
				ext := filepath.Ext(relPath)
				if ext != "" && strings.HasSuffix(pattern, ext) {
					return true
				}
			}
		}
	}

	// Check if the file is the output file itself
	if filepath.Base(relPath) == outputFile {
		return true
	}

	return false
}
