package html

import (
	"embed"
	"io"
	"log"
	"text/template"

	"licenseplate.wtf/model"
	"licenseplate.wtf/util"
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
	data := genericParams{
		Page: PageData{
			Title:     "WTF is licenseplate.wtf",
			Canonical: fullURL(),
		},
	}
	log.Println("Rendering home.html")
	home.Execute(w, data)
}

var plateShow = parse("plates/show.html")

func PlateShow(w io.Writer, data ParamsMap) {
	plate := data["Plate"].(*model.Plate)
	code := plate.Code
	data["Page"] = PageData{
		Title:     "What does the license plate " + code + " mean?",
		Canonical: fullURL(plate.URL()),
	}
	util.LogTime("Rendering plate/show.html", func() {
		if err := plateShow.Execute(w, data); err != nil {
			log.Println(err)
		}
	})
}

var plateList = parse("plates/list.html")

func PlateList(w io.Writer, data ParamsMap) {
	data["Page"] = PageData{
		Title:     "All license plates",
		Canonical: fullURL("plates"),
	}
	data["Plates"] = data["Plates"].([]model.Plate)

	util.LogTime("Rendering plate/list.html", func() {
		if err := plateList.Execute(w, data); err != nil {
			log.Println(err)
		}
	})
}
