package tplimpl

import (
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/tpl"
	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
)

func (t *templateNamespace) readTemplateInto(templ *TemplInfo) error {
	if err := func() error {
		meta := templ.Fi.Meta()
		f, err := meta.Open()
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		templ.Content = removeLeadingBOM(string(b))
		if !templ.NoBaseOf {
			templ.NoBaseOf = !needsBaseTemplate(templ.Content)
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}

// The tweet and twitter shortcodes were deprecated in favor of the x shortcode
// in v0.141.0. We can remove these aliases in v0.155.0 or later.
var embeddedTemplatesAliases = map[string][]string{
	"_shortcodes/twitter.html": {"_shortcodes/tweet.html"},
}

func (t *templateNamespace) parseTemplate(ti *TemplInfo) error {
	if !ti.NoBaseOf || ti.Category == CategoryBaseof {
		// Delay parsing until we have the base template.
		return nil
	}
	pi := ti.PathInfo
	name := pi.PathNoLeadingSlash()
	if ti.isLegacyMapped {
		// When mapping the old taxonomy structure to the new one, we may map the same path to multiple templates per kind.
		// Append the kind here to make the name unique.
		name += ("-" + ti.D.Kind)
	}

	var (
		templ tpl.Template
		err   error
	)

	if ti.D.IsPlainText {
		prototype := t.parseText
		templ, err = prototype.New(name).Parse(ti.Content)
		if err != nil {
			return err
		}
	} else {
		prototype := t.parseHTML
		templ, err = prototype.New(name).Parse(ti.Content)
		if err != nil {
			return err
		}

		if ti.SubCategory == SubCategoryEmbedded {
			// In Hugo 0.146.0 we moved the internal templates around.
			// For the "_internal/twitter_cards.html" style templates, they
			// were moved to the _partials directory.
			// But we need to make them accessible from the old path for a while.
			if pi.Type() == paths.TypePartial {
				aliasName := strings.TrimPrefix(name, "_partials/")
				aliasName = "_internal/" + aliasName
				_, err = prototype.AddParseTree(aliasName, templ.(*htmltemplate.Template).Tree)
				if err != nil {
					return err
				}
			}

			// This was also possible before Hugo 0.146.0, but this should be deprecated.
			if pi.Type() == paths.TypeShortcode {
				aliasName := strings.TrimPrefix(name, "_shortcodes/")
				aliasName = "_internal/shortcodes/" + aliasName
				_, err = prototype.AddParseTree(aliasName, templ.(*htmltemplate.Template).Tree)
				if err != nil {
					return err
				}
			}

		}
	}

	ti.Template = templ

	return nil
}

func (t *templateNamespace) applyBaseTemplate(overlay *TemplInfo, base keyTemplateInfo) error {
	tb := &TemplWithBaseApplied{
		Overlay: overlay,
		Base:    base.Info,
	}

	base.Info.Overlays = append(base.Info.Overlays, overlay)

	var templ tpl.Template
	if overlay.D.IsPlainText {
		tt := texttemplate.Must(t.parseText.Clone()).New(overlay.PathInfo.PathNoLeadingSlash())
		var err error
		tt, err = tt.Parse(base.Info.Content)
		if err != nil {
			return err
		}
		tt, err = tt.Parse(overlay.Content)
		if err != nil {
			return err
		}
		templ = tt
		t.baseofTextClones = append(t.baseofTextClones, tt)
	} else {
		tt := htmltemplate.Must(t.parseHTML.CloneShallow()).New(overlay.PathInfo.PathNoLeadingSlash())
		var err error
		tt, err = tt.Parse(base.Info.Content)
		if err != nil {
			return err
		}
		tt, err = tt.Parse(overlay.Content)
		if err != nil {
			return err
		}
		templ = tt

		t.baseofHtmlClones = append(t.baseofHtmlClones, tt)

	}

	tb.Template = &TemplInfo{
		Template: templ,
		Base:     base.Info,
		PathInfo: overlay.PathInfo,
		Fi:       overlay.Fi,
		D:        overlay.D,
		NoBaseOf: true,
	}

	variants := overlay.BaseVariants.Get(base.Key)
	if variants == nil {
		variants = make(map[TemplateDescriptor]*TemplWithBaseApplied)
		overlay.BaseVariants.Insert(base.Key, variants)
	}
	variants[base.Info.D] = tb
	return nil
}

func (t *templateNamespace) templatesIn(in tpl.Template) []tpl.Template {
	var templs []tpl.Template
	if textt, ok := in.(*texttemplate.Template); ok {
		for _, t := range textt.Templates() {
			templs = append(templs, t)
		}
	}
	if htmlt, ok := in.(*htmltemplate.Template); ok {
		for _, t := range htmlt.Templates() {
			templs = append(templs, t)
		}
	}
	return templs
}

/*


func (t *templateHandler) applyBaseTemplate(overlay, base templateInfo) (tpl.Template, error) {
	if overlay.isText {
		var (
			templ = t.main.getPrototypeText(prototypeCloneIDBaseof).New(overlay.name)
			err   error
		)

		if !base.IsZero() {
			templ, err = templ.Parse(base.template)
			if err != nil {
				return nil, base.errWithFileContext("text: base: parse failed", err)
			}
		}

		templ, err = texttemplate.Must(templ.Clone()).Parse(overlay.template)
		if err != nil {
			return nil, overlay.errWithFileContext("text: overlay: parse failed", err)
		}

		// The extra lookup is a workaround, see
		// * https://github.com/golang/go/issues/16101
		// * https://github.com/gohugoio/hugo/issues/2549
		// templ = templ.Lookup(templ.Name())

		return templ, nil
	}

	var (
		templ = t.main.getPrototypeHTML(prototypeCloneIDBaseof).New(overlay.name)
		err   error
	)

	if !base.IsZero() {
		templ, err = templ.Parse(base.template)
		if err != nil {
			return nil, base.errWithFileContext("html: base: parse failed", err)
		}
	}

	templ, err = htmltemplate.Must(templ.Clone()).Parse(overlay.template)
	if err != nil {
		return nil, overlay.errWithFileContext("html: overlay: parse failed", err)
	}

	// The extra lookup is a workaround, see
	// * https://github.com/golang/go/issues/16101
	// * https://github.com/gohugoio/hugo/issues/2549
	templ = templ.Lookup(templ.Name())

	return templ, err
}

*/

var baseTemplateDefineRe = regexp.MustCompile(`^{{-?\s*define`)

// needsBaseTemplate returns true if the first non-comment template block is a
// define block.
func needsBaseTemplate(templ string) bool {
	idx := -1
	inComment := false
	for i := 0; i < len(templ); {
		if !inComment && strings.HasPrefix(templ[i:], "{{/*") {
			inComment = true
			i += 4
		} else if !inComment && strings.HasPrefix(templ[i:], "{{- /*") {
			inComment = true
			i += 6
		} else if inComment && strings.HasPrefix(templ[i:], "*/}}") {
			inComment = false
			i += 4
		} else if inComment && strings.HasPrefix(templ[i:], "*/ -}}") {
			inComment = false
			i += 6
		} else {
			r, size := utf8.DecodeRuneInString(templ[i:])
			if !inComment {
				if strings.HasPrefix(templ[i:], "{{") {
					idx = i
					break
				} else if !unicode.IsSpace(r) {
					break
				}
			}
			i += size
		}
	}

	if idx == -1 {
		return false
	}

	return baseTemplateDefineRe.MatchString(templ[idx:])
}

func removeLeadingBOM(s string) string {
	const bom = '\ufeff'

	for i, r := range s {
		if i == 0 && r != bom {
			return s
		}
		if i > 0 {
			return s[i:]
		}
	}

	return s
}

type templateNamespace struct {
	parseText     *texttemplate.Template
	parseHTML     *htmltemplate.Template
	prototypeText *texttemplate.Template
	prototypeHTML *htmltemplate.Template

	standaloneText *texttemplate.Template

	baseofTextClones []*texttemplate.Template
	baseofHtmlClones []*htmltemplate.Template
}

func (t *templateNamespace) createPrototypesParse() error {
	if t.prototypeHTML == nil {
		panic("prototypeHTML not set")
	}
	t.parseHTML = htmltemplate.Must(t.prototypeHTML.Clone())
	t.parseText = texttemplate.Must(t.prototypeText.Clone())
	return nil
}

func (t *templateNamespace) createPrototypes(init bool) error {
	if init {
		t.prototypeHTML = htmltemplate.Must(t.parseHTML.Clone())
		t.prototypeText = texttemplate.Must(t.parseText.Clone())
	}
	// t.execHTML = htmltemplate.Must(t.parseHTML.Clone())
	// t.execText = texttemplate.Must(t.parseText.Clone())

	return nil
}

func newTemplateNamespace(funcs map[string]any) *templateNamespace {
	return &templateNamespace{
		parseHTML:      htmltemplate.New("").Funcs(funcs),
		parseText:      texttemplate.New("").Funcs(funcs),
		standaloneText: texttemplate.New("").Funcs(funcs),
	}
}
