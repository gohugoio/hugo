package transform

import (
	"github.com/spf13/viper"
	"sync"
)

var absURLInit sync.Once
var ar *absURLReplacer

func AbsURL() (trs []link, err error) {
	initAbsURLReplacer()
	return absURLFromReplacer(ar)
}

func absURLFromURL(URL string) (trs []link, err error) {
	return absURLFromReplacer(newAbsURLReplacer(URL))
}

func absURLFromReplacer(ar *absURLReplacer) (trs []link, err error) {
	trs = append(trs, func(rw contentTransformer) {
		ar.replaceInHTML(rw)
	})
	return
}

func AbsURLInXML() (trs []link, err error) {
	initAbsURLReplacer()
	return absURLInXMLFromReplacer(ar)
}

func absURLInXMLFromURL(URL string) (trs []link, err error) {
	return absURLInXMLFromReplacer(newAbsURLReplacer(URL))
}

func absURLInXMLFromReplacer(ar *absURLReplacer) (trs []link, err error) {
	trs = append(trs, func(rw contentTransformer) {
		ar.replaceInXML(rw)
	})
	return
}

func initAbsURLReplacer() {
	absURLInit.Do(func() {
		ar = newAbsURLReplacer(viper.GetString("BaseURL"))
	})
}
