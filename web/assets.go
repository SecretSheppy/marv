package web

import "embed"

//go:embed icons
var IconsFS embed.FS

//go:embed scripts
var ScriptsFS embed.FS

//go:embed static
var StaticFS embed.FS

//go:embed styles
var StylesFS embed.FS

//go:embed themes
var ThemesFS embed.FS
