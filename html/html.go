package html

import (
	"embed"
	"io"
	"text/template"
)

//go:embed *.html */*.html
var files embed.FS

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").Funcs(FuncMap).ParseFS(
			files, "layout.html", file,
		),
	)
}

var home = parse("home.html")

func Home(w io.Writer) {
	home.Execute(w, struct{}{})
}

var plateShow = parse("plates/show.html")

func PlateShow(w io.Writer) {
	plateShow.Execute(w, struct{}{})
}

var plateList = parse("plates/list.html")

func PlateList(w io.Writer) {
	plateList.Execute(w, struct{}{})
}
