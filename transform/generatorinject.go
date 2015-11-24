package transform

import (
	"bytes"
	"fmt"

	"github.com/spf13/hugo/helpers"
)

func GeneratorInject(ct contentTransformer) {
	match := []byte("</head>")
	replace := []byte(fmt.Sprintf(`<meta name="generator" content="Hugo %s" /></head>`, helpers.HugoVersion()))

	newcontent := bytes.Replace(ct.Content(), match, replace, -1)
	if len(newcontent) == len(ct.Content()) {
		match := []byte("</HEAD>")
		newcontent = bytes.Replace(ct.Content(), match, replace, -1)
	}

	ct.Write(newcontent)
}
