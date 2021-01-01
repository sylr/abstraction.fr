package static

import (
	"embed"
)

//go:embed css images favicon.ico
var StaticFS embed.FS
