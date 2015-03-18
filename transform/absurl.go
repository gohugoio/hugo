package transform

import (
	"sync"
)

var absURLInit sync.Once
var ar *absURLReplacer

// for performance reasons, we reuse the first baseURL given
func initAbsURLReplacer(baseURL string) {
	absURLInit.Do(func() {
		ar = newAbsURLReplacer(baseURL)
	})
}

func AbsURL(absURL string) (trs []link, err error) {
	initAbsURLReplacer(absURL)

	trs = append(trs, func(content []byte) []byte {
		return ar.replaceInHTML(content)
	})
	return
}

func AbsURLInXML(absURL string) (trs []link, err error) {
	initAbsURLReplacer(absURL)

	trs = append(trs, func(content []byte) []byte {
		return ar.replaceInXML(content)
	})
	return
}
