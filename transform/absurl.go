package transform

import (
	"sync"
)

var absUrlInit sync.Once
var ar *absurlReplacer

// for performance reasons, we reuse the first baseUrl given
func initAbsurlReplacer(baseURL string) {
	absUrlInit.Do(func() {
		ar = newAbsurlReplacer(baseURL)
	})
}

func AbsURL(absURL string) (trs []link, err error) {
	initAbsurlReplacer(absURL)

	trs = append(trs, func(content []byte) []byte {
		return ar.replaceInHtml(content)
	})
	return
}

func AbsURLInXML(absURL string) (trs []link, err error) {
	initAbsurlReplacer(absURL)

	trs = append(trs, func(content []byte) []byte {
		return ar.replaceInXml(content)
	})
	return
}
