# Codebaseify

A simple Go utility to convert your project's codebase into a single Markdown file, 
making it easier to share with AI assistants like ChatGPT, Claude, or other Large Language Models.

## Features

- Respects `.gitignore` rules when determining which files to include
- Automatically ignores binary files
- Creates a single Markdown file with the entire codebase
- Each file is properly formatted with code blocks and relative paths

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/codebaseify
cd codebaseify

# Build the binary
go build

# Optionally install it to your GOPATH
go install
```

## Usage

```bash
# Navigate to the root of your project
cd /path/to/your/project

# Run codebaseify
codebaseify

# Or specify a custom output file name
codebaseify -output custom-name.md
```

The tool will generate a file named project-prompt.md (or your custom name) in the current directory,
containing the contents of all files in the project, organized by their paths.

## Example Output
<!-- markdown content -->
```text
`main.go`
```go
package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, world!")
}

`util/helper.go`
```go
package util

func Helper() string {
    return "I'm helping!"
}
```

