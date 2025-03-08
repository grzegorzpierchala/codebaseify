package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	outputFile = "project-prompt.md"
	rootDir    string
)

func main() {
	// Parse command line flags
	flag.StringVar(&outputFile, "output", outputFile, "Name of the output markdown file")
	flag.Parse()

	var err error
	rootDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Parse .gitignore if it exists
	if _, err := os.Stat(".gitignore"); err == nil {
		parseGitignore(".gitignore")
		fmt.Println("Loaded .gitignore file")
	} else {
		fmt.Println("No .gitignore file found. Will include all files.")
	}

	// Create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("Generating project documentation in Markdown format...")

	// Process the project files
	err = processProject(rootDir, file)
	if err != nil {
		fmt.Printf("Error processing project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s with project documentation.\n", outputFile)
}
