package hugolib

import (
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

type Sitemap struct {
	ChangeFreq string
	Priority   float64
}

func parseSitemap(input map[string]interface{}) Sitemap {
	sitemap := Sitemap{Priority: -1}

	for key, value := range input {
		switch key {
		case "changefreq":
			sitemap.ChangeFreq = cast.ToString(value)
		case "priority":
			sitemap.Priority = cast.ToFloat64(value)
		default:
			jww.WARN.Printf("Unknown Sitemap field: %s\n", key)
		}
	}

	return sitemap
}
