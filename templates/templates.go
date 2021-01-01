package templates

import (
	"embed"
)

var (
	// StaticFS contains templates
	//go:embed *.html
	StaticFS embed.FS
)
