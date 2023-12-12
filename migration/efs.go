package migration

import (
	"embed"
)

//go:embed "sqlite"
var EmbeddedFiles embed.FS
