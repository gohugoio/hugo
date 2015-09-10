package transform

var ar *absURLReplacer = newAbsURLReplacer()

var AbsURL = func(ct contentTransformer) {
	ar.replaceInHTML(ct)
}

var AbsURLInXML = func(ct contentTransformer) {
	ar.replaceInXML(ct)
}
