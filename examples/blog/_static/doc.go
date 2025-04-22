// Package static contains static assets for the blog example.
package static

import "embed"

//go:embed dist
var Assets embed.FS