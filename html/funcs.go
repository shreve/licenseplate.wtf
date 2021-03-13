package html

import (
	"text/template"
	"time"
)

var FuncMap = template.FuncMap{
	"time": func(format string) string {
		return time.Now().Format(format)
	},
}
