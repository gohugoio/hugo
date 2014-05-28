package transform

import (
	"bytes"

	"github.com/spf13/viper"
)

func LiveReloadInject(content []byte) []byte {
	match := []byte("</body>")
	port := viper.GetString("port")
	replace := []byte(`<script>document.write('<script src="http://'
        + (location.host || 'localhost').split(':')[0]
		+ ':` + port + `/livereload.js?mindelay=10"></'
        + 'script>')</script></body>`)
	newcontent := bytes.Replace(content, match, replace, -1)

	if len(newcontent) == len(content) {
		match := []byte("</BODY>")
		newcontent = bytes.Replace(content, match, replace, -1)
	}

	return newcontent
}
