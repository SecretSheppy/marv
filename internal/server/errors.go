package server

import "errors"

var (
	ErrFailedToReadFile      = errors.New("failed to read file")
	ErrFailedToMinifyOrCache = errors.New("failed to minify or cache")
	ErrFailedToInitRender    = errors.New("failed to init render")
	ErrFailedToRenderHTML    = errors.New("failed to render html")
)
