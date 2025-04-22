# Code Generation

Twerge provides powerful code generation capabilities, allowing you to generate Go code from class mappings for improved performance and type safety.

## The Problem

When using shortened class names in a production environment, you need:

1. A way to consistently map original Tailwind classes to short names
2. Fast lookups without runtime overhead
3. Type safety and compiler checks
4. Integration with build processes

## How Twerge Solves It

Twerge can generate Go code that contains the class mappings, providing compile-time checking and improved performance:

```go
package main

import "github.com/conneroisu/twerge"

func main() {
	if err := twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		"input.css",
		"classes/classes.html",
		views.View(),
	); err != nil {
		panic(err)
	}
}
```

## Benefits of Code Generation

Using generated code provides several advantages:

1. **Performance** - No runtime hash computation or map lookups
2. **Type Safety** - Compile-time checking of class names
3. **Smaller Binary** - Compiled code can be optimized by the Go compiler
4. **Build-time Validation** - Issues are caught during the build process
5. **IDE Support** - Auto-completion and refactoring support in IDEs
