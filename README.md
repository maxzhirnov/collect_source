# GoSourceCollector

Simple utility to collect all source files from your Go project into a single text file. Perfect for code sharing, documentation, and review purposes.

## Features

- Collects all source files (configurable extensions: `.go`, `.html`, etc.)
- Generates project tree structure
- Adds meta information (generation time, Go version)
- Preserves relative paths for each file
- Excludes itself from parsing
- Skips `.git` and `vendor` directories

## Usage

1. Copy `collect_source.go` to your project directory
2. Run:

```bash
go run collect_source.go
```

3. Find the generated `project_source.txt` file in the same directory

## Configuration

Edit these constants in the code to customize behavior:

```go
const (
    outputFile    = "project_source.txt"  // Output file name
    thisFileName  = "collect_source.go"   // This script's filename to exclude
)

var fileExtensions = []string{
    ".go",
    ".html",
    // Add more extensions here
}
```

## Example Output

```
Project Source Code Export
Generated: 2024-11-21 14:00:00
Parsing extensions: [.go .html]
Working Directory: /your/project/path

----------------------------------------

Project Tree:
=============
├── main.go (1024 bytes)
├── handlers/
│   ├── api.go (2048 bytes)
│   └── templates/
│       └── index.html (512 bytes)
└── go.mod (128 bytes)

----------------------------------------

Source Code:
============

// File: main.go
// Size: 1024 bytes
// Extension: .go
----------------------------------------
[file content here]

=====================================
```

## License

MIT License - feel free to use and modify as you need.

## Author

Maksim Zhirnov