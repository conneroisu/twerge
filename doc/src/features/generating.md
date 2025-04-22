# Generating Short Class Names

Twerge can generate short, unique class names based on hashes of the merged Tailwind classes, allowing for smaller HTML output and improved performance.

## The Problem

Tailwind CSS is incredibly powerful, but it can lead to long class strings:

```html
<div
  class="flex flex-col items-center justify-between p-4 rounded-lg shadow-lg bg-white hover:bg-gray-50 dark:bg-gray-800 dark:hover:bg-gray-700 transition-all duration-300"
>
  <!-- Content -->
</div>
```

These long strings increase HTML size, reduce readability, and can impact performance.

## How Twerge Solves It

The `Generate` function creates short, unique class names based on the hash of the merged Tailwind classes:

```go
import "github.com/conneroisu/twerge"

// Long class string
classes := "flex flex-col items-center justify-between p-4 rounded-lg shadow-lg bg-white hover:bg-gray-50"

// Generate a short unique class name
shortClassName := twerge.It(classes)

// Result: "tw-1" (example - actual output will vary)
```

## Benefits of Generated Class Names

- **Smaller HTML** - Dramatically reduces HTML file size
- **Better Caching** - Improves browser caching of HTML
- **Consistent Naming** - Same class string always generates the same short name
- **Collision-Free** - Hash-based generation ensures uniqueness
- **Automatic Conflict Resolution** - Classes are merged before hashing

## How It Works

1. First, Twerge merges the provided classes to resolve conflicts
2. The merged class string is hashed using a fast algorithm (default is FNV-1a)
3. The hash is converted to a short alphanumeric string
4. A prefix is added (default is "tw-")
5. The mapping from original classes to the generated name is stored

## Customizing Generation

## Integration Example

In a Go-templ template:

```go
// Instead of this with long class strings
<div class="flex items-center justify-between p-4 bg-white dark:bg-gray-800">...</div>

// You can do this with a short generated class
<div class={ twerge.It("flex items-center justify-between p-4 bg-white dark:bg-gray-800") }>...</div>

// Which results in HTML like:
// <div class="tw-1">...</div>
// <div class="tw-2">...</div>
```

## Performance Considerations

- Generated class names are cached for performance
- The same input always produces the same output
- Short names reduce HTML size and parsing time
