package html

import (
	"strings"
)

type PageData struct {
	Title     string
	Canonical string
}

type genericParams struct {
	Page   PageData
	Errors []string
}

var domain = "licenseplate.wtf"
var baseURL = "https://" + domain

func fullURL(bits ...string) string {
	bits = append([]string{"https:/", domain}, bits...)
	return strings.Join(bits, "/")
}

type ParamsMap map[string]interface{}
