package media

type BuiltinTypes struct {
	CalendarType   Type
	CSSType        Type
	SCSSType       Type
	SASSType       Type
	CSVType        Type
	HTMLType       Type
	JavascriptType Type
	TypeScriptType Type
	TSXType        Type
	JSXType        Type

	JSONType           Type
	WebAppManifestType Type
	RSSType            Type
	XMLType            Type
	SVGType            Type
	TextType           Type
	TOMLType           Type
	YAMLType           Type

	// Common image types
	PNGType  Type
	JPEGType Type
	GIFType  Type
	TIFFType Type
	BMPType  Type
	WEBPType Type

	// Common font types
	TrueTypeFontType Type
	OpenTypeFontType Type

	// Common document types
	PDFType              Type
	MarkdownType         Type
	EmacsOrgModeType     Type
	AsciiDocType         Type
	PandocType           Type
	ReStructuredTextType Type

	// Common video types
	AVIType  Type
	MPEGType Type
	MP4Type  Type
	OGGType  Type
	WEBMType Type
	GPPType  Type

	// wasm
	WasmType Type

	OctetType Type
}

var Builtin = BuiltinTypes{
	CalendarType:   Type{Type: "text/calendar"},
	CSSType:        Type{Type: "text/css"},
	SCSSType:       Type{Type: "text/x-scss"},
	SASSType:       Type{Type: "text/x-sass"},
	CSVType:        Type{Type: "text/csv"},
	HTMLType:       Type{Type: "text/html"},
	JavascriptType: Type{Type: "text/javascript"},
	TypeScriptType: Type{Type: "text/typescript"},
	TSXType:        Type{Type: "text/tsx"},
	JSXType:        Type{Type: "text/jsx"},

	JSONType:           Type{Type: "application/json"},
	WebAppManifestType: Type{Type: "application/manifest+json"},
	RSSType:            Type{Type: "application/rss+xml"},
	XMLType:            Type{Type: "application/xml"},
	SVGType:            Type{Type: "image/svg+xml"},
	TextType:           Type{Type: "text/plain"},
	TOMLType:           Type{Type: "application/toml"},
	YAMLType:           Type{Type: "application/yaml"},

	// Common image types
	PNGType:  Type{Type: "image/png"},
	JPEGType: Type{Type: "image/jpeg"},
	GIFType:  Type{Type: "image/gif"},
	TIFFType: Type{Type: "image/tiff"},
	BMPType:  Type{Type: "image/bmp"},
	WEBPType: Type{Type: "image/webp"},

	// Common font types
	TrueTypeFontType: Type{Type: "font/ttf"},
	OpenTypeFontType: Type{Type: "font/otf"},

	// Common document types
	PDFType:              Type{Type: "application/pdf"},
	MarkdownType:         Type{Type: "text/markdown"},
	AsciiDocType:         Type{Type: "text/asciidoc"}, // https://github.com/asciidoctor/asciidoctor/issues/2502
	PandocType:           Type{Type: "text/pandoc"},
	ReStructuredTextType: Type{Type: "text/rst"}, // https://docutils.sourceforge.io/FAQ.html#what-s-the-official-mime-type-for-restructuredtext-data
	EmacsOrgModeType:     Type{Type: "text/org"},

	// Common video types
	AVIType:  Type{Type: "video/x-msvideo"},
	MPEGType: Type{Type: "video/mpeg"},
	MP4Type:  Type{Type: "video/mp4"},
	OGGType:  Type{Type: "video/ogg"},
	WEBMType: Type{Type: "video/webm"},
	GPPType:  Type{Type: "video/3gpp"},

	// Web assembly.
	WasmType: Type{Type: "application/wasm"},

	OctetType: Type{Type: "application/octet-stream"},
}

var defaultMediaTypesConfig = map[string]any{
	"text/calendar":   map[string]any{"suffixes": []string{"ics"}},
	"text/css":        map[string]any{"suffixes": []string{"css"}},
	"text/x-scss":     map[string]any{"suffixes": []string{"scss"}},
	"text/x-sass":     map[string]any{"suffixes": []string{"sass"}},
	"text/csv":        map[string]any{"suffixes": []string{"csv"}},
	"text/html":       map[string]any{"suffixes": []string{"html", "htm"}},
	"text/javascript": map[string]any{"suffixes": []string{"js", "jsm", "mjs"}},
	"text/typescript": map[string]any{"suffixes": []string{"ts"}},
	"text/tsx":        map[string]any{"suffixes": []string{"tsx"}},
	"text/jsx":        map[string]any{"suffixes": []string{"jsx"}},

	"application/json":          map[string]any{"suffixes": []string{"json"}},
	"application/manifest+json": map[string]any{"suffixes": []string{"webmanifest"}},
	"application/rss+xml":       map[string]any{"suffixes": []string{"xml", "rss"}},
	"application/xml":           map[string]any{"suffixes": []string{"xml"}},
	"image/svg+xml":             map[string]any{"suffixes": []string{"svg"}},
	"text/plain":                map[string]any{"suffixes": []string{"txt"}},
	"application/toml":          map[string]any{"suffixes": []string{"toml"}},
	"application/yaml":          map[string]any{"suffixes": []string{"yaml", "yml"}},

	// Common image types
	"image/png":  map[string]any{"suffixes": []string{"png"}},
	"image/jpeg": map[string]any{"suffixes": []string{"jpg", "jpeg", "jpe", "jif", "jfif"}},
	"image/gif":  map[string]any{"suffixes": []string{"gif"}},
	"image/tiff": map[string]any{"suffixes": []string{"tif", "tiff"}},
	"image/bmp":  map[string]any{"suffixes": []string{"bmp"}},
	"image/webp": map[string]any{"suffixes": []string{"webp"}},

	// Common font types
	"font/ttf": map[string]any{"suffixes": []string{"ttf"}},
	"font/otf": map[string]any{"suffixes": []string{"otf"}},

	// Common document types
	"application/pdf": map[string]any{"suffixes": []string{"pdf"}},
	"text/markdown":   map[string]any{"suffixes": []string{"md", "mdown", "markdown"}},
	"text/asciidoc":   map[string]any{"suffixes": []string{"adoc", "asciidoc", "ad"}},
	"text/pandoc":     map[string]any{"suffixes": []string{"pandoc", "pdc"}},
	"text/rst":        map[string]any{"suffixes": []string{"rst"}},
	"text/org":        map[string]any{"suffixes": []string{"org"}},

	// Common video types
	"video/x-msvideo": map[string]any{"suffixes": []string{"avi"}},
	"video/mpeg":      map[string]any{"suffixes": []string{"mpg", "mpeg"}},
	"video/mp4":       map[string]any{"suffixes": []string{"mp4"}},
	"video/ogg":       map[string]any{"suffixes": []string{"ogv"}},
	"video/webm":      map[string]any{"suffixes": []string{"webm"}},
	"video/3gpp":      map[string]any{"suffixes": []string{"3gpp", "3gp"}},

	// wasm
	"application/wasm": map[string]any{"suffixes": []string{"wasm"}},

	"application/octet-stream": map[string]any{},
}
