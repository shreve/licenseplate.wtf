//go:build !dev

package views

import "embed"

//go:embed html/*/*.html
var files embed.FS
