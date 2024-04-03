package golang

import "embed"

//go:embed *.go.tmpl
var FS embed.FS
