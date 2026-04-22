package web

import "embed"

//go:embed scripts
var ScriptsFS embed.FS

//go:embed static
var StaticFS embed.FS

//go:embed styles
var StylesFS embed.FS
