package assets

import "embed"

//go:embed web/templates/*.tmpl
var TemplateFS embed.FS
