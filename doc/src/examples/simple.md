# Basic Examples

This page provides basic examples of how to use Twerge in your Go applications.

## Class Merging Example


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
