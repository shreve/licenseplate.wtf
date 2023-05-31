package html

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"
	"text/template"

	"licenseplate.wtf/model"
	"licenseplate.wtf/util"
)

//go:embed templates/*.html templates/*/*.html
var files embed.FS
var tmplFS fs.FS
var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)

	if os.Getenv("ENV") == "production" {
		tmplFS = files
	} else {
		_, file, _, _ := runtime.Caller(0)
		dir := path.Dir(file)
		tmplFS = os.DirFS(dir)
	}
}

func parse(file string) *template.Template {
	file = path.Join("templates", file)
	log.Printf("Parsing %s from %v", file, tmplFS)
	tmpl, err := template.New("layout.html").Funcs(FuncMap).ParseFS(
		tmplFS, "templates/layout.html", file,
	)

	return template.Must(tmpl, err)
}

func getTemplate(file string) *template.Template {
	if os.Getenv("ENV") != "production" {
		return parse(file)
	}
	if tmpl, ok := templates[file]; ok {
		return tmpl
	}
	tmpl := parse(file)
	templates[file] = tmpl
	return tmpl
}

func Home(w io.Writer) {
	data := genericParams{
		Page: PageData{
			Title:     "WTF is licenseplate.wtf",
			Canonical: fullURL(),
		},
	}
	log.Println("Rendering home.html")
	getTemplate("home.html").Execute(w, data)
}

func PlateShow(w io.Writer, data ParamsMap) {
	plate := data["Plate"].(*model.Plate)
	code := plate.Code
	data["Page"] = PageData{
		Title:     "What does the license plate " + code + " mean?",
		Canonical: fullURL(plate.URL()),
	}
	util.LogTime("Rendering plate/show.html", func() {
		if err := getTemplate("plates/show.html").Execute(w, data); err != nil {
			log.Println(err)
		}
	})
}

func PlateList(w io.Writer, data ParamsMap) {
	data["Page"] = PageData{
		Title:     "All license plates",
		Canonical: fullURL("plates"),
	}
	data["Plates"] = data["Plates"].([]model.Plate)

	util.LogTime("Rendering plate/list.html", func() {
		if err := getTemplate("plates/list.html").Execute(w, data); err != nil {
			log.Println(err)
		}
	})
}
