package transform

import (
	"bytes"
	"github.com/spf13/viper"
)

func LiveReloadInject(ct contentTransformer) {
	match := []byte("</body>")
	port := viper.GetString("port")
	replace := []byte(`<script data-no-instant>document.write('<script src="http://'
        + (location.host || 'localhost').split(':')[0]
		+ ':` + port + `/livereload.js?mindelay=10"></'
        + 'script>')</script></body>`)
	newcontent := bytes.Replace(ct.Content(), match, replace, -1)

	if len(newcontent) == len(ct.Content()) {
		match := []byte("</BODY>")
		newcontent = bytes.Replace(ct.Content(), match, replace, -1)
	}

	ct.Write(newcontent)
}
