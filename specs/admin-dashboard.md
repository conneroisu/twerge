# Admin Dashboard Use Case

## Overview
A responsive admin dashboard demonstrating Twerge's ability to handle complex layouts with conflicting responsive classes, sidebar navigation, and data visualization components.

## Key Features
- Responsive sidebar that collapses on mobile
- Complex grid layouts for metrics and charts
- Dark/light mode toggle
- Multi-level navigation with hover states
- Data tables with sortable columns

## Twerge Benefits Demonstrated

### 1. Responsive Class Conflict Resolution
```go
// Sidebar responsive behavior
templ Sidebar(collapsed bool) {
    <aside class={twerge.It(
        "w-64 lg:w-80 xl:w-96 " +
        "md:translate-x-0 " +
        twerge.If(collapsed, "-translate-x-full md:translate-x-0", "translate-x-0") +
        " transition-transform duration-300 ease-in-out " +
        "bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700"
    )}>
        // Sidebar content
    </aside>
}
```

**Without Twerge**: Manual class conflict resolution, verbose conditionals
**With Twerge**: Automatic precedence handling, clean conditional logic

### 2. Complex Layout Grid System
```go
// Dashboard metrics grid
templ MetricsGrid() {
    <div class={twerge.It(
        "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 " +
        "gap-4 lg:gap-6 xl:gap-8 " +
        "p-4 lg:p-6 xl:p-8"
    )}>
        @MetricCard("Users", "12,345", "↗ 12%", true)
        @MetricCard("Revenue", "$56,789", "↗ 8%", true)
        @MetricCard("Orders", "1,234", "↘ 3%", false)
        @MetricCard("Conversion", "3.2%", "↗ 0.5%", true)
    </div>
}
```

**Class Conflicts Resolved**:
- `gap-4 lg:gap-6 xl:gap-8` → Progressive gap sizing
- `p-4 lg:p-6 xl:p-8` → Responsive padding scaling

### 3. Interactive Navigation States
```go
templ NavItem(label, href string, active, hasSubmenu bool) {
    <li>
        <a href={templ.SafeURL(href)} class={twerge.It(
            "flex items-center px-4 py-2 text-sm font-medium rounded-lg " +
            "transition-colors duration-200 " +
            twerge.If(active, 
                "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200",
                "text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-800"
            ) +
            twerge.If(hasSubmenu, " group", "")
        )}>
            { label }
            if hasSubmenu {
                <svg class={twerge.It("ml-auto h-4 w-4 transition-transform group-hover:rotate-90")}>
                    // Chevron icon
                </svg>
            }
        </a>
    </li>
}
```

### 4. Data Table with Dynamic States
```go
templ DataTable(columns []string, sortColumn string, sortDesc bool) {
    <div class={twerge.It("overflow-x-auto bg-white dark:bg-gray-800 shadow-lg rounded-lg")}>
        <table class={twerge.It("min-w-full divide-y divide-gray-200 dark:divide-gray-700")}>
            <thead class={twerge.It("bg-gray-50 dark:bg-gray-900")}>
                <tr>
                    for _, col := range columns {
                        @TableHeader(col, col == sortColumn, sortDesc)
                    }
                </tr>
            </thead>
            <tbody class={twerge.It("bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700")}>
                // Table rows
            </tbody>
        </table>
    </div>
}

templ TableHeader(label string, sorted, desc bool) {
    <th class={twerge.It(
        "px-6 py-3 text-left text-xs font-medium uppercase tracking-wider " +
        "cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 " +
        twerge.If(sorted,
            "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20",
            "text-gray-500 dark:text-gray-400"
        )
    )}>
        <div class={twerge.It("flex items-center space-x-1")}>
            <span>{ label }</span>
            if sorted {
                <svg class={twerge.It(
                    "h-4 w-4 transition-transform " +
                    twerge.If(desc, "rotate-180", "rotate-0")
                )}>
                    // Sort icon
                </svg>
            }
        </div>
    </th>
}
```

## Performance Benefits

### Class Consolidation
**Before Twerge** (hypothetical output):
```css
.sidebar-mobile { /* 15 declarations */ }
.sidebar-tablet { /* 12 declarations */ }
.sidebar-desktop { /* 18 declarations */ }
.nav-item-active { /* 8 declarations */ }
.nav-item-hover { /* 6 declarations */ }
/* + dozens more classes */
```

**With Twerge**:
```css
/* twerge:begin */
.tw-1 { @apply w-64 lg:w-80 xl:w-96 md:translate-x-0 transition-transform duration-300 ease-in-out bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700; }
.tw-2 { @apply w-64 lg:w-80 xl:w-96 -translate-x-full md:translate-x-0 transition-transform duration-300 ease-in-out bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700; }
.tw-3 { @apply grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 lg:gap-6 xl:gap-8 p-4 lg:p-6 xl:p-8; }
/* twerge:end */
```

### Bundle Size Impact
- **Estimated CSS reduction**: 40-60% fewer class definitions
- **Runtime performance**: O(1) class lookups vs. string concatenation
- **Maintainability**: Single source of truth for component styles

## Code Generation Workflow

```go
//go:build ignore

package main

import (
    "github.com/conneroisu/twerge"
    "github.com/yourproject/admin/views"
)

func main() {
    if err := twerge.CodeGen(
        twerge.Default(),
        "styles/generated.go",
        "styles/input.css", 
        "styles/classes.html",
        views.Dashboard(),
        views.Sidebar(false),
        views.Sidebar(true),
        views.MetricsGrid(),
        views.DataTable([]string{"Name", "Email", "Role"}, "Name", false),
        views.NavItem("Dashboard", "/", true, false),
        views.NavItem("Users", "/users", false, true),
    ); err != nil {
        panic(err)
    }
}
```

## Expected Outcomes

1. **Developer Experience**: Simplified responsive design with automatic conflict resolution
2. **Performance**: Smaller CSS bundles and faster runtime class lookups
3. **Maintainability**: Consistent styling patterns across complex UI components
4. **Scalability**: Easy to add new dashboard sections without class naming conflicts