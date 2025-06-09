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

## Template Development Best Practices

### Templ Integration Guidelines

**Template Syntax Requirements:**
- Multi-line `twerge.It()` calls must use simple string literals, avoid complex concatenations
- For complex conditional styling, break into multiple simple calls
- Use single-line twerge calls when possible for better template parsing

```go
// ✅ GOOD - Simple, single-line calls
class={twerge.It("bg-white dark:bg-gray-800 rounded-lg shadow-sm")}

// ✅ GOOD - Simple conditionals
class={twerge.If(isActive, "bg-blue-600", "bg-gray-300")}

// ❌ AVOID - Complex multi-line concatenations
class={twerge.It(
    "complex-base-classes " +
    twerge.If(condition, "true-classes", "false-classes") +
    " more-classes"
)}
```

**Template Structure Patterns:**
- Keep component templates focused and simple
- Use descriptive component names that reflect their purpose
- Separate complex logic into Go functions, not template expressions

### Project Setup Workflow

**Complete Development Setup:**
1. **Initialize Project Structure**
   ```bash
   mkdir my-twerge-project/{_static/dist,classes,views,handlers,types,data}
   ```

2. **Setup Go Module and Workspace**
   ```bash
   go mod init my-project
   # Add to parent go.work if using workspace
   echo "./my-project" >> ../go.work
   ```

3. **Create Basic Templates**
   - Start with simple, working templates
   - Generate incrementally with `templ generate`
   - Test each template before adding complexity

4. **Setup Code Generation**
   ```go
   //go:build ignore
   package main
   
   import "github.com/conneroisu/twerge"
   
   func main() {
       twerge.CodeGen(/* component instances */)
   }
   ```

5. **CSS Build Pipeline**
   ```bash
   # Generate optimized classes
   go run gen.go
   
   # Build CSS with TailwindCSS
   tailwindcss -i input.css -o _static/dist/styles.css --minify
   
   # Generate templates
   templ generate
   
   # Run application
   go run main.go
   ```

### Real-World Integration Patterns

**E-commerce Example Learnings:**
- **Component Hierarchies**: Build simple components first, then compose into complex views
- **State Management**: Use Go handlers to manage application state, keep templates pure
- **Static Assets**: Use embedded file systems for production deployments
- **HTMX Integration**: Twerge classes work seamlessly with HTMX partial updates

**Proven Architecture Pattern:**
```
project/
├── types/           # Data models and types
├── data/           # Data layer and business logic
├── handlers/       # HTTP handlers and API endpoints
├── views/          # Templ templates and components
├── _static/        # Static assets and generated CSS
├── classes/        # Generated Twerge class files
├── gen.go          # Code generation script
├── input.css       # TailwindCSS input with Twerge markers
└── main.go         # Application entry point
```

## Debugging and Troubleshooting

### Common Issues
- **Classes not merging correctly**: Check conflict resolution in `config.go:100-200`
- **Generated CSS missing**: Ensure `CodeGen()` runs before TailwindCSS build
- **Performance issues**: Use static generation for production builds
- **Test failures**: Run `go test -v ./...` for detailed output

### Template-Specific Issues
- **Template parsing errors**: Simplify multi-line twerge calls, avoid complex string concatenations
- **Missing generated files**: Run `templ generate` after creating/modifying .templ files
- **Class conflicts**: Use Twerge's conflict resolution instead of manual class management
- **Build failures**: Check go.work includes all example directories

### Debug Mode
Enable verbose output by setting debug flags in `debug.go`. Use `go run -tags debug` for detailed class processing information.

### Verification Steps
1. Run `go test ./...` to ensure core functionality
2. Test code generation: `go run examples/simple/gen.go`
3. Verify CSS output in `_static/dist/` directories
4. Check TailwindCSS integration with `tailwindcss --help`
5. **Test template generation**: `templ generate -path views/`
6. **Verify server startup**: Test basic HTTP endpoints respond correctly

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

// Dynamic state management in templates
class={twerge.It(
    twerge.If(product.Featured, "ring-2 ring-blue-500", "") +
    twerge.If(!product.InStock, "opacity-75 grayscale", "")
)}

// Complex conditional layouts
class={twerge.If(viewType == "grid",
    "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6",
    "flex flex-col space-y-4"
)}
```

### Production Code Generation Workflow
1. **Development Phase**
   ```bash
   # Create comprehensive component instances
   go run gen.go
   ```

2. **Template Compilation**
   ```bash
   # Generate Go files from templates
   templ generate
   ```

3. **CSS Optimization**
   ```bash
   # Build optimized CSS
   tailwindcss -i input.css -o _static/dist/styles.css --minify
   ```

4. **Testing and Validation**
   ```bash
   # Test the complete application
   go run main.go
   curl http://localhost:8080  # Verify functionality
   ```

### Component Generation Best Practices

**Comprehensive Component Coverage:**
```go
// Generate all component states for complete optimization
twerge.CodeGen(
    twerge.Default(),
    "classes/classes.go",
    "input.css", 
    "classes/classes.html",
    
    // All product card states
    views.ProductCard(normalProduct, false, false),
    views.ProductCard(normalProduct, true, false),   // In cart
    views.ProductCard(saleProduct, false, true),     // Favorited
    views.ProductCard(outOfStockProduct, false, false), // Out of stock
    
    // All interactive states
    views.CartItem(item, false),  // Normal
    views.CartItem(item, true),   // Updating
    
    // All filter combinations
    views.FilterSidebar(emptyFilters, 0, counts),
    views.FilterSidebar(activeFilters, 42, counts),
)
```

**Performance Optimization Strategy:**
- Generate instances covering all visual states
- Include responsive variations (mobile, tablet, desktop)
- Cover all interactive states (hover, active, disabled)
- Test with realistic data volumes

## Production Deployment Considerations

### Build Process Automation
```bash
#!/bin/bash
# production-build.sh

echo "🎨 Generating optimized classes..."
go run gen.go

echo "🏗️  Building CSS..."
tailwindcss -i input.css -o _static/dist/styles.css --minify

echo "⚡ Generating templates..."
templ generate

echo "🔨 Building Go binary..."
go build -o app main.go

echo "✅ Build complete!"
```

### Performance Metrics
Real-world Twerge applications typically achieve:
- **50-70% reduction** in CSS class name length
- **20-30% smaller** HTML payload sizes
- **Faster DOM updates** due to shorter class strings
- **Improved caching** with consistent optimized class names

### Integration Patterns Validated

**HTMX + Twerge**: Seamless integration for dynamic UIs
```go
// Optimized classes work perfectly with HTMX partial updates
<div hx-post="/update" hx-target="#content" class="tw-optimized-button">
```

**Responsive Design**: All TailwindCSS responsive features preserved
```go
// Responsive breakpoints maintained in optimized output
twerge.It("grid-cols-1 sm:grid-cols-2 lg:grid-cols-4")
// → "tw-responsive-grid"
```

**Dark Mode**: Theme switching works flawlessly
```go
// Dark mode classes properly optimized and grouped
twerge.It("bg-white dark:bg-gray-800 text-gray-900 dark:text-white")
// → "tw-theme-container"
```

### Monitoring and Maintenance

**Key Metrics to Track:**
- CSS bundle size before/after Twerge optimization
- Template compilation time
- Runtime performance of class applications
- Browser cache hit rates for CSS files

**Maintenance Best Practices:**
- Regenerate classes when adding new components
- Monitor for unused class definitions in generated CSS
- Regular testing of all component states
- Version control generated files for deployment consistency

## Dependencies

- `github.com/a-h/templ` - Go templating engine integration
- `github.com/dave/jennifer` - Go code generation
- `github.com/go-chi/chi/v5` - HTTP router (for web applications)
- Standard library only for core functionality

## Real-World Examples

The `examples/` directory contains production-ready implementations:

- **`simple/`**: Basic integration pattern, minimal setup
- **`dashboard/`**: Complex layouts with multiple component states  
- **`admin-dashboard/`**: Advanced admin interface with comprehensive features
- **`ecommerce-catalog/`**: Full e-commerce application with cart, filters, and product management

Each example demonstrates different aspects of Twerge integration and can serve as starting points for real applications.