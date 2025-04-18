// Package main is a simple example of how to use twerge.
// It generates some class names and saves them to a file.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/conneroisu/twerge"
)

func main() {
	// Generate some class names
	class1 := twerge.Merge("text-red-500 bg-blue-500")
	class2 := twerge.Merge("text-green-300 p-4")
	class3 := twerge.Merge("flex items-center justify-between")
	fmt.Println("Generated class names:")
	fmt.Printf("text-red-500 bg-blue-500 -> %s\n", class1)
	fmt.Printf("text-green-300 p-4 -> %s\n", class2)
	fmt.Printf("flex items-center justify-between -> %s\n", class3)

	// Test class merging functionality
	fmt.Println("\nMerged classes:")
	merged := twerge.Merge("text-red-500 text-blue-700")
	fmt.Printf("text-red-500 text-blue-700 -> %s\n", merged)

	// Class conflict resolution
	merged = twerge.Merge("p-4 p-8")
	fmt.Printf("p-4 p-8 -> %s\n", merged)

	// Generate and save the class map code
	outPath := filepath.Join(os.TempDir(), "class_map_generated.go")
	fmt.Printf("\nGenerating class map code to %s\n", outPath)

	// Generate class names using Generate instead of Merge to populate ClassMapStr
	fmt.Println("\nAdding entries to ClassMapStr:")
	class1Gen := twerge.It("text-red-500 bg-blue-500")
	class2Gen := twerge.It("text-green-300 p-4")
	class3Gen := twerge.It("flex items-center justify-between")
	class4Gen := twerge.It("text-red-500 text-blue-700")
	class5Gen := twerge.It("p-4 p-8")

	fmt.Printf("text-red-500 bg-blue-500 -> %s\n", class1Gen)
	fmt.Printf("text-green-300 p-4 -> %s\n", class2Gen)
	fmt.Printf("flex items-center justify-between -> %s\n", class3Gen)
	fmt.Printf("text-red-500 text-blue-700 -> %s\n", class4Gen)
	fmt.Printf("p-4 p-8 -> %s\n", class5Gen)

	// Generate the code with the populated ClassMapStr
	err := twerge.GenerateGo("./class_map_generated.go")
	if err != nil {
		fmt.Println(err)
		return
	}

}
