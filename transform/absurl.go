package transform

import (
	"github.com/spf13/viper"
	"sync"
)

// to be used in tests; the live site will get its value from Viper.
var AbsBaseUrl string

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
	trs = append(trs, func(ct contentTransformer) {
		ar.replaceInHTML(ct)
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
	trs = append(trs, func(ct contentTransformer) {
		ar.replaceInXML(ct)
	})
	return
}

func initAbsURLReplacer() {
	absURLInit.Do(func() {
		var url string

		if AbsBaseUrl != "" {
			url = AbsBaseUrl
		} else {
			url = viper.GetString("BaseURL")
		}

		ar = newAbsURLReplacer(url)
	})
}
