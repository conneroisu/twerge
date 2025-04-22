# Tailwind Build Integration

This page shows how to integrate Twerge with your Tailwind CSS build process for optimal performance and developer experience.

## Build Process Overview

A typical Twerge-Tailwind integration includes these steps:

1. Scan your templates for Tailwind classes
2. Generate optimized class mappings in Go code
3. Process CSS with Tailwind CLI
4. Serve the optimized CSS and HTML

## Basic Build Script

Here's a simple build script that handles code generation and Tailwind processing:

```go title="gen.go"
//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/conneroisu/twerge"
	"github.com/conneroisu/twerge/examples/simple/views"
)

var cwd = flag.String("cwd", "", "current working directory")

func main() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(update-css) Done in %s.\n", elapsed)
	}()
	flag.Parse()
	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Updating Generated Code...")
	start = time.Now()
	if err := twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		"input.css",
		"classes/classes.html",
		views.View(),
	); err != nil {
		panic(err)
	}
	fmt.Println("Done Generating Code. (took", time.Since(start), ")")

	fmt.Println("Running Tailwind...")
	start = time.Now()
	runTailwind()
	fmt.Println("Done Running Tailwind. (took", time.Since(start), ")")
}

func runTailwind() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(tailwind) Done in %s.\n", elapsed)
	}()
	cmd := exec.Command("tailwindcss", "-i", "input.css", "-o", "_static/dist/styles.css")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
```

## Tailwind Configuration

Here's a sample `tailwind.config.js` that works well with Twerge:

```js title="tailwind.config.js"
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './classes/classes.html', // Generated HTML classes from Twerge
    './views/**/*.templ',     // Optional - you can include your templates directly too
  ],
  theme: {
    extend: {
      colors: {
        primary: '#3b82f6',
        secondary: '#10b981',
        accent: '#8b5cf6',
      },
    },
  },
  plugins: [],
}
```

## Development Workflow

For development, you'll want to automatically rebuild when files change. Here's a simple watch script you can use:

```go title="watch.go"
//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start the build once at beginning
	runBuild()

	// Watch directories
	dirsToWatch := []string{"./views", "./input.css", "./tailwind.config.js"}
	for _, dir := range dirsToWatch {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Watch for changes recursively in views directory
	err = filepath.Walk("./views", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Debounce builds
	var lastBuild time.Time
	debounceInterval := 500 * time.Millisecond

	log.Println("Watching for changes...")
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Only rebuild on write or create events for .templ, .css, or .js files
			if (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) &&
				(filepath.Ext(event.Name) == ".templ" || filepath.Ext(event.Name) == ".css" || filepath.Ext(event.Name) == ".js") {
				// Debounce
				if time.Since(lastBuild) > debounceInterval {
					lastBuild = time.Now()
					log.Println("Change detected, rebuilding...")
					runBuild()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func runBuild() {
	// First generate templ components
	cmdTempl := exec.Command("templ", "generate", "./views")
	cmdTempl.Stdout = os.Stdout
	cmdTempl.Stderr = os.Stderr
	if err := cmdTempl.Run(); err != nil {
		log.Println("Error generating templ:", err)
		return
	}
	
	// Then run our build script
	cmdGen := exec.Command("go", "run", "gen.go")
	cmdGen.Stdout = os.Stdout
	cmdGen.Stderr = os.Stderr
	if err := cmdGen.Run(); err != nil {
		log.Println("Error running gen.go:", err)
		return
	}
	
	log.Println("Build completed successfully")
}
```

## Production Build Optimization

For production builds, you'll want to minify your CSS and use Tailwind's purge feature to remove unused styles:

```go title="gen_prod.go"
//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/conneroisu/twerge"
	"github.com/conneroisu/twerge/examples/simple/views"
)

var (
	cwd = flag.String("cwd", "", "current working directory")
	prod = flag.Bool("prod", false, "production build (minified)")
)

func main() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(update-css) Done in %s.\n", elapsed)
	}()
	flag.Parse()
	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	
	fmt.Println("Updating Generated Code...")
	start = time.Now()
	if err := twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		"input.css",
		"classes/classes.html",
		views.View(),
	); err != nil {
		panic(err)
	}
	fmt.Println("Done Generating Code. (took", time.Since(start), ")")

	fmt.Println("Running Tailwind...")
	start = time.Now()
	runTailwind(*prod)
	fmt.Println("Done Running Tailwind. (took", time.Since(start), ")")
}

func runTailwind(prod bool) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(tailwind) Done in %s.\n", elapsed)
	}()
	
	args := []string{"-i", "input.css", "-o", "_static/dist/styles.css"}
	if prod {
		args = append(args, "--minify")
	}
	
	cmd := exec.Command("tailwindcss", args...)
	cmd.Env = append(os.Environ(), "NODE_ENV=production")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
```

Usage:
```sh
# Development build
go run gen_prod.go

# Production build (minified)
go run gen_prod.go -prod
```

## Makefile Integration

You can create a Makefile to simplify common build tasks:

```makefile title="Makefile"
.PHONY: dev build watch clean prod

dev:
	templ generate ./views
	go run gen.go

watch:
	go run watch.go

build:
	templ generate ./views
	go run gen.go
	go build -o app ./main.go

prod:
	templ generate ./views
	go run gen_prod.go -prod
	go build -o app -ldflags="-s -w" ./main.go

clean:
	rm -f app
	rm -rf _static/dist/*
```

## GitHub Actions Example

You can automate the build process in CI/CD with GitHub Actions:

```yaml title=".github/workflows/build.yml"
name: Build and Deploy

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Install templ
      run: go install github.com/a-h/templ/cmd/templ@latest

    - name: Install tailwindcss
      run: npm install -g tailwindcss

    - name: Generate templ components
      run: templ generate ./views

    - name: Build production CSS
      run: go run gen_prod.go -prod

    - name: Build application
      run: go build -o app -ldflags="-s -w" ./main.go

    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: app-build
        path: |
          app
          _static/dist
```

## Docker Integration

For containerized deployments, here's a sample Dockerfile:

```dockerfile title="Dockerfile"
# Build stage
FROM golang:1.24-alpine AS builder

# Install Node.js and npm for Tailwind
RUN apk add --no-cache nodejs npm

# Install Tailwind CSS
RUN npm install -g tailwindcss

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Generate templ files
RUN templ generate ./views

# Build with Twerge
RUN go run gen_prod.go -prod

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app -ldflags="-s -w" ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy only necessary files from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/_static/dist _static/dist

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./app"]
```

## Multi-Theme Support

You can extend your build script to generate multiple theme variants:

```go title="gen_themes.go"
//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/conneroisu/twerge"
	"github.com/conneroisu/twerge/examples/simple/views"
)

var (
	cwd = flag.String("cwd", "", "current working directory")
	theme = flag.String("theme", "default", "theme name (default, dark, custom)")
)

func main() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(update-css) Done in %s.\n", elapsed)
	}()
	flag.Parse()
	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	
	// Determine input file based on theme
	inputCss := "input.css"
	if *theme != "default" {
		inputCss = fmt.Sprintf("input-%s.css", *theme)
	}
	
	// Determine output directory based on theme
	outputDir := filepath.Join("_static", "dist")
	outputCss := filepath.Join(outputDir, "styles.css")
	if *theme != "default" {
		outputCss = filepath.Join(outputDir, fmt.Sprintf("styles-%s.css", *theme))
	}
	
	fmt.Printf("Building theme: %s\n", *theme)
	fmt.Printf("Input CSS: %s\n", inputCss)
	fmt.Printf("Output CSS: %s\n", outputCss)
	
	fmt.Println("Updating Generated Code...")
	start = time.Now()
	if err := twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		inputCss,
		"classes/classes.html",
		views.View(),
	); err != nil {
		panic(err)
	}
	fmt.Println("Done Generating Code. (took", time.Since(start), ")")

	fmt.Println("Running Tailwind...")
	start = time.Now()
	runTailwind(inputCss, outputCss)
	fmt.Println("Done Running Tailwind. (took", time.Since(start), ")")
}

func runTailwind(input, output string) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(tailwind) Done in %s.\n", elapsed)
	}()
	
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}
	
	cmd := exec.Command("tailwindcss", "-i", input, "-o", output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
```

Usage:
```sh
# Build default theme
go run gen_themes.go

# Build dark theme
go run gen_themes.go -theme dark
```

This example demonstrates how to integrate Twerge into your Tailwind CSS build process for both development and production environments.