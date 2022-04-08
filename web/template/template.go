package template

import "embed"

//go:embed *.template.html
//go:embed *.template.js
var Content embed.FS
