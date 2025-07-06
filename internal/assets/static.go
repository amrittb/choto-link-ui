package assets

import "embed"

//go:embed web/static/*
var StaticFS embed.FS
