package twerge

import (
	"crypto/sha1"
	"encoding/base64"
	"os"
	"strings"
	"sync"

	"github.com/dave/jennifer/jen"
)

// ClassMap is a mapping of original class strings to generated class names
type ClassMap map[string]string

var (
	// Global mapping of original class strings to generated class names
	globalClassMap = make(ClassMap)

	// Mutex to protect globalClassMap
	mapMutex sync.RWMutex

	// cache for generated classes
	genCache = newCache(1000)
)

// Generate creates a short unique CSS class name from the merged classes
func Generate(classes string) string {
	// First, merge the classes
	merged := Merge(classes)

	// Check if we've already generated a class for this set of merged classes
	mapMutex.RLock()
	if cachedClass := genCache.Get(merged); cachedClass != "" {
		mapMutex.RUnlock()
		return cachedClass
	}
	mapMutex.RUnlock()

	// Generate a hash of the merged classes
	hash := sha1.Sum([]byte(merged))

	// Use URL-safe base64 encoding and trim to 7 characters for brevity
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	classname := "tw-" + encoded[:7]

	// Store the mapping
	mapMutex.Lock()
	globalClassMap[classes] = classname
	genCache.Set(merged, classname)
	mapMutex.Unlock()

	return classname
}

// GetMapping returns the current mapping from original class strings to generated class names
func GetMapping() ClassMap {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	// Create a copy to avoid concurrent map access issues
	mapping := make(ClassMap, len(globalClassMap))
	for k, v := range globalClassMap {
		mapping[k] = v
	}

	return mapping
}

// GenerateClassMapCode generates Go code for a variable containing the class mapping
func GenerateClassMapCode() string {
	mapping := GetMapping()

	// Create a new file
	f := jen.NewFile("twerge")

	// Add a package comment
	f.PackageComment("Code generated by twerge. DO NOT EDIT.")

	// Create the ClassMapStr variable
	f.Var().Id("ClassMapStr").Op("=").Map(jen.String()).String().Values(jen.DictFunc(func(d jen.Dict) {
		// Sort keys for deterministic output
		var keys []string
		for k := range mapping {
			keys = append(keys, k)
		}

		// Add each key-value pair
		for _, k := range keys {
			d[jen.Lit(k)] = jen.Lit(mapping[k])
		}
	}))

	// Generate the code
	buf := &strings.Builder{}
	err := f.Render(buf)
	if err != nil {
		return "// Error generating code: " + err.Error()
	}

	return buf.String()
}

// WriteClassMapFile writes the generated class map to the specified file
func WriteClassMapFile(filepath string) error {
	code := GenerateClassMapCode()

	// Just write the code directly to file
	return os.WriteFile(filepath, []byte(code), 0644)
}

