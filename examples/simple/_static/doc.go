// Package static contains static assets for the example application, simple.
package static

import "embed"

//go:embed dist/*
var Assets embed.FS
