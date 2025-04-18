package twerge

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceBetweenMarkers(t *testing.T) {
	// Test with existing markers
	original := []byte("Some content\n" + twergeBeginMarker + "\nold content\n" + twergeEndMarker + "\nMore content")
	replacement := []byte("new content")

	result, err := replaceBetweenMarkers(original, replacement)
	assert.NoError(t, err)
	assert.Contains(t, string(result), "new content")
	assert.NotContains(t, string(result), "old content")

	// Test with no markers
	original = []byte("Some content without markers")
	result, err = replaceBetweenMarkers(original, replacement)
	assert.NoError(t, err)
	assert.Contains(t, string(result), twergeBeginMarker)
	assert.Contains(t, string(result), twergeEndMarker)
	assert.Contains(t, string(result), "new content")
}

func TestGenerateInputCSSForTailwind(t *testing.T) {
	// Create temporary input and output files
	inputFile, err := os.CreateTemp("", "twerge-input-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer func() { _ = os.Remove(inputFile.Name()) }()

	// templ file
	templFile, err := os.CreateTemp("", "twerge-templ-*.templ")
	if err != nil {
		t.Fatalf("Failed to create temp templ file: %v", err)
	}
	// defer print(templFile.Name())
	defer func() { _ = os.Remove(templFile.Name()) }()

	// Write some content to the input file
	inputContent := `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
.custom-class {
  color: blue;
}

` + twergeBeginMarker + `
/* Old generated content */
` + twergeEndMarker + `

/* More styles */
`
	err = os.WriteFile(inputFile.Name(), []byte(inputContent), 0644)
	assert.NoError(t, err)

	// Create a test class map
	GenClassMergeStr = map[string]string{
		"tw-test1": "text-red-500",
		"tw-test2": "bg-blue-500",
	}

	// Generate input CSS
	err = GenerateCSS(inputFile.Name())
	assert.NoError(t, err)

	// Read the output file
	outputContent, err := os.ReadFile(inputFile.Name())
	assert.NoError(t, err)

	// Check content
	outputStr := string(outputContent)
	assert.Contains(t, outputStr, "@tailwind base")
	assert.Contains(t, outputStr, ".custom-class")
	assert.Contains(t, outputStr, ".tw-test1")
	assert.Contains(t, outputStr, ".tw-test2")
	assert.NotContains(t, outputStr, "Old generated content")
	assert.Contains(t, outputStr, "/* More styles */")

	err = GenerateHTML(templFile.Name())
	assert.NoError(t, err)
}

func TestGenerate(t *testing.T) {
	// Reset the class map for testing
	mapMutex.Lock()
	ClassMapStr = make(map[string]string)
	mapMutex.Unlock()

	// Test that Generate creates a consistent class name for the same input
	class1 := It("text-red-500 bg-blue-500")
	class2 := It("text-red-500 bg-blue-500")
	assert.Equal(t, class1, class2, "Generate should return the same class name for the same input")

	// Test that Generate handles class merging correctly
	class3 := It("text-red-500 text-blue-700")
	assert.NotEqual(t, class1, class3, "Generate should return different class names for different inputs")

	// Test that the generated class name format is correct
	assert.True(t, strings.HasPrefix(class1, "tw-"), "Generated class should start with 'tw-'")
}

func TestGetMapping(t *testing.T) {
	// Reset the class map for testing
	mapMutex.Lock()
	ClassMapStr = make(map[string]string)
	mapMutex.Unlock()

	// Generate some class names and store them directly in the map for testing
	class1 := "tw-abcdefg"
	class2 := "tw-hijklmn"

	mapMutex.Lock()
	ClassMapStr["text-red-500 bg-blue-500"] = class1
	ClassMapStr["text-green-300 p-4"] = class2
	mapMutex.Unlock()

	// Get the mapping
	mapping := getMapping()

	// Check that the mapping contains the expected entries
	assert.Equal(t, class1, mapping["text-red-500 bg-blue-500"], "Mapping should contain the original class string and generated class name")
	assert.Equal(t, class2, mapping["text-green-300 p-4"], "Mapping should contain the original class string and generated class name")
	assert.Equal(t, 2, len(mapping), "Mapping should contain 2 entries")
}
