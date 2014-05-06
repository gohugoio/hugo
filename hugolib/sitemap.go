package hugolib

import jww "github.com/spf13/jwalterweatherman"

type Sitemap struct {
	ChangeFreq string
	Priority   float32
}

func (s Sitemap) Validate() {
	if s.Priority < 0 {
		jww.WARN.Printf("Sitemap priority should be greater than 0, found: %f", s.Priority)
		s.Priority = 0
	} else if s.Priority > 1 {
		jww.WARN.Printf("Sitemap priority should be lesser than 1, found: %f", s.Priority)
		s.Priority = 1
	}
}
