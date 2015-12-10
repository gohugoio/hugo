package hugolib

import (
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

type Sitemap struct {
	ChangeFreq string
	Priority   float64
	Filename   string
}

func parseSitemap(input map[string]interface{}) Sitemap {
	sitemap := Sitemap{Priority: -1, Filename: "sitemap.xml"}

	for key, value := range input {
		switch key {
		case "changefreq":
			sitemap.ChangeFreq = cast.ToString(value)
		case "priority":
			sitemap.Priority = cast.ToFloat64(value)
		case "filename":
			sitemap.Filename = cast.ToString(value)
		default:
			jww.WARN.Printf("Unknown Sitemap field: %s\n", key)
		}
	}

	return sitemap
}
