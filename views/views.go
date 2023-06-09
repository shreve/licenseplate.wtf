package views

import (
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"

	"github.com/kataras/blocks"
)

func renderString(in string, params *Params) string {
	out := &strings.Builder{}
	tmpl := template.Must(template.New("").Parse(in))
	tmpl.Execute(out, params)
	return out.String()
}

var views *blocks.Blocks

func init() {

	if os.Getenv("ENV") == "production" {
		views = blocks.New(files)
	} else {
		_, file, _, _ := runtime.Caller(0)
		views = blocks.New(path.Dir(file)).Reload(true)
	}

	views = views.RootDir("html").Funcs(FuncMap).DefaultLayout("main")

	err := views.Load()
	if err != nil {
		log.Fatalf("Error loading views: %v", err)
	}

	names := make([]string, 0, len(views.Templates))
	for n := range views.Templates {
		names = append(names, n)
	}
	log.Println("Loaded views:", names)
}

type Params map[string]interface{}

type View struct {
	Title       string
	Description string
	Canonical   string
	Template    string
	Layout      string
}

func (v View) Render(out io.Writer, params Params) error {
	if params == nil {
		params = Params{}
	}
	params["Template"] = v.Template
	params["Layout"] = v.Layout
	params["Title"] = renderString(v.Title, &params)
	params["Description"] = renderString(v.Description, &params)
	params["Canonical"] = renderString(v.Canonical, &params)
	log.Println("Rendering template", v.Template)
	return views.ExecuteTemplate(out, v.Template, v.Layout, params)
}

var Home = View{
	Template:    "home",
	Title:       "WTF is licenseplate.wtf?",
	Description: "licenseplate.wtf is a place to share your thoughts on custom license plates.",
	Canonical:   "https://licenseplate.wtf",
}

var PlateShow = View{
	Template:    "plates/show",
	Title:       "What does the license plate [{{.Plate.Code}}] mean?",
	Description: "See what other people think about the license plate [{{.Plate.Code}}].",
	Canonical:   "https://licenseplate.wtf/plates/{{.Plate.Code}}",
}

var PlateList = View{
	Template:    "plates/list",
	Title:       "WTF do these plates mean?",
	Description: "See what other people think about these license plates.",
	Canonical:   "https://licenseplate.wtf/plates",
}
