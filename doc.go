// Package twerge provides a tailwind merger for go-templ with class generation and runtime static hashmap.
//
// # Overview
//
// Twerge optimizes TailwindCSS usage in Go templ applications by intelligently merging Tailwind classes,
// resolving conflicts according to Tailwind's precedence rules, and generating short, unique class names
// for improved runtime performance and CSS output size.
//
// It performs four key functions:
//   - Merges TailwindCSS classes intelligently (resolving conflicts between competing classes)
//   - Generates short unique CSS class names from the merged classes (e.g., "tw-1", "tw-2")
//   - Creates a mapping from original class strings to generated class names
//   - Provides a code generation step for efficient runtime lookups
//
// # Basic Usage
//
// In your .templ files:
//
//	import "github.com/conneroisu/twerge"
//
//	templ Button(primary bool) {
//		<button class={twerge.It("px-4 py-2 rounded " + twerge.If(primary, "bg-blue-500 text-white", "bg-gray-200 text-gray-800"))}>
//			{ children... }
//		</button>
//	}
//
// # Class Merging Logic
//
// When multiple Tailwind classes that target the same CSS property are provided, twerge will
// intelligently merge them according to Tailwind's precedence rules:
//
//	// Input: "text-red-500 bg-blue-500 text-blue-700"
//	// Output: "bg-blue-500 text-blue-700"
//	// (text-blue-700 takes precedence over text-red-500)
//
// # Code Generation
//
// For optimal runtime performance, twerge includes a code generation step that creates:
//
//   - A CSS file with all the generated classes
//
//   - A Go file with a static map for class lookups
//
//   - An HTML file for rendering
//
//     g := twerge.Default()
//     err := g.CodeGen(
//     g,
//     "views/gen_twerge.go",  // Generated Go file path
//     "views/twerge.css",     // Generated CSS file path
//     "views/twerge.html",    // Generated HTML file path
//     views.Button(true),     // Components to analyze
//     views.Card(),
//     // ... other components
//     )
//
// The generated CSS will look like:
//
//	/* twerge:begin */
//	/* from px-4 py-2 rounded bg-blue-500 text-white */
//	.tw-1 {
//		@apply px-4 py-2 rounded bg-blue-500 text-white;
//	}
//	/* from px-4 py-2 rounded bg-gray-200 text-gray-800 */
//	.tw-2 {
//		@apply px-4 py-2 rounded bg-gray-200 text-gray-800;
//	}
//	/* twerge:end */
//
// # API Reference
//
// ## Basic Functions
//
//	// It returns a short unique CSS class name from the merged classes.
//	// This is the most commonly used function in templates.
//	func It(raw string) string
//
//	// If returns a class based on a condition.
//	// Useful for conditional styling.
//	func If(ok bool, trueClass string, falseClass string) string
//
//	// CodeGen generates all the code needed to use Twerge statically.
//	func CodeGen(g *Generator, goPath string, cssPath string, htmlPath string, comps ...templ.Component) error
//
// ## Generator
//
// The Generator struct provides more advanced functionality:
//
//	// Default returns the default Generator instance.
//	func Default() *Generator
//
//	// New creates a new Generator with the given non-nil Handler.
//	func New(h Handler) *Generator
//
//	// Cache returns the cache of the Generator.
//	func (Generator) Cache() map[string]CacheValue
//
//	// It returns a short unique CSS class name from the merged classes.
//	func (g *Generator) It(classes string) string
//
// ## Configuration
//
// Although most users will use the default configuration, customization is possible
// by implementing the Handler interface:
//
//	type Handler interface {
//		It(string) string
//		Cache() map[string]CacheValue
//		SetCache(map[string]CacheValue)
//	}
//
// # Implementation Details
//
// Twerge uses a sophisticated algorithm to:
//   - Parse and understand Tailwind class semantics
//   - Identify and resolve conflicting classes
//   - Handle complex class relationships (e.g., inset-x-1 and left-1)
//   - Support arbitrary values, modifiers, and variants
//   - Maintain proper precedence of important (!) modifiers
//
// # Best Practices
//
// 1. Run CodeGen as part of your build process
// 2. Use twerge.It() for all Tailwind classes in your templates
// 3. Use twerge.If() for conditional class application
// 4. Avoid manual string concatenation for class names when possible
//
// # Example Workflow
//
// A typical development workflow with twerge:
//
// 1. Write your templ components using twerge.It() for class handling
// 2. During build or development, run CodeGen to generate static assets
// 3. Include the generated CSS in your application
// 4. The templ components will use the generated optimized class names at runtime
package twerge

//go:generate gomarkdoc -o README.md -e .
