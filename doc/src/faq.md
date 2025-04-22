# Frequently Asked Questions

This page answers common questions about using Twerge with Tailwind CSS and Go.

## General Questions

### What is Twerge?

Twerge is a Go library that optimizes Tailwind CSS class usage in Go web applications. It provides class merging, short class name generation, and code generation features to improve performance and developer experience.

### Why should I use Twerge?

Twerge solves several common challenges when working with Tailwind CSS:
- It correctly merges Tailwind classes, resolving conflicts according to Tailwind's specificity rules
- It generates short, unique class names to reduce HTML size
- It creates Go code with class mappings for improved performance and type safety
- It integrates smoothly with Go build processes and templ templates

### How does Twerge compare to similar tools?

Unlike JavaScript-based tools like `tailwind-merge` or `clsx`, Twerge is designed specifically for Go applications. It integrates natively with Go code, templ templates, and build processes to provide a seamless development experience.

## Technical Questions

### Does Twerge support all Tailwind CSS features?

Twerge supports all standard Tailwind CSS classes, including:
- Layout and positioning
- Flexbox and Grid
- Typography
- Colors and backgrounds
- Borders and shadows
- Transitions and animations
- Responsive variants
- Dark mode variants
- Hover, focus, and other state variants

For custom Tailwind plugins or very specific edge cases, check the documentation or submit an issue on GitHub.

### How does class conflict resolution work?

Twerge follows Tailwind's own conflict resolution rules:
1. The last conflicting class wins (e.g., `text-red-500 text-blue-500` results in blue text)
2. Classes are grouped by category for readability and optimization
3. An internal mapping handles all standard Tailwind class conflicts

### What's the performance impact?

Twerge is designed for performance:
- Class merging uses an LRU cache for frequently used combinations
- Generated code offers zero runtime overhead for class lookups
- Short class names reduce HTML size and improve parsing time
- The build-time approach means no client-side JavaScript overhead

Tests show that using Twerge can reduce HTML size by 30-50% for Tailwind-heavy pages.

### Can I use Twerge with existing projects?

Yes, you can integrate Twerge into existing Go web projects that use Tailwind CSS. The integration process is straightforward:

1. Install Twerge: `go get github.com/conneroisu/twerge`
2. Update your templates to use `twerge.It()` or `twerge.Merge()`
3. Set up code generation in your build process
4. Configure TailwindCSS to use the generated HTML classes

See the Examples section for detailed integration guides.

## Common Issues

### Classes aren't being applied correctly

If classes aren't being applied correctly, check:
1. Make sure the generated CSS is being included in your HTML
2. Verify that your templates are using `twerge.It()` correctly
3. Check that code generation is running before the Tailwind build
4. Inspect the generated classes.html file to ensure your classes are included

### Build errors in code generation

If you're seeing build errors:
1. Ensure you're using the latest version of Twerge
2. Check that your templ templates are valid and compiling correctly
3. Verify that the paths in your CodeGen function are correct
4. Check the generated code for any issues

### Performance concerns

If you're concerned about performance:
1. Use the code generation approach for production builds
2. Configure appropriate cache sizes for your application
3. Consider enabling minification in your Tailwind build
4. Use HTTP compression for serving HTML and CSS

## Usage Examples

### How do I merge classes conditionally?

```go
// Conditional class merging
func Button(primary bool) templ.Component {
    classes := "px-4 py-2 rounded"
    if primary {
        classes = twerge.Merge(classes + " bg-blue-600 text-white")
    } else {
        classes = twerge.Merge(classes + " bg-gray-200 text-gray-800")
    }
    
    return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return templ.Tag("button", templ.Attributes{
            "class": classes,
        }).Render(ctx, w)
    })
}
```

### How do I integrate with CI/CD pipelines?

For CI/CD pipelines, include the Twerge code generation step in your build process:

```yaml
# GitHub Actions example
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - name: Install dependencies
        run: |
          go install github.com/a-h/templ/cmd/templ@latest
          npm install -g tailwindcss
      - name: Generate templ code
        run: templ generate ./views
      - name: Run Twerge code generation
        run: go run ./gen.go
      - name: Build application
        run: go build -o app
```

## Getting Help

If you have questions not covered in this FAQ:

1. Check the [documentation](https://github.com/conneroisu/twerge/docs)
2. Search for [existing issues](https://github.com/conneroisu/twerge/issues)
3. Open a new issue if you think you've found a bug
4. Join the community discussions on GitHub or Discord