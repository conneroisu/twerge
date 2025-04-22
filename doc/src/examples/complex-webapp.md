# Complex Web Application

This page demonstrates using Twerge in a more complex web application scenario, showcasing a dashboard example with multiple components and dynamic data.

## Dashboard Example Overview

The dashboard example demonstrates a modern analytics UI with:

- Responsive layout using Tailwind CSS
- Multiple templ components for different sections
- Metrics cards with dynamic data
- Data tables for displaying information
- Twerge optimization for class handling

## Project Structure

```
dashboard/
├── _static/
│   └── dist/             # Directory for compiled CSS
├── classes/
│   ├── classes.go        # Generated Go code with class mappings
│   └── classes.html      # HTML output of class definitions 
├── gen.go                # Code generation script
├── go.mod                # Go module file
├── input.css             # TailwindCSS input file
├── main.go               # Web server implementation
├── tailwind.config.js    # TailwindCSS configuration
└── views/
    ├── dashboard.templ   # Dashboard page component
    ├── report.templ      # Report page component
    ├── settings.templ    # Settings page component
    └── view.templ        # Layout component
```

## Code Generation with Multiple Components

The dashboard example uses multiple components, all processed by Twerge:

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
	"github.com/conneroisu/twerge/examples/dashboard/views"
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
		views.Dashboard(),
		views.Settings(),
		views.Report(),
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

## Dashboard Component

The dashboard component displays metrics cards and a data table:

```go title="views/dashboard.templ (excerpt)"
package views

import "github.com/conneroisu/twerge"

templ Dashboard() {
	@Layout("Dashboard") {
		<main class={ twerge.It("container mx-auto px-4 py-6 flex-grow") }>
			<div class={ twerge.It("mb-6") }>
				<h2 class={ twerge.It("text-2xl font-bold text-gray-800 mb-4") }>Overview</h2>
				<div class={ twerge.It("grid grid-cols-1 md:grid-cols-3 gap-6") }>
					<div class={ twerge.It("bg-white rounded-lg shadow-md p-6") }>
						<div class={ twerge.It("flex items-center justify-between") }>
							<h3 class={ twerge.It("text-gray-500 text-sm font-medium") }>Total Users</h3>
							<span class={ twerge.It("bg-green-100 text-green-800 text-xs font-semibold px-2 py-1 rounded") }>+12%</span>
						</div>
						<p class={ twerge.It("text-3xl font-bold text-gray-800 mt-2") }>24,521</p>
						<div class={ twerge.It("mt-4 text-sm text-gray-500") }>1,250 new users this week</div>
					</div>
					<!-- Additional metric cards -->
				</div>
			</div>
			<div class={ twerge.It("mb-6") }>
				<h2 class={ twerge.It("text-2xl font-bold text-gray-800 mb-4") }>Recent Activity</h2>
				<div class={ twerge.It("bg-white rounded-lg shadow-md overflow-hidden") }>
					<table class={ twerge.It("min-w-full divide-y divide-gray-200") }>
						<!-- Table header and rows -->
					</table>
				</div>
			</div>
		</main>
	}
}
```

## Settings Component

The settings page demonstrates form handling with Tailwind:

```go title="views/settings.templ (excerpt)"
package views

import "github.com/conneroisu/twerge"

templ Settings() {
	@Layout("Settings") {
		<main class={ twerge.It("container mx-auto px-4 py-6 flex-grow") }>
			<h2 class={ twerge.It("text-2xl font-bold text-gray-800 mb-6") }>Account Settings</h2>
			
			<div class={ twerge.It("bg-white rounded-lg shadow-md p-6 mb-6") }>
				<h3 class={ twerge.It("text-lg font-medium text-gray-800 mb-4") }>Profile Information</h3>
				<form>
					<div class={ twerge.It("grid grid-cols-1 gap-6 md:grid-cols-2") }>
						<div>
							<label for="name" class={ twerge.It("block text-sm font-medium text-gray-700 mb-1") }>Name</label>
							<input type="text" id="name" name="name" value="John Doe" 
								class={ twerge.It("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500") }/>
						</div>
						<div>
							<label for="email" class={ twerge.It("block text-sm font-medium text-gray-700 mb-1") }>Email</label>
							<input type="email" id="email" name="email" value="john@example.com" 
								class={ twerge.It("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500") }/>
						</div>
						<!-- More form fields -->
					</div>
					<div class={ twerge.It("mt-6") }>
						<button type="submit" class={ twerge.It("px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2") }>
							Save Changes
						</button>
					</div>
				</form>
			</div>
			
			<!-- Additional settings sections -->
		</main>
	}
}
```

## Report Component

The report page shows how to present data visualizations:

```go title="views/report.templ (excerpt)"
package views

import "github.com/conneroisu/twerge"

templ Report() {
	@Layout("Analytics Report") {
		<main class={ twerge.It("container mx-auto px-4 py-6 flex-grow") }>
			<div class={ twerge.It("flex justify-between items-center mb-6") }>
				<h2 class={ twerge.It("text-2xl font-bold text-gray-800") }>Analytics Report</h2>
				<div class={ twerge.It("flex space-x-2") }>
					<button class={ twerge.It("px-3 py-1 bg-white border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50") }>
						Export PDF
					</button>
					<button class={ twerge.It("px-3 py-1 bg-white border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50") }>
						Export CSV
					</button>
				</div>
			</div>
			
			<div class={ twerge.It("grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6") }>
				<div class={ twerge.It("bg-white rounded-lg shadow-md p-6") }>
					<h3 class={ twerge.It("text-lg font-medium text-gray-800 mb-4") }>Revenue Overview</h3>
					<div class={ twerge.It("h-64 bg-gray-100 rounded flex items-center justify-center") }>
						<p class={ twerge.It("text-gray-500") }>Chart Placeholder</p>
					</div>
				</div>
				
				<div class={ twerge.It("bg-white rounded-lg shadow-md p-6") }>
					<h3 class={ twerge.It("text-lg font-medium text-gray-800 mb-4") }>User Growth</h3>
					<div class={ twerge.It("h-64 bg-gray-100 rounded flex items-center justify-center") }>
						<p class={ twerge.It("text-gray-500") }>Chart Placeholder</p>
					</div>
				</div>
			</div>
			
			<!-- Additional report sections -->
		</main>
	}
}
```

## Layout Component with Template Composition

The layout component used by all pages for consistent structure:

```go title="views/view.templ"
package views

import "github.com/conneroisu/twerge"

templ Layout(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title } | Dashboard</title>
			<link rel="stylesheet" href="/dist/styles.css"/>
		</head>
		<body class={ twerge.It("bg-gray-100 text-gray-900 flex") }>
			<!-- Sidebar -->
			<aside class={ twerge.It("hidden md:flex md:flex-col w-64 bg-gray-800 text-white") }>
				<div class={ twerge.It("flex items-center justify-center h-16 border-b border-gray-700") }>
					<h1 class={ twerge.It("text-xl font-bold") }>Dashboard</h1>
				</div>
				<nav class={ twerge.It("flex-grow") }>
					<ul class={ twerge.It("mt-6") }>
						<li>
							<a href="/" class={ twerge.It("flex items-center px-4 py-3 text-gray-300 hover:bg-gray-700") }>
								<span class={ twerge.It("ml-2") }>Dashboard</span>
							</a>
						</li>
						<li>
							<a href="/reports" class={ twerge.It("flex items-center px-4 py-3 text-gray-300 hover:bg-gray-700") }>
								<span class={ twerge.It("ml-2") }>Reports</span>
							</a>
						</li>
						<li>
							<a href="/settings" class={ twerge.It("flex items-center px-4 py-3 text-gray-300 hover:bg-gray-700") }>
								<span class={ twerge.It("ml-2") }>Settings</span>
							</a>
						</li>
					</ul>
				</nav>
			</aside>
			
			<!-- Main content -->
			<div class={ twerge.It("flex flex-col flex-grow min-h-screen") }>
				<!-- Top navbar -->
				<header class={ twerge.It("bg-white shadow h-16 flex items-center justify-between px-6") }>
					<div class={ twerge.It("flex items-center") }>
						<button class={ twerge.It("md:hidden mr-4 text-gray-600") }>
							<svg class={ twerge.It("h-6 w-6") } fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
							</svg>
						</button>
						<h2 class={ twerge.It("text-lg font-medium") }>{ title }</h2>
					</div>
					<div class={ twerge.It("flex items-center") }>
						<button class={ twerge.It("p-1 mr-4 text-gray-500") }>
							<svg class={ twerge.It("h-6 w-6") } fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
							</svg>
						</button>
						<div class={ twerge.It("relative") }>
							<button class={ twerge.It("flex items-center") }>
								<img class={ twerge.It("h-8 w-8 rounded-full") } src="https://randomuser.me/api/portraits/men/1.jpg" alt="User profile"/>
								<span class={ twerge.It("ml-2 text-sm") }>John Doe</span>
							</button>
						</div>
					</div>
				</header>
				
				<!-- Page content -->
				{ children... }
			</div>
		</body>
	</html>
}
```

## Web Server Implementation

The main.go file sets up routes for each component:

```go title="main.go"
package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/conneroisu/twerge/examples/dashboard/views"
)

func main() {
	// Static file handler
	fs := http.FileServer(http.Dir("_static"))
	http.Handle("/dist/", http.StripPrefix("/dist/", fs))
	
	// Routes
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/reports", handleReports)
	http.HandleFunc("/settings", handleSettings)
	
	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	err := views.Dashboard().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleReports(w http.ResponseWriter, r *http.Request) {
	err := views.Report().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	err := views.Settings().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper to get environment variable with default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
```

## Dynamic Data Integration

In a real application, you'd likely integrate with a database or API. Here's an example of how to pass dynamic data to templates:

```go title="models/data.go"
package models

type User struct {
	ID       int
	Name     string
	Email    string
	Role     string
	Avatar   string
	LastSeen string
}

type MetricCard struct {
	Title    string
	Value    string
	Change   string
	IsUp     bool
	Subtitle string
}

type Order struct {
	ID       string
	Customer string
	Amount   string
	Status   string
	Date     string
}

// Repository interface
type DataRepository interface {
	GetUsers() []User
	GetMetrics() []MetricCard
	GetRecentOrders() []Order
}
```

```go title="data/repository.go"
package data

import "github.com/example/dashboard/models"

type Repository struct {
	// Database connection or other dependencies
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetMetrics() []models.MetricCard {
	return []models.MetricCard{
		{
			Title:    "Total Users",
			Value:    "24,521",
			Change:   "+12%",
			IsUp:     true,
			Subtitle: "1,250 new users this week",
		},
		{
			Title:    "Total Revenue",
			Value:    "$45,428",
			Change:   "+5%",
			IsUp:     true,
			Subtitle: "$2,150 new revenue this week",
		},
		{
			Title:    "Total Orders",
			Value:    "12,234",
			Change:   "-2%",
			IsUp:     false,
			Subtitle: "345 new orders this week",
		},
	}
}

// Implement other methods...
```

```go title="views/dashboard.templ (with dynamic data)"
package views

import (
	"github.com/conneroisu/twerge"
	"github.com/example/dashboard/models"
)

templ MetricCard(card models.MetricCard) {
	<div class={ twerge.It("bg-white rounded-lg shadow-md p-6") }>
		<div class={ twerge.It("flex items-center justify-between") }>
			<h3 class={ twerge.It("text-gray-500 text-sm font-medium") }>{ card.Title }</h3>
			if card.IsUp {
				<span class={ twerge.It("bg-green-100 text-green-800 text-xs font-semibold px-2 py-1 rounded") }>{ card.Change }</span>
			} else {
				<span class={ twerge.It("bg-red-100 text-red-800 text-xs font-semibold px-2 py-1 rounded") }>{ card.Change }</span>
			}
		</div>
		<p class={ twerge.It("text-3xl font-bold text-gray-800 mt-2") }>{ card.Value }</p>
		<div class={ twerge.It("mt-4 text-sm text-gray-500") }>{ card.Subtitle }</div>
	</div>
}

templ Dashboard(metrics []models.MetricCard, orders []models.Order) {
	@Layout("Dashboard") {
		<main class={ twerge.It("container mx-auto px-4 py-6 flex-grow") }>
			<div class={ twerge.It("mb-6") }>
				<h2 class={ twerge.It("text-2xl font-bold text-gray-800 mb-4") }>Overview</h2>
				<div class={ twerge.It("grid grid-cols-1 md:grid-cols-3 gap-6") }>
					for _, metric := range metrics {
						@MetricCard(metric)
					}
				</div>
			</div>
			<!-- Table with orders... -->
		</main>
	}
}
```

## Benefits of This Approach

The complex dashboard example demonstrates several benefits of using Twerge:

1. **Component Reusability** - Common UI elements can be extracted into reusable components
2. **Optimized Output** - Large Tailwind class strings are converted to short, efficient codes
3. **Type Safety** - Generated Go code provides compile-time checking
4. **Performance** - HTML output is smaller and faster to parse
5. **Maintainability** - Templates remain readable with full Tailwind class names
6. **Dynamic Data Integration** - Easy to integrate with databases or APIs

## Running the Example

1. Navigate to the example directory:
```sh
cd examples/dashboard
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

The dashboard example demonstrates how Twerge can be used in complex web applications with multiple components, layouts, and dynamic data integration.