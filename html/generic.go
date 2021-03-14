package html

import (
	"strings"
)

type PageData struct {
	Title     string
	Canonical string
}

type genericParams struct {
	Page PageData
}

var domain = "licenseplate.wtf"

func fullURL(bits ...string) string {
	bits = append([]string{"https:/", domain}, bits...)
	return strings.Join(bits, "/")
}
