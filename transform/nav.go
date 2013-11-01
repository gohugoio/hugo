package transform

import (
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"fmt"
)

func NavActive(section, attrName string) (tr []*htmltran.Transform) {
	ma := htmltran.MustTrans(htmltran.ModifyAttrib("class", "active"), fmt.Sprintf("li[%s=%s]", attrName, section))
	tr = append(tr, ma)
	return
}
