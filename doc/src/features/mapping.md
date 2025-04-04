# Class Mapping

Twerge maintains mappings between original Tailwind class strings and their generated short names, providing powerful tools for managing these mappings.

## What are Class Mappings?

Class mappings are associations between:

- Original Tailwind class strings (e.g., `"flex items-center justify-between p-4"`)
- Generated short class names (e.g., `"tw-a1b2c3d4"`)

These mappings are essential for:

1. Consistent class name generation
2. CSS generation from shortened class names
3. Optimizing build processes
4. Debugging and development

## Working with Mappings

### Getting the Current Mapping

```go
import "github.com/conneroisu/twerge"

// Get the current mapping of original classes to generated names
mapping := twerge.GetMapping()

// Result: map[string]string{
//   "flex items-center justify-between p-4": "tw-a1b2c3d4",
//   "text-lg font-bold text-gray-800": "tw-e5f6g7h8",
//   ...
// }
```

### Manipulating Mappings

```go
// Merge multiple class maps
map1 := twerge.GetMapping()
map2 := loadMappingsFromSomewhere()
mergedMap := twerge.MergeCSSMaps(map1, map2)

// Result: A combined map with all entries from both maps
```

### Registering Custom Mappings

```go
// Define custom class mappings
customClasses := map[string]string{
    "flex items-center justify-between": "tw-header",
    "text-sm font-medium text-gray-500": "tw-label",
}

// Register these custom mappings
twerge.RegisterClasses(customClasses)
```

## Persistence and Sharing

### Exporting Mappings to CSS

```go
// Export the current mapping to CSS
err := twerge.ExportCSS("styles.css")
if err != nil {
    // Handle error
}

// Export a specific mapping to CSS
err = twerge.ExportCSSWithMap("component-styles.css", componentMap)
if err != nil {
    // Handle error
}
```

### Appending to Existing CSS

```go
// Append classes to an existing CSS file
err := twerge.AppendClassesToFile("styles.css", newClasses, "/* New Component */")
if err != nil {
    // Handle error
}
```

## Code Generation with Mappings

One of the most powerful features of Twerge is the ability to generate Go code from class mappings:

```go
// Generate Go code for a variable containing the class mapping
code := twerge.GenerateClassMapCode()

// Write the generated code to a file
err := twerge.WriteClassMapFile("tailwind_classes_gen.go")
if err != nil {
    // Handle error
}
```

The generated code might look something like:

```go
// Code generated by Twerge - DO NOT EDIT.

package mypackage

var TailwindClasses = map[string]string{
    "flex items-center justify-between p-4": "tw-a1b2c3d4",
    "text-lg font-bold text-gray-800": "tw-e5f6g7h8",
    // ...
}
```

## Use Cases for Mappings

### Component Libraries

Create component libraries with consistent class naming:

```go
// In your component package
package components

import "github.com/conneroisu/twerge"

// Generate component-specific class names
func init() {
    componentClasses := map[string]string{
        "flex items-center p-2 rounded": "btn-base",
        "bg-blue-500 text-white": "btn-primary",
        "bg-gray-200 text-gray-800": "btn-secondary",
    }
    twerge.RegisterClasses(componentClasses)
}

// In your component
templ Button(variant string) {
    <button class={ twerge.RuntimeGenerate("flex items-center p-2 rounded " + variant) }>
        { children... }
    </button>
}
```

### Optimized Build Process

Use mappings to optimize your build process:

1. During development, use the runtime functions
2. At build time, extract the mappings
3. Generate optimized code with the mappings
4. Use the generated code in production

This approach gives you the flexibility of runtime generation during development with the performance of build-time optimization in production.
