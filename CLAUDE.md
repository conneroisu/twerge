# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Twerge is a Go library that optimizes TailwindCSS usage in Go templ applications by intelligently merging Tailwind classes and generating short, unique class names for improved runtime performance. It provides conflict resolution for competing CSS classes and enables a static code generation workflow for production builds.

## File Structure

```
twerge/
├── twerge.go          # Core class processing engine with It() and If() functions
├── tw.go              # Code generation system with CodeGen() function
├── config.go          # TailwindCSS configuration and conflict resolution
├── config_test.go     # Configuration tests
├── twerge_test.go     # Main functionality tests
├── debug.go           # Debug utilities and verbose output
├── doc.go             # Package documentation
├── examples/          # Example integrations
│   ├── simple/        # Basic usage example
│   ├── dashboard/     # Complex dashboard example
│   ├── admin-dashboard/ # Advanced admin interface
│   └── ecommerce-catalog/ # E-commerce catalog example
├── internal/files/    # Internal file utilities
├── specs/             # Feature specifications
└── benchmarks/        # Performance benchmarks
```

### Key Files for Development
- `twerge.go:15-25` - Main It() function implementation
- `tw.go:45-60` - CodeGen() function for static generation
- `config.go:100-200` - Class conflict resolution logic
- `examples/*/gen.go` - Code generation patterns

## Build Commands

### Core Development Commands
- Run tests: `go test ./...`
- Run single test: `go test -run TestName ./path/to/package`
- Format code: `go fmt ./...`
- Generate code: `go generate ./...`

### Nix Development Environment
The project supports Nix Flake development with predefined scripts:
- `dx` - Edit flake.nix
- `tests` - Run all go tests with verbose output
- `lint` - Run golangci-lint, statix, and deadnix
- `format` - Format Go, JS, TS, CSS, MD, and JSON files
- `generate-all` - Generate documentation with gomarkdoc
- `coverage-tests` - Run tests with coverage report

### Testing Commands
- Run all tests: `go test ./...`
- Run tests with coverage: `go test -cover ./...`
- Run specific test function: `go test -run TestIt ./...`
- Run benchmark tests: `go test -bench=. ./benchmarks/...`
- Test with race detection: `go test -race ./...`
- Verbose test output: `go test -v ./...`

### Testing Patterns
- **Unit Tests**: Focus on `twerge_test.go` for core functionality
- **Config Tests**: Use `config_test.go` for class resolution logic  
- **Integration Tests**: Test examples with `go run examples/*/gen.go`
- **Benchmark Tests**: Performance tests in `benchmarks/` directory

### TailwindCSS Integration
- Build CSS: `tailwindcss -i input.css -o _static/dist/styles.css`
- Example generation: `go run examples/dashboard/gen.go`
- Watch mode: `tailwindcss -i input.css -o _static/dist/styles.css --watch`

## Architecture

### Core Components

**Class Processing Engine (`twerge.go`)**:
- `It(classes string)` - Main function that converts TailwindCSS classes to optimized class names
- `If(condition, trueClass, falseClass)` - Conditional class application
- Thread-safe caching with `sync.RWMutex` for concurrent access
- Class conflict resolution based on TailwindCSS precedence rules

**Code Generation System (`tw.go`)**:
- `CodeGen()` - Generates CSS, Go, and HTML files for static builds
- CSS generation with `@apply` directives between `/* twerge:begin */` and `/* twerge:end */` markers
- Go file generation creates static cache maps for runtime performance
- HTML file generation for comprehensive TailwindCSS purging

**Configuration System (`config.go`)**:
- Comprehensive TailwindCSS class group definitions and conflict mappings
- Validator functions for arbitrary values, colors, sizes, and complex patterns
- Modifier handling for responsive, pseudo-states, and custom variants

### Data Flow

1. **Development**: Use `twerge.It("bg-blue-500 text-white")` in templ files
2. **Code Generation**: Run `CodeGen()` to analyze all components and generate:
   - CSS file with optimized `@apply` rules
   - Go file with static class mappings
   - HTML file for TailwindCSS purging
3. **Build**: Process generated CSS through TailwindCSS CLI
4. **Runtime**: Fast lookups using pre-generated static maps

### Class Merging Logic

The merger resolves conflicts intelligently:
- `"text-red-500 text-blue-500"` → `"text-blue-500"` (later class wins)
- `"px-4 pl-2"` → `"px-4 pl-2"` (specific overrides general)
- Handles modifiers: `"hover:bg-red-500 focus:bg-blue-500"` → both preserved
- Important modifier support: `"!text-red-500"`

### Example Integration

Typical workflow in a templ project:
1. Create `gen.go` with `//go:build ignore` for code generation
2. Use `twerge.CodeGen()` to process all components
3. Run TailwindCSS build process on generated CSS
4. Import generated Go file for static caching in production

### Performance Considerations

- Thread-safe concurrent access to class cache
- Static code generation eliminates runtime class parsing
- Short generated class names reduce CSS bundle size
- Automatic conflict resolution prevents CSS specificity issues

## Code Patterns and Conventions

### Core API Usage
```go
// Basic class merging
result := twerge.It("bg-blue-500 text-white hover:bg-blue-600")

// Conditional classes
result := twerge.If(isActive, "bg-green-500", "bg-gray-500")

// Complex merging with conflicts
result := twerge.It("px-4 pl-2 text-red-500 text-blue-500") // → "px-4 pl-2 text-blue-500"
```

### Code Generation Pattern
```go
//go:build ignore

package main

import "github.com/conneroisu/twerge"

func main() {
    twerge.CodeGen()
}
```

### Testing Patterns
```go
func TestClassMerging(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"conflict resolution", "text-red-500 text-blue-500", "text-blue-500"},
        {"modifier preservation", "hover:bg-red-500 focus:bg-blue-500", "hover:bg-red-500 focus:bg-blue-500"},
    }
    // Test implementation...
}
```

### File Naming Conventions
- `gen.go` files use `//go:build ignore` for code generation
- Test files follow `*_test.go` pattern
- Example directories contain complete working implementations
- Generated files use `_templ.go` suffix for templ integration

### Import Patterns
```go
import (
    "github.com/conneroisu/twerge"
    "github.com/a-h/templ"
)
```

## Debugging and Troubleshooting

### Common Issues
- **Classes not merging correctly**: Check conflict resolution in `config.go:100-200`
- **Generated CSS missing**: Ensure `CodeGen()` runs before TailwindCSS build
- **Performance issues**: Use static generation for production builds
- **Test failures**: Run `go test -v ./...` for detailed output

### Debug Mode
Enable verbose output by setting debug flags in `debug.go`. Use `go run -tags debug` for detailed class processing information.

### Verification Steps
1. Run `go test ./...` to ensure core functionality
2. Test code generation: `go run examples/simple/gen.go`
3. Verify CSS output in `_static/dist/` directories
4. Check TailwindCSS integration with `tailwindcss --help`

## Quick Reference Examples

### Basic Integration
```go
// In your templ component
@twerge.It("flex items-center justify-between p-4")

// With conditions
@twerge.If(user.IsAdmin, "bg-red-500", "bg-blue-500")
```

### Advanced Usage
```go
// Complex class merging
classes := twerge.It("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-6")

// With responsive modifiers
classes := twerge.It("text-sm md:text-base lg:text-lg font-medium")
```

### Code Generation Workflow
1. Create `gen.go` with `//go:build ignore`
2. Import and call `twerge.CodeGen()`
3. Run `go run gen.go` to generate files
4. Build CSS with `tailwindcss -i input.css -o output.css`

## Dependencies

- `github.com/a-h/templ` - Go templating engine integration
- `github.com/dave/jennifer` - Go code generation
- Standard library only for core functionality