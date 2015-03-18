package transform

import (
	"bytes"
	"github.com/spf13/viper"
)

func LiveReloadInject(rw contentRewriter) {
	match := []byte("</body>")
	port := viper.GetString("port")
	replace := []byte(`<script>document.write('<script src="http://'
        + (location.host || 'localhost').split(':')[0]
		+ ':` + port + `/livereload.js?mindelay=10"></'
        + 'script>')</script></body>`)
	newcontent := bytes.Replace(rw.Content(), match, replace, -1)

	if len(newcontent) == len(rw.Content()) {
		match := []byte("</BODY>")
		newcontent = bytes.Replace(rw.Content(), match, replace, -1)
	}

	rw.Write(newcontent)
}
