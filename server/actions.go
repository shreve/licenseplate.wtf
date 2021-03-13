package server

import (
	"licenseplate.wtf/html"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	html.Home(w)
}

func plateList(w http.ResponseWriter, r *http.Request) {
	html.PlateList(w)
}

func plateShow(w http.ResponseWriter, r *http.Request) {
	html.PlateShow(w)
}
