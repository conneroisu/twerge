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
	"github.com/conneroisu/twerge/examples/admin-dashboard/views"
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
	
	// Ensure directories exist
	if err := os.MkdirAll("static/dist", 0755); err != nil {
		panic(err)
	}
	
	fmt.Println("Updating Generated Code...")
	start = time.Now()
	
	// Generate code for all dashboard states
	if err := twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		"input.css",
		"classes/classes.html",
		// Dashboard with sidebar open
		views.Dashboard(false),
		// Dashboard with sidebar collapsed
		views.Dashboard(true),
		// Individual components with various states
		views.Sidebar(false),
		views.Sidebar(true),
		views.Header(),
		views.MetricsGrid(),
		views.NavItem("Dashboard", "/", true, false),
		views.NavItem("Users", "/users", false, true),
		views.NavItem("Products", "/products", false, true),
		views.NavItem("Orders", "/orders", false, false),
		views.MetricCard("Users", "12,345", "↗ 12%", true),
		views.MetricCard("Revenue", "$56,789", "↘ 8%", false),
		views.DataTable([]string{"Order ID", "Customer", "Amount", "Status"}, "Amount", true),
		views.DataTable([]string{"Product", "Sales", "Revenue", "Stock"}, "Revenue", false),
		views.TableHeader("Sample", true, true),
		views.TableHeader("Sample", true, false),
		views.TableHeader("Sample", false, false),
		views.StatusBadge(true),
		views.StatusBadge(false),
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
	cmd := exec.Command("tailwindcss", "-i", "input.css", "-o", "static/dist/styles.css")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}