# Basic Examples

This page provides basic examples of how to use Twerge in your Go applications.

## Simple Website Example

The "simple" example demonstrates a basic website layout with header, main content, and footer sections. It shows how to use Twerge to optimize Tailwind CSS classes in a Go templ application.

### Project Structure

```
simple/
├── _static/
│   └── dist/           # Directory for compiled CSS
├── classes/
│   ├── classes.go      # Generated Go code with class mappings
│   └── classes.html    # HTML output of class definitions 
├── gen.go              # Code generation script
├── go.mod              # Go module file
├── input.css           # TailwindCSS input file
├── main.go             # Web server
├── tailwind.config.js  # TailwindCSS configuration
└── views/
    ├── view.templ      # Templ template file
    └── view_templ.go   # Generated Go code from templ
```

### Code Generation

The `gen.go` file handles Twerge code generation and TailwindCSS processing:

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

### Template Usage

The `view.templ` file shows how to use Twerge in a templ template:

```go title="views/view.templ (excerpt)"
package views

import "github.com/conneroisu/twerge"

templ View() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>stellar</title>
			<link rel="stylesheet" href="/dist/styles.css"/>
		</head>
		<body class={ twerge.It("bg-gray-50 text-gray-900 flex flex-col min-h-screen") }>
			<header class={ twerge.It("bg-indigo-600 text-white shadow-md") }>
				<!-- Header content -->
			</header>
			<main class={ twerge.It("container mx-auto px-4 py-6 flex-grow") }>
				<!-- Page content -->
			</main>
			<footer class={ twerge.It("bg-gray-800 text-white py-6") }>
				<!-- Footer content -->
			</footer>
		</body>
	</html>
}
```

### Benefits Demonstrated

- **Class Optimization** - Long Tailwind class strings are converted to short, efficient class names
- **Build Integration** - Twerge integrates with the build process to generate optimized CSS
- **Maintainability** - Templates remain readable with full Tailwind class names
- **Performance** - Final HTML output uses short class names for improved performance

### Running the Example

1. Navigate to the example directory:
```sh
cd examples/simple
```

2. Generate the templ components:
```sh
templ generate ./views
```

3. Run the code generation:
```sh
go run gen.go
```

4. Run the server:
```sh
go run main.go
```

5. Open your browser and navigate to http://localhost:8080

## Multiple Component Example

For applications with multiple components, you can pass all components to the `CodeGen` function:

```go
if err := twerge.CodeGen(
	twerge.Default(),
	"classes/classes.go",
	"input.css", 
	"classes/classes.html",
	views.Header(),
	views.Footer(),
	views.Sidebar(),
	views.Content(),
); err != nil {
	panic(err)
}
```

This ensures all components' Tailwind classes are included in the optimization process.